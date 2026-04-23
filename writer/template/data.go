package template

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/vibridi/gomock/v3/internal/fn"
)

// Holds the data needed to execute the mock template.
type data struct {
	Qualify       bool
	Export        bool
	Disambiguate  bool
	Package       string
	ServiceName   string
	InterfaceName string
	FuncDefs      []*funcDef
	UnnamedSig    bool
	Underlying    map[string]string
	Aliases       map[string]string
	PrefixPackage bool
	TypeParamList string // full type parameter list as it appears in the interface declaration
	TypeArguments string // type argument list as it appears in the method receiver

	// unexported
	typeParamSet map[string]struct{}
}

// Populates TypeParamList and TypeArguments from the given list of type parameters.
func (td *data) AddTypeParameters(typeParams []*ast.Field) {
	if len(typeParams) == 0 {
		return
	}

	typeParamNames := make([]string, 0, len(typeParams))
	typeParamDefs := make([]string, 0, len(typeParams))

	for _, tp := range typeParams {
		for _, n := range tp.Names {
			typeParamNames = append(typeParamNames, n.Name)
			typeParamDefs = append(typeParamDefs, n.Name+" "+td.expressionType(tp.Type))
		}
	}
	td.TypeArguments = "[" + strings.Join(typeParamNames, ",") + "]"
	td.TypeParamList = "[" + strings.Join(typeParamDefs, ", ") + "]"

	td.typeParamSet = make(map[string]struct{})
	for _, n := range typeParamNames {
		td.typeParamSet[n] = struct{}{}
	}
}

func (td *data) AppendFuncDef(field *ast.Field) {
	ftype, ok := field.Type.(*ast.FuncType)
	if !ok {
		return
	}

	funcDef := &funcDef{}
	funcDef.ServiceName = td.ServiceName
	funcDef.Name = field.Names[0].Name

	paramNames := make([]ParamName, 0, len(ftype.Params.List))
	paramTypes := make([]string, 0, len(ftype.Params.List))

	for i, p := range ftype.Params.List {
		if len(p.Names) == 0 {
			paramNames = append(paramNames, paramName(p.Type, "p"+strconv.Itoa(i)))
			paramTypes = append(paramTypes, td.expressionType(p.Type))

		} else {
			for _, n := range p.Names {
				paramNames = append(paramNames, paramName(p.Type, n.Name))
				paramTypes = append(paramTypes, td.expressionType(p.Type))
			}
		}
	}

	if !td.UnnamedSig {
		funcDef.Signature = strings.Join(fn.Zips(justNames(paramNames), paramTypes, " "), ", ")
	} else {
		funcDef.Signature = strings.Join(paramTypes, ", ")
	}

	funcDef.Args = strings.Join(expandNames(paramNames), ", ")

	if ftype.Results == nil {
		td.FuncDefs = append(td.FuncDefs, funcDef)
		return
	}

	returnTypes := make([]string, 0, len(ftype.Results.List))
	returnValues := make([]string, 0, len(ftype.Results.List))

	for _, r := range ftype.Results.List {
		returnTypes = append(returnTypes, td.expressionType(r.Type))
		returnValues = append(returnValues, td.returnValue(r.Type))
	}

	funcDef.Return = formatReturnTypes(returnTypes)
	funcDef.ReturnValues = strings.Join(returnValues, ", ")

	td.FuncDefs = append(td.FuncDefs, funcDef)
}

func (td *data) expressionType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if td.Qualify && ast.IsExported(t.Name) && !td.isTypeParam(t) {
			return td.Package + "." + t.Name
		}
		return t.Name

	case *ast.SelectorExpr:
		pkg := t.X.(*ast.Ident).Name
		if alias, ok := td.Aliases[pkg]; ok {
			pkg = alias
		}
		return pkg + "." + t.Sel.Name

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

	case *ast.UnaryExpr:
		switch t.Op {
		case token.TILDE:
			return "~" + td.expressionType(t.X)
		default:
			return td.expressionType(t.X)
		}

	default:
		return ""
	}
}

func (td *data) functionType(fn *ast.FuncType) string {
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
	s += strings.Join(rs, ", ")
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

func (td *data) returnValue(expr ast.Expr) string {
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
			if td.isTypeParam(t) {
				return "*new(" + t.Name + ")"
			}
			qname := td.qualifiedName(t)
			u, ok := td.Underlying[qname]
			if !ok {
				return qname + "{}"
			}
			// then consider the underlying type
			return td.returnValue(&ast.Ident{Name: u})
		}

	case *ast.SelectorExpr:
		pkg := t.X.(*ast.Ident).Name
		if alias, ok := td.Aliases[pkg]; ok {
			pkg = alias
		}
		tname := pkg + "." + t.Sel.Name
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

func (td *data) qualifiedName(ident *ast.Ident) string {
	if td.Package == "" {
		return ident.Name
	}
	if td.Qualify {
		return td.Package + "." + ident.Name
	}
	return ident.Name
}

func (td *data) isTypeParam(t *ast.Ident) bool {
	_, ok := td.typeParamSet[t.Name]
	return ok
}

func paramName(expr ast.Expr, name string) ParamName {
	_, isVararg := expr.(*ast.Ellipsis)
	return ParamName{name, isVararg}
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

func formatReturnTypes(r []string) string {
	switch len(r) {
	case 0:
		return ""
	case 1:
		return r[0]
	default:
		return "(" + strings.Join(r, ", ") + ")"
	}
}
