package writer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	gomock "gomock/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	templParamType = `
package test
type TestInterface interface { 
	Get(a %s) 
}
type TestStruct struct {
}
`

	templMethod = `
package test
type TestInterface interface {
	%s
}
`
)

func TestParser(t *testing.T) {

	t.Run("write func type", func(t *testing.T) {
		cases := []string{
			"func()",
			"func(a string)",
			"func(c complex128)",
			"func(string)",
			"func(string, int)",
			"func(a, b string)",
			"func(a, b, c, d string, i int)",
			"func(a string, i int, g test.TestStruct)",
			"func(a string, i int, g test.TestStruct)",
			"func(a string, i int, g, h test.TestStruct)",
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
			"func(a map[string]*string)",
			"func(a map[string]struct{})",
			"func(i interface{})",
			"func(a map[*zap.Logger]test.Test)",
			"func(ch chan int)",
			"func(ch chan map[string]chan<- *zap.Logger)",
		}

		for _, p := range cases {
			src := fmt.Sprintf(templParamType, p)

			f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
			assert.Nil(t, err)

			spec, err := gomock.GetInterfaceSpec(f, "")
			assert.Nil(t, err)

			ft := spec.Type.(*ast.InterfaceType).
				Methods.List[0].Type.(*ast.FuncType).
				Params.List[0].Type.(*ast.FuncType)

			td := &TemplateData{}

			s := td.functionType(ft, true)
			assert.Equal(t, p, s)
		}
	})

	t.Run("write method", func(t *testing.T) {
		cases := map[string]string{
			"Get()":                                     "_",
			"Get(a, b string)":                          "Get(a string, b string)",
			"Get(string)":                               "Get(p0 string)",
			"Get(string, string)":                       "Get(p0 string, p1 string)",
			"Get(string, string) string":                "Get(p0 string, p1 string) string",
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

			spec, err := gomock.GetInterfaceSpec(f, "")
			assert.Nil(t, err)

			m := spec.Type.(*ast.InterfaceType).Methods.List[0]

			td := TemplateData{}

			data := td.toFuncDef(m, false)
			if expected == "_" {
				assert.Equal(t, method, data.String())
			} else {
				assert.Equal(t, expected, data.String())
			}
		}
	})

	t.Run("write return values", func(t *testing.T) {
		cases := map[string]string{
			"Get() string":                 "\"\"",
			"Get() int":                    "0",
			"Get() int8":                   "0",
			"Get() int16":                  "0",
			"Get() int32":                  "0",
			"Get() int64":                  "0",
			"Get() map[string]int":         "nil",
			"Get() Test":                   "Test{}",
			"Get() test.Test":              "test.Test{}",
			"Get() *test.Test":             "nil",
			"Get() []test.Test":            "nil",
			"Get() struct{}":               "struct{}{}",
			"Get() interface{}":            "nil",
			"Get() <-chan amqp.Delivery":   "nil",
			"Get() [2]string":              "[2]string{}",
			"Get() [2]map[string]struct{}": "[2]map[string]struct{}{}",
		}

		for method, expected := range cases {
			src := fmt.Sprintf(templMethod, method)

			f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
			assert.Nil(t, err)

			spec, err := gomock.GetInterfaceSpec(f, "")
			assert.Nil(t, err)

			m := spec.Type.(*ast.InterfaceType).Methods.List[0]
			r := m.Type.(*ast.FuncType).Results.List[0]

			td := TemplateData{}

			retval := td.returnValue(r.Type, true)
			assert.Equal(t, expected, retval)
		}
	})

	t.Run("write mock template", func(t *testing.T) {

		src := fmt.Sprintf(templMethod, "Get() string")

		md, err := gomock.Parse("", src, "")
		assert.Nil(t, err)

		out, err := Write(md, false)
		assert.Nil(t, err)
		assert.NotEqual(t, "", out)
	})

}
