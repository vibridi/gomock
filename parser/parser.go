package parser

import (
	"go/ast"
	"go/parser"
	"go/token"

	throws "github.com/vibridi/gomock/error"
)

type MockData struct {
	PackageName   string
	InterfaceName string
	MethodFields  []*ast.Field
}

func Parse(srcFile string, src interface{}, target string) (*MockData, error) {
	f, err := parser.ParseFile(token.NewFileSet(), srcFile, src, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	md := &MockData{}
	md.PackageName = f.Name.Name

	spec, err := GetInterfaceSpec(f, target)
	if err != nil {
		return nil, err
	}
	interfaceType := spec.Type.(*ast.InterfaceType)

	if interfaceType.Incomplete {
		return nil, throws.NoMethods
	}

	md.InterfaceName = spec.Name.Name

	for _, field := range interfaceType.Methods.List {
		md.MethodFields = append(md.MethodFields, field)
	}

	return md, nil
}

func GetInterfaceSpec(f *ast.File, target string) (*ast.TypeSpec, error) {
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
