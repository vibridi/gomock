package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	regexp "regexp"
	"strings"
)

const (

	mockTemplate = `
type mock@@service@@ struct {
	options mock@@service@@Options
}

type mock@@service@@Options struct {
	@@funcdefs@@
}

var defaultMock@@service@@Options = mock@@service@@Options{
	@@funcdefaults@@
}

type mock@@service@@Option func(*mock@@service@@Options)

@@funcwiths@@

@@funcimpls@@

func newMock@@service@@(opt ...mock@@service@@Option) @@service@@ {
	opts := defaultMock@@service@@Options
	for _, o := range opt {
		o(&opts)
	}
	return &mock@@service@@{
		options: opts,
	}
}`


	mockFunctionDefinition = `func@@method@@  func(@@signature@@) @@return@@`

	mockFunctionDefault = `
func@@method@@: func(@@signature@@) @@return@@ {
	return @@return_clause@@
},`


	mockWithFunc = `
func withFunc@@method@@(f func(@@signature@@) @@return@@) mock@@service@@Option {
	return func(o *mock@@service@@Options) {
		o.func@@method@@ = f
	}
}`

	mockImpl = `
func (m *mock@@service@@) @@method@@(@@signature@@) @@return@@ {
	return m.options.func@@method@@()
}`

)


func main() {
	app := cli.NewApp()

	app.Name = "gomock"

	app.Action = func(c *cli.Context) error {

		var re *regexp.Regexp
		intf := strings.TrimSpace(c.Args().Get(0))

		re = regexp.MustCompile(".*type\\s+(.*)\\s+interface")
		service := re.FindStringSubmatch(intf)[1]

		re = regexp.MustCompile("\\{((.|\n)*)\\}")
		decl := re.FindStringSubmatch(intf)[1]

		funcs := make([]string, 0, 2)
		for _, s := range strings.Split(decl, "\n") {
			if len(s) == 0 {
				continue
			}
			s = strings.TrimSpace(s)
			funcs = append(funcs, s)
		}

		funcdefs := make([]string, 0, 2)
		funcdefaults := make([]string, 0, 2)
		funcwiths := make([]string, 0, 2)
		funcimpls := make([]string, 0, 2)

		re = regexp.MustCompile("(.*?)\\((.*?)\\)\\s*(.*)\\s*")
		for _, s := range funcs {
			subs := re.FindStringSubmatch(s)
			mName := strings.TrimSpace(subs[1])
			mSign := subs[2]
			mRet := ""
			if len(subs) > 3 {
				mRet = subs[3]
			}

			var fdef string
			fdef = strings.Replace(mockFunctionDefinition, "@@method@@", mName, -1)
			fdef = strings.Replace(fdef, "@@signature@@", mSign, -1)
			fdef = strings.Replace(fdef, "@@return@@", mRet, -1)
			funcdefs = append(funcdefs, fdef)

			var dft string
			dft = strings.Replace(mockFunctionDefault, "@@method@@", mName, -1)
			dft = strings.Replace(dft, "@@signature@@", mSign, -1)
			dft = strings.Replace(dft, "@@return@@", mRet, -1)

			rr := make([]string, 0, len(strings.Split(mRet, ",")))
			for i := 0; i < cap(rr); i++ {
				rr = append(rr, "nil")
			}

			dft = strings.Replace(dft, "@@return_clause@@", strings.Join(rr, ","), -1)
			funcdefaults = append(funcdefaults, dft)

			var withf string
			withf = strings.Replace(mockWithFunc, "@@method@@", mName, -1)
			withf = strings.Replace(withf, "@@signature@@", mSign, -1)
			withf = strings.Replace(withf, "@@return@@", mRet, -1)
			withf = strings.Replace(withf, "@@service@@", service, -1)
			funcwiths = append(funcwiths, withf)

			var impl string
			impl = strings.Replace(mockImpl, "@@method@@", mName, -1)
			impl = strings.Replace(impl, "@@signature@@", mSign, -1)
			impl = strings.Replace(impl, "@@return@@", mRet, -1)
			impl = strings.Replace(impl, "@@service@@", service, -1)
			funcimpls = append(funcimpls, impl)

		}

		var out string
		out = strings.Replace(mockTemplate, "@@service@@", service, -1)
		out = strings.Replace(out, "@@funcdefs@@", strings.Join(funcdefs, "\n"), -1)
		out = strings.Replace(out, "@@funcdefaults@@", strings.Join(funcdefaults, "\n"), -1)
		out = strings.Replace(out, "@@funcwiths@@", strings.Join(funcwiths, "\n"), -1)
		out = strings.Replace(out, "@@funcimpls@@", strings.Join(funcimpls, "\n"), -1)

		fmt.Println(out)

		return nil
	}

	app.Run(os.Args)
}
