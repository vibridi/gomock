package templates

import (
	"go/ast"
	"strconv"
	"strings"

	"github.com/vibridi/gomock/v3/helper"
)

type Data struct {
	Qualify     bool
	Export      bool
	Package     string
	ServiceName string
	FuncDefs    []*FuncDef
	UnnamedSig  bool
	Underlying  map[string]string
}

func (td *Data) ToFuncDef(field *ast.Field) *FuncDef {

	fn := field.Type.(*ast.FuncType)

	funcDef := &FuncDef{}
	funcDef.ServiceName = td.ServiceName
	funcDef.Name = field.Names[0].Name

	paramNames := make([]ParamName, 0, len(fn.Params.List))
	paramTypes := make([]string, 0, len(fn.Params.List))

	for i, p := range fn.Params.List {
		if len(p.Names) == 0 {
			paramNames = append(paramNames, td.expressionName(p.Type, "p"+strconv.Itoa(i)))
			paramTypes = append(paramTypes, td.expressionType(p.Type))

		} else {
			for _, n := range p.Names {
				paramNames = append(paramNames, td.expressionName(p.Type, n.Name))
				paramTypes = append(paramTypes, td.expressionType(p.Type))
			}
		}
	}

	if !td.UnnamedSig {
		funcDef.Signature = strings.Join(helper.Zips(justNames(paramNames), paramTypes, " "), ", ")
	} else {
		funcDef.Signature = strings.Join(paramTypes, ", ")
	}

	funcDef.Args = strings.Join(expandNames(paramNames), ", ")

	if fn.Results == nil {
		return funcDef
	}

	returnTypes := make([]string, 0, len(fn.Results.List))
	returnValues := make([]string, 0, len(fn.Results.List))

	for _, r := range fn.Results.List {
		returnTypes = append(returnTypes, td.expressionType(r.Type))
		returnValues = append(returnValues, td.returnValue(r.Type))
	}

	funcDef.Return = helper.ReturnTypesToString(returnTypes)
	funcDef.ReturnValues = strings.Join(returnValues, ", ")

	return funcDef
}

func (td *Data) expressionType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if td.Qualify && ast.IsExported(t.Name) {
			return td.Package + "." + t.Name
		}
		return t.Name

	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name

	case *ast.FuncType:
		return td.functionType(t)

	case *ast.ArrayType:
		return arrayLength(t) + td.expressionType(t.Elt)

	case *ast.StarExpr:
		return "*" + td.expressionType(t.X)

	case *ast.MapType:
		return "map[" + td.expressionType(t.Key) + "]" + td.expressionType(t.Value)

	case *ast.StructType:
		return "struct{}"

	case *ast.InterfaceType:
		return "interface{}"

	case *ast.ChanType:
		return chanType(t) + " " + td.expressionType(t.Value)

	case *ast.Ellipsis:
		return "..." + td.expressionType(t.Elt)

	default:
		return ""
	}
}

func (td *Data) functionType(fn *ast.FuncType) string {
	s := "func("

	pdecl := make([]string, 0)
	for _, param := range fn.Params.List {
		pn := make([]string, 0)
		for _, n := range param.Names {
			pn = append(pn, n.Name)
		}
		p := ""
		if len(pn) > 0 {
			p += strings.Join(pn, ", ")
			p += " "
		}
		p += td.expressionType(param.Type)
		pdecl = append(pdecl, p)
	}
	s += strings.Join(pdecl, ", ")
	s += ")"

	rs := make([]string, 0)
	if fn.Results != nil {
		for _, r := range fn.Results.List {
			rs = append(rs, td.expressionType(r.Type))
		}
	}

	if len(rs) == 0 {
		return s
	}
	s += " ("
	s += strings.Join(rs, ",")
	s += ")"
	return s
}

func chanType(ch *ast.ChanType) string {
	switch ch.Dir {
	case 1:
		return "chan<-"
	case 2:
		return "<-chan"
	default:
		return "chan"
	}
}

func arrayLength(arr *ast.ArrayType) string {
	if arr.Len != nil {
		return "[" + arr.Len.(*ast.BasicLit).Value + "]"
	}
	return "[]"
}

func (td *Data) returnValue(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// the case of an identifier that's not a selector expression is matched by
		// a named type that belongs to the same package as the mocked interface

		switch t.Name {
		case "string":
			return `""`

		case "bool":
			return "false"

		case "error":
			return "nil"

		case
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"complex64", "complex128":
			return "0"

		case "float32", "float64":
			return "0.0"

		default:
			qname := td.qualifiedName(t)
			u, ok := td.Underlying[qname]
			if !ok {
				return qname + "{}"
			}
			// then consider the underlying type
			return td.returnValue(&ast.Ident{Name: u})
		}

	case *ast.SelectorExpr:
		tname := t.X.(*ast.Ident).Name + "." + t.Sel.Name
		u, ok := td.Underlying[tname]
		if ok {
			return td.returnValue(&ast.Ident{Name: u})
		}
		return tname + "{}"

	case
		*ast.StarExpr,
		*ast.FuncType,
		*ast.MapType,
		*ast.InterfaceType,
		*ast.ChanType:
		return "nil"

	case *ast.ArrayType:
		if t.Len != nil {
			return "[" + t.Len.(*ast.BasicLit).Value + "]" + td.expressionType(t.Elt) + "{}"
		}
		return "nil"

	case *ast.StructType:
		return "struct{}{}"

	default:
		return ""
	}
}

func (td *Data) expressionName(expr ast.Expr, name string) ParamName {
	_, isVararg := expr.(*ast.Ellipsis)
	return ParamName{name, isVararg}
}

func (td *Data) qualifiedName(ident *ast.Ident) string {
	if td.Package == "" {
		return ident.Name
	}
	if td.Qualify {
		return td.Package + "." + ident.Name
	}
	return ident.Name
}

func justNames(paramNames []ParamName) []string {
	ss := make([]string, 0, len(paramNames))
	for _, n := range paramNames {
		ss = append(ss, n.string)
	}
	return ss
}

func expandNames(paramNames []ParamName) []string {
	ss := make([]string, 0, len(paramNames))
	for _, n := range paramNames {
		ss = append(ss, n.Expand())
	}
	return ss
}
