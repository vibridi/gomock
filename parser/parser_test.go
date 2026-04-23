package parser

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ()

func TestParser(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		const src = `
package test
type TestInterface interfac
`
		md, err := Parse("", src, "")
		assert.Nil(t, md)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("no interface", func(t *testing.T) {
		const src = `
package test
type MyInt uint8
`
		md, err := Parse("", src, "")
		assert.Nil(t, md)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("parse", func(t *testing.T) {
		const src = `
package test
type foo interface {
	Do() error
}
type TestInterface interface { 
	foo
	bar.B
	Get(a string) 
}`
		md, err := Parse("", src, "TestInterface")
		assert.NotNil(t, md)
		assert.Nil(t, err)
		assert.Equal(t, "test", md.PackageName)
		assert.Equal(t, "TestInterface", md.InterfaceName)
		assert.Len(t, md.MethodFields, 1)
		assert.Equal(t, 3, md.Len())
	})

	t.Run("default to first", func(t *testing.T) {
		const src = `
package test
type Foo interface { 
	Get() 
}
type Bar interface {
	Do() error
}
`
		md, err := Parse("", src, "")
		assert.NotNil(t, md)
		assert.Nil(t, err)
		assert.Equal(t, "test", md.PackageName)
		assert.Equal(t, "Foo", md.InterfaceName)
	})

	t.Run("target not found", func(t *testing.T) {
		const src = `
package test
type Foo interface { 
	Get() 
}
type Bar interface {
	Do() error
}
`
		md, err := Parse("", src, "Baz")
		assert.Nil(t, md)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(), "target not found")
	})

	t.Run("can parse generic interface", func(t *testing.T) {
		const src = `
package test
type TestInterface[T any] interface { 
	Get(a T) 
}`
		md, err := Parse("", src, "")
		assert.NotNil(t, md)
		assert.Nil(t, err)
		assert.Equal(t, "test", md.PackageName)
		assert.Equal(t, "TestInterface", md.InterfaceName)
		assert.Len(t, md.MethodFields, 1)
	})

	t.Run("find all interfaces", func(t *testing.T) {
		pkg := &ast.Package{
			Name:    "",
			Scope:   nil,
			Imports: nil,
			Files: map[string]*ast.File{
				"foo.go": {
					Decls: []ast.Decl{
						&ast.GenDecl{
							Tok: token.TYPE,
							Specs: []ast.Spec{
								&ast.TypeSpec{
									Name: &ast.Ident{Name: "foo"},
									Type: &ast.InterfaceType{},
								},
								&ast.TypeSpec{
									Name: &ast.Ident{Name: "bar"},
									Type: &ast.StructType{},
								},
							},
						},
					},
				},
			},
		}
		m := findAllInterfaces(pkg)
		assert.Len(t, m, 1)
		assert.IsType(t, &ast.TypeSpec{}, m["foo"])
	})
}

func TestParseDirContent(t *testing.T) {
	tmp := t.TempDir()

	// Root package
	rootFile := `
	package rootpkg
	
	type RootInterface interface {
		RootMethod()
	}
	`
	err := os.WriteFile(filepath.Join(tmp, "root.go"), []byte(rootFile), 0644)
	require.Nil(t, err)

	// Subdirectory package
	subDir := filepath.Join(tmp, "subpkg")
	err = os.Mkdir(subDir, 0755)
	require.Nil(t, err)

	subFile := `
	package subpkg
	
	type SubInterface interface {
		SubMethod()
	}
	`
	err = os.WriteFile(filepath.Join(subDir, "sub.go"), []byte(subFile), 0644)
	require.Nil(t, err)

	// MockData configured to pick up both interfaces
	md := &MockData{
		PackageName: "rootpkg",
		Components: []*ast.Field{
			{
				Type: &ast.Ident{Name: "RootInterface"},
			},
		},
		ExternalComponents: []*ast.Field{
			{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "subpkg"},
					Sel: &ast.Ident{Name: "SubInterface"},
				},
			},
		},
	}

	result, err := parseDirContent(md, tmp)
	require.Nil(t, err)
	assert.Equal(t, "RootMethod", result["RootInterface"][0].Names[0].Name)
	assert.Equal(t, "SubMethod", result["SubInterface"][0].Names[0].Name)
}
