package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type MockData struct {
	PackageName           string
	InterfaceName         string
	MethodFields          []*ast.Field
	Components            []*ast.Field
	ExternalComponents    []*ast.Field
	InheritedMethodFields map[string][]*ast.Field
}

func (md *MockData) Len() int {
	return len(md.MethodFields) + len(md.Components) + len(md.ExternalComponents)
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
		return nil, errors.New("source interface declares no methods")
	}

	md.InterfaceName = spec.Name.Name

	for _, field := range interfaceType.Methods.List {
		switch field.Type.(type) {
		case *ast.FuncType:
			md.MethodFields = append(md.MethodFields, field)

		case *ast.Ident:
			md.Components = append(md.Components, field)

		case *ast.SelectorExpr:
			md.ExternalComponents = append(md.ExternalComponents, field)
		}
	}

	// If the interface contains any identifier, detect composition
	if len(md.Components) > 0 || len(md.ExternalComponents) > 0 {
		inheritedMethods, err := parseDirContent(md, filepath.Dir(srcFile))
		if err != nil {
			return nil, err
		}
		md.InheritedMethodFields = inheritedMethods
	}

	return md, nil
}

func GetInterfaceSpec(f *ast.File, target string) (*ast.TypeSpec, error) {
	interfaces, first := findInterfaces(f)

	var spec *ast.TypeSpec
	switch len(interfaces) {
	case 0:
		return nil, errors.New("source contains no type declarations")

	case 1:
		spec = interfaces[first]

	default:
		if target == "" {
			// When in doubt, work on the first interface
			spec = interfaces[first]
			break
		}
		for _, s := range interfaces {
			if s.Name.Name == target {
				spec = s
			}
		}
	}
	if spec == nil {
		return nil, errors.New("no suitable type declaration was found in source")
	}

	return spec, nil
}

func parseDirContent(md *MockData, srcDir string) (map[string][]*ast.Field, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), srcDir, nil, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	inheritedMethods := findInheritedMethodFields(md, pkgs)

	d, err := os.Open(srcDir)
	if err != nil {
		return nil, err
	}
	ff, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, f := range ff {
		if f.IsDir() {
			subPkgMethods, err := parseDirContent(md, srcDir+"/"+f.Name())
			if err != nil {
				return nil, err
			}
			for k, v := range subPkgMethods {
				inheritedMethods[k] = v
			}
		}
	}
	return inheritedMethods, nil
}

func findInheritedMethodFields(md *MockData, pkgs map[string]*ast.Package) map[string][]*ast.Field {
	inheritedMethodFields := make(map[string][]*ast.Field)

	if len(md.Components) > 0 {
		local := findLocalInheritedMethodFields(md.Components, pkgs[md.PackageName])
		for name, intf := range local {
			inheritedMethodFields[name] = intf.Type.(*ast.InterfaceType).Methods.List
		}
	}

	if len(md.ExternalComponents) > 0 {
		imported := findImportedInheritedMethodFields(md.ExternalComponents, pkgs)
		for name, intf := range imported {
			inheritedMethodFields[name] = intf.Type.(*ast.InterfaceType).Methods.List
		}
	}
	return inheritedMethodFields
}

func findLocalInheritedMethodFields(components []*ast.Field, pkg *ast.Package) map[string]*ast.TypeSpec {
	allInPkg := findAllInterfaces(pkg)
	filtered := make(map[string]*ast.TypeSpec, 0)

	// component fields are *ast.Ident
	for _, comp := range components {
		ident := comp.Type.(*ast.Ident).Name
		if intf, found := allInPkg[ident]; found {
			filtered[ident] = intf
		}
	}
	return filtered
}

func findImportedInheritedMethodFields(imported []*ast.Field, pkgs map[string]*ast.Package) map[string]*ast.TypeSpec {
	filtered := make(map[string]*ast.TypeSpec, 0)

	for _, imp := range imported {
		selexp := imp.Type.(*ast.SelectorExpr)
		pkgname := selexp.X.(*ast.Ident).Name
		ident := selexp.Sel.Name

		allInPkg := findAllInterfaces(pkgs[pkgname])

		if intf, found := allInPkg[ident]; found {
			filtered[ident] = intf
		}
	}
	return filtered
}

func findAllInterfaces(pkg *ast.Package) map[string]*ast.TypeSpec {
	if pkg == nil {
		return nil
	}
	all := make(map[string]*ast.TypeSpec, 0)
	for _, file := range pkg.Files {
		inFile, _ := findInterfaces(file)
		for name, spec := range inFile {
			all[name] = spec
		}
	}
	return all
}

func findInterfaces(file *ast.File) (interfaces map[string]*ast.TypeSpec, first string) {
	interfaces = make(map[string]*ast.TypeSpec, 0)

	for _, d := range file.Decls {
		if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			specs := gd.Specs[0].(*ast.TypeSpec)

			if _, ok := specs.Type.(*ast.InterfaceType); !ok {
				continue
			}
			interfaces[specs.Name.Name] = specs
			if first == "" {
				first = specs.Name.Name
			}
		}
	}
	return
}

func filterInterfaces(name string, file *ast.File) {

}
