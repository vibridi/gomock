package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	throws "gomock/error"
	"strings"
)

type MockData struct {
	PackageName   string
	InterfaceName string
	Methods       []*Method
}

type Method struct {
	Name    string
	Params  []*Param
	Returns []string
}

type Param struct {
	Name  string
	Ptype string
}

func (m *Method) String() string {
	s := m.Name
	s += "("
	pp := make([]string, 0)
	for _, v := range m.Params {
		p := ""
		if v.Name != "" {
			p += v.Name + " "
		}
		p += v.Ptype
		pp = append(pp, p)
	}
	s += strings.Join(pp, ", ")
	s += ")"

	switch len(m.Returns) {
	case 0:
		return s
	case 1:
		return s + " " + m.Returns[0]
	default:
		return s + " (" + strings.Join(m.Returns, ", ") + ")"
	}
}

func Parse(src string, target string) (*ast.File, error) {

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	data := &MockData{}
	data.PackageName = f.Name.Name

	spec, err := getInterfaceSpec(f, target)
	if err != nil {
		return nil, err
	}
	interfaceType := spec.Type.(*ast.InterfaceType)

	if interfaceType.Incomplete {
		return nil, throws.NoMethods
	}

	data.InterfaceName = spec.Name.Name

	methods := make([]*Method, 0)
	for _, field := range interfaceType.Methods.List {
		methods = append(methods, toMethodData(field))
	}
	data.Methods = methods

	return f, nil
}

func getInterfaceSpec(f *ast.File, target string) (*ast.TypeSpec, error) {
	typeSpecs := make([]*ast.TypeSpec, 0)

	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			specs := gd.Specs[0].(*ast.TypeSpec)

			if _, ok := specs.Type.(*ast.InterfaceType); !ok {
				continue
			}

			typeSpecs = append(typeSpecs, specs)
		}
	}

	var spec *ast.TypeSpec
	switch len(typeSpecs) {
	case 0:
		return nil, throws.NoTypeDeclarations

	case 1:
		spec = typeSpecs[0]

	default:
		if target == "" {
			return nil, throws.InterfaceNotSpecified
		}
		for _, s := range typeSpecs {
			if s.Name.Name == target {
				spec = s
			}
		}
	}
	if spec == nil {
		return nil, throws.InterfaceNotFound
	}

	return spec, nil
}

func toMethodData(field *ast.Field) *Method {
	m := &Method{}
	if len(field.Names) > 0 {
		m.Name = field.Names[0].Name
	}
	m.Params = make([]*Param, 0)
	m.Returns = make([]string, 0)

	fn := field.Type.(*ast.FuncType)
	for _, p := range fn.Params.List {

		if len(p.Names) == 0 {
			par := &Param{
				Ptype: expressionType(p.Type),
			}
			m.Params = append(m.Params, par)

		} else {
			for _, n := range p.Names {
				par := &Param{
					Name:  n.Name,
					Ptype: expressionType(p.Type),
				}
				m.Params = append(m.Params, par)
			}
		}
	}

	if fn.Results != nil {
		for _, r := range fn.Results.List {
			m.Returns = append(m.Returns, expressionType(r.Type))
		}
	}

	return m
}

func functionType(fn *ast.FuncType) string {
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
		p += expressionType(param.Type)
		pdecl = append(pdecl, p)
	}
	s += strings.Join(pdecl, ", ")
	s += ")"

	rs := make([]string, 0)
	if fn.Results != nil {
		for _, r := range fn.Results.List {
			rs = append(rs, expressionType(r.Type))
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

func expressionType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name

	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name

	case *ast.FuncType:
		return functionType(t)

	case *ast.ArrayType:
		return arrayLength(t) + expressionType(t.Elt)

	case *ast.StarExpr:
		return "*" + expressionType(t.X)

	default:
		return ""
	}
}

func arrayLength(arr *ast.ArrayType) string {
	if arr.Len != nil {
		return "[" + arr.Len.(*ast.BasicLit).Value + "]"
	}
	return "[]"
}
