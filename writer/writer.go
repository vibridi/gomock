package writer

import (
	"bytes"
	"go/ast"
	"strconv"
	"strings"
	"text/template"

	"github.com/vibridi/gomock/helper"
	"github.com/vibridi/gomock/parser"
)

const (
	mockTemplate = `
type mock{{.ServiceName}} struct {
	options mock{{.ServiceName}}Options
}

type mock{{.ServiceName}}Options struct {
	{{range .FuncDefs}}func{{.Name}}  func({{.Signature}}) {{.Return}}
	{{end}}
}

var defaultMock{{.ServiceName}}Options = mock{{.ServiceName}}Options{
	{{range .FuncDefs}}func{{.Name}}: func({{.SignatureUnnamed}}) {{.Return}} {
		return {{.ReturnValues}}
	},
	{{end}}
}

type mock{{.ServiceName}}Option func(*mock{{.ServiceName}}Options)

{{range .FuncDefs}}
func {{if $.Export}}W{{else}}w{{end}}ithFunc{{.Name}}(f func({{.SignatureUnnamed}}) {{.Return}}) mock{{.ServiceName}}Option {
	return func(o *mock{{.ServiceName}}Options) {
		o.func{{.Name}} = f
	}
}
{{end}}

{{range .FuncDefs}}
func (m *mock{{.ServiceName}}) {{.Name}}({{.Signature}}) {{.Return}} {
	return {{if .Return}}m.options.func{{.Name}}({{.Args}}){{end}}
}
{{end}}

func {{if .Export}}N{{else}}n{{end}}ewMock{{.ServiceName}}(opt ...mock{{.ServiceName}}Option) {{if .Qualify}}{{.Package}}.{{end}}{{.ServiceName}} {
	opts := defaultMock{{.ServiceName}}Options
	for _, o := range opt {
		o(&opts)
	}
	return &mock{{.ServiceName}}{
		options: opts,
	}
}`
)

type TemplateData struct {
	Qualify     bool
	Export      bool
	Package     string
	ServiceName string
	FuncDefs    []*FuncDef
}

type FuncDef struct {
	ServiceName      string
	Name             string
	Signature        string
	SignatureUnnamed string
	Return           string
	Args             string
	ReturnValues     string
}

func (fd FuncDef) String() string {
	s := fd.Name + "(" + fd.Signature + ") " + fd.Return
	return strings.TrimSpace(s)
}

type ParamName struct {
	string
	IsVararg bool
}

func (pn ParamName) Expand() string {
	if pn.IsVararg {
		return pn.string + "..."
	}
	return pn.string
}

func Write(data *parser.MockData, qualify bool, export bool) (string, error) {

	if len(data.MethodFields) == 0 {
		return "", nil
	}

	d := toTemplateData(data, qualify, export)

	var buf bytes.Buffer
	t := template.Must(template.New("mock").Parse(mockTemplate))
	if err := t.Execute(&buf, d); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toTemplateData(data *parser.MockData, qualify bool, export bool) *TemplateData {
	d := &TemplateData{}
	d.Qualify = qualify
	d.Export = export
	d.Package = data.PackageName
	d.ServiceName = data.InterfaceName

	funcDefs := make([]*FuncDef, 0, len(data.MethodFields))
	for _, field := range data.MethodFields {
		funcDefs = append(funcDefs, d.toFuncDef(field, qualify))
	}

	d.FuncDefs = funcDefs
	return d
}

func (td *TemplateData) toFuncDef(field *ast.Field, qualify bool) *FuncDef {

	fn := field.Type.(*ast.FuncType)

	funcDef := &FuncDef{}
	funcDef.ServiceName = td.ServiceName
	funcDef.Name = field.Names[0].Name

	paramNames := make([]ParamName, 0, len(fn.Params.List))
	paramTypes := make([]string, 0, len(fn.Params.List))

	for i, p := range fn.Params.List {
		if len(p.Names) == 0 {
			paramNames = append(paramNames, td.expressionName(p.Type, "p"+strconv.Itoa(i)))
			paramTypes = append(paramTypes, td.expressionType(p.Type, qualify))

		} else {
			for _, n := range p.Names {
				paramNames = append(paramNames, td.expressionName(p.Type, n.Name))
				paramTypes = append(paramTypes, td.expressionType(p.Type, qualify))
			}
		}
	}

	funcDef.Signature = strings.Join(helper.Zips(justNames(paramNames), paramTypes, " "), ", ")
	funcDef.SignatureUnnamed = strings.Join(paramTypes, ", ")
	funcDef.Args = strings.Join(expandNames(paramNames), ", ")

	if fn.Results == nil {
		return funcDef
	}

	returnTypes := make([]string, 0, len(fn.Results.List))
	returnValues := make([]string, 0, len(fn.Results.List))

	for _, r := range fn.Results.List {
		returnTypes = append(returnTypes, td.expressionType(r.Type, qualify))
		returnValues = append(returnValues, td.returnValue(r.Type, qualify))
	}

	funcDef.Return = helper.ReturnTypesToString(returnTypes)
	funcDef.ReturnValues = strings.Join(returnValues, ", ")

	return funcDef
}

func (td *TemplateData) expressionType(expr ast.Expr, qualify bool) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if qualify && ast.IsExported(t.Name) {
			return td.Package + "." + t.Name
		}
		return t.Name

	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name

	case *ast.FuncType:
		return td.functionType(t, qualify)

	case *ast.ArrayType:
		return arrayLength(t) + td.expressionType(t.Elt, qualify)

	case *ast.StarExpr:
		return "*" + td.expressionType(t.X, qualify)

	case *ast.MapType:
		return "map[" + td.expressionType(t.Key, qualify) + "]" + td.expressionType(t.Value, qualify)

	case *ast.StructType:
		return "struct{}"

	case *ast.InterfaceType:
		return "interface{}"

	case *ast.ChanType:
		return chanType(t) + " " + td.expressionType(t.Value, qualify)

	case *ast.Ellipsis:
		return "..." + td.expressionType(t.Elt, qualify)

	default:
		return ""
	}
}

func (td *TemplateData) functionType(fn *ast.FuncType, qualify bool) string {
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
		p += td.expressionType(param.Type, qualify)
		pdecl = append(pdecl, p)
	}
	s += strings.Join(pdecl, ", ")
	s += ")"

	rs := make([]string, 0)
	if fn.Results != nil {
		for _, r := range fn.Results.List {
			rs = append(rs, td.expressionType(r.Type, qualify))
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

func (td *TemplateData) returnValue(expr ast.Expr, qualify bool) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if t.Obj != nil {
			if qualify {
				return td.Package + "." + t.Name + "{}"
			}
			return t.Name + "{}"
		}

		switch t.Name {
		case "string":
			return "\"\""

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
			return t.Name + "{}"
		}

	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name + "{}"

	case
		*ast.StarExpr,
		*ast.FuncType,
		*ast.MapType,
		*ast.InterfaceType,
		*ast.ChanType:
		return "nil"

	case *ast.ArrayType:
		if t.Len != nil {
			return "[" + t.Len.(*ast.BasicLit).Value + "]" + td.expressionType(t.Elt, qualify) + "{}"
		}
		return "nil"

	case *ast.StructType:
		return "struct{}{}"

	default:
		return ""
	}
}

func (td *TemplateData) expressionName(expr ast.Expr, name string) ParamName {
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
