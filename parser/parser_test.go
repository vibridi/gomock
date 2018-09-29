package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/stretchr/testify/assert"

	//"go/ast"
	"testing"
)

const (
	templParamType = `
package test
type TestInterface interface { 
	Get(a %s) 
}`

	templMethod = `
package test
type TestInterface interface {
	%s
}
`
)

type TestInterface interface {
	Get() string
	Set(value string, marcus func(a string, b func(a string, b func(a int))))
}

func TestMethodDataToString(t *testing.T) {

	t.Run("vanilla method", func(t *testing.T) {
		m := Method{
			Name:    "Get",
			Params:  []*Param{},
			Returns: []string{},
		}
		assert.Equal(t, "Get()", m.String())
	})

	t.Run("nameless params", func(t *testing.T) {
		m := Method{
			Name: "Get",
			Params: []*Param{
				{
					Name:  "",
					Ptype: "string",
				},
				{
					Name:  "",
					Ptype: "func()",
				},
			},
			Returns: []string{},
		}
		assert.Equal(t, "Get(string, func())", m.String())
	})

	t.Run("named params and one return", func(t *testing.T) {
		m := Method{
			Name: "Get",
			Params: []*Param{
				{
					Name:  "a",
					Ptype: "string",
				},
				{
					Name:  "f",
					Ptype: "func()",
				},
			},
			Returns: []string{"int32"},
		}
		assert.Equal(t, "Get(a string, f func()) int32", m.String())
	})

	t.Run("named params and multiple returns", func(t *testing.T) {
		m := Method{
			Name: "Get",
			Params: []*Param{
				{
					Name:  "a",
					Ptype: "string",
				},
				{
					Name:  "f",
					Ptype: "func()",
				},
			},
			Returns: []string{"int32", "error"},
		}
		assert.Equal(t, "Get(a string, f func()) (int32, error)", m.String())
	})
}

func TestParser(t *testing.T) {

	t.Run("parse func type", func(t *testing.T) {
		cases := []string{
			"func()",
			"func(a string)",
			"func(string)",
			"func(string, int)",
			"func(a, b string)",
			"func(a, b, c, d string, i int)",
			"func(a string, i int, g Test)",
			"func(a string, i int, g test.Test)",
			"func(a string, i int, g, h test.Test)",
			"func(f func())",
			"func(a string, f func())",
			"func(string, func())",
			"func(f func(f func()))",
			"func(f func(f func(a string, b, c int, k func(z zap.Logger))))",
			"func(a []string)",
			"func(a [][]string)",
			"func(a [1]string)",
			"func(a, b [1]string)",
			"func(a, b []func(a string))",
			"func(a, b []func(a []string, b func([]int)))",
			"func(a *string)",
			"func(a []*string)",
		}

		for _, p := range cases {
			src := fmt.Sprintf(templParamType, p)

			f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
			assert.Nil(t, err)

			spec, err := getInterfaceSpec(f, "")
			assert.Nil(t, err)

			ft := spec.Type.(*ast.InterfaceType).
				Methods.List[0].Type.(*ast.FuncType).
				Params.List[0].Type.(*ast.FuncType)

			s := functionType(ft)
			assert.Equal(t, p, s)
		}
	})

	t.Run("parse method data", func(t *testing.T) {
		cases := map[string]string{
			"Get()":                                     "_",
			"Get(a, b string)":                          "Get(a string, b string)",
			"Get(string)":                               "_",
			"Get(string, string)":                       "_",
			"Get(string, string) string":                "_",
			"Get(a, b string, c, d string) string":      "Get(a string, b string, c string, d string) string",
			"Get(f func(), s string) string":            "_",
			"Get(f func(string, int), s string) string": "_",
			"Get(ss []string) (string)":                 "Get(ss []string) string",
			"Get(star *string) (string, error)":         "_",
			"Get() (string, int64, *zap.Logger, error)": "_",
		}

		for method, expected := range cases {
			src := fmt.Sprintf(templMethod, method)

			f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
			assert.Nil(t, err)

			spec, err := getInterfaceSpec(f, "")
			assert.Nil(t, err)

			m := spec.Type.(*ast.InterfaceType).Methods.List[0]

			data := toMethodData(m)
			if expected == "_" {
				assert.Equal(t, method, data.String())
			} else {
				assert.Equal(t, expected, data.String())
			}
		}

	})

	t.Run("what", func(t *testing.T) {
		t.Skip()
		p := "*zap.Logger"
		src := fmt.Sprintf(templParamType, p)

		f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
		assert.Nil(t, err)

		ast.Print(nil, f)
		//
		// PackageName := f.Name.Name
		// assert.Equal(t, "test", PackageName)

	})

}
