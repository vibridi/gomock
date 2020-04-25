package templates

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "github.com/vibridi/gomock/parser"
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

func TestData(t *testing.T) {
	t.Run("write func type", func(t *testing.T) {
		cases := []string{
			"func()",                                      // no args
			"func(a string)",                              // one arg identifier
			"func(c complex128)",                          // one arg identifier
			"func(string)",                                // one arg no label
			"func(string, int)",                           // multiple args no label
			"func(a, b string)",                           // comma-delimited same-type args
			"func(a, b, c, d string, i int)",              // multiple comma-delimited same-type args
			"func(a string, i int, g test.TestStruct)",    // qualified name
			"func(a string, i int, g, h test.TestStruct)", // comma-delimited qualified names
			"func(f func())",                              // function
			"func(a string, f func())",                    // function as second arg
			"func(string, func())",                        // function arg with no label
			"func(f func(f func()))",                      // nested funcs
			"func(f func(f func(a string, b, c int, k func(z zap.Logger))))", // nested funcs with multiple args
			"func(a []string)",                             // slice
			"func(a [][]string)",                           // 2D slice
			"func(a [1]string)",                            // array
			"func(a, b [1]string)",                         // comma-delimited same-type arrays
			"func(a, b []func(a string))",                  // slice of funcs
			"func(a, b []func(a []string, b func([]int)))", // slice of nested funcs
			"func(a *string)",                              // pointer
			"func(a []*string)",                            // slice of pointers
			"func(a map[string]*string)",                   // map of pointers
			"func(a map[string]struct{})",                  // map of structs
			"func(i interface{})",                          // interface
			"func(a map[*zap.Logger]test.Test)",            // qualified pointer
			"func(ch chan int)",                            // chan
			"func(ch chan map[string]chan<- *zap.Logger)",  // directed chan of complex type
			"func(ss ...string)",                           // vararg
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

			td := &Data{
				Qualify: true,
			}

			s := td.functionType(ft)
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
			"Get(ss ...string)":                         "_",
			"Get(ff ...func(string, string))":           "_",
		}

		for method, expected := range cases {
			src := fmt.Sprintf(templMethod, method)

			f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.DeclarationErrors)
			assert.Nil(t, err)

			spec, err := gomock.GetInterfaceSpec(f, "")
			assert.Nil(t, err)

			m := spec.Type.(*ast.InterfaceType).Methods.List[0]

			td := &Data{}

			data := td.ToFuncDef(m)
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
			"Get() error":                  "nil",
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

			td := &Data{
				Qualify: true,
			}

			retval := td.returnValue(r.Type)
			assert.Equal(t, expected, retval)
		}
	})
}
