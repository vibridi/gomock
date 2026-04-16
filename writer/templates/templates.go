package templates

const Options = `
type mock{{.ServiceName}}{{.TypeParamList}} struct {
	options mock{{.ServiceName}}Options{{.TypeArguments}}
}

type mock{{.ServiceName}}Options{{.TypeParamList}} struct {
	{{range .FuncDefs}}func{{.Name}}  func({{.Signature}}) {{.Return}}
	{{end}}
}

{{if eq .TypeParamList ""}}
var defaultMock{{.ServiceName}}Options = mock{{.ServiceName}}Options{
	{{range .FuncDefs}}func{{.Name}}: func({{.Signature}}) {{.Return}} {
		return {{.ReturnValues}}
	},
	{{end}}
}
{{else}}
func newDefaultMock{{.ServiceName}}Options{{.TypeParamList}}() mock{{.ServiceName}}Options{{.TypeArguments}} {
	return mock{{.ServiceName}}Options{{.TypeArguments}}{
		{{range .FuncDefs}}func{{.Name}}: func({{.Signature}}) {{.Return}} {
			return {{.ReturnValues}}
		},
		{{end}}
	}
}
{{end}}

type mock{{.ServiceName}}Option{{.TypeParamList}} func(*mock{{.ServiceName}}Options{{.TypeArguments}})

{{range .FuncDefs}}
func {{if $.Export}}W{{else}}w{{end}}ithFunc{{if $.Disambiguate}}{{.ServiceName}}{{.Name}}{{else}}{{.Name}}{{end}}{{$.TypeParamList}}(f func({{.Signature}}) {{.Return}}) mock{{.ServiceName}}Option{{$.TypeArguments}} {
	return func(o *mock{{.ServiceName}}Options{{$.TypeArguments}}) {
		o.func{{.Name}} = f
	}
}
{{end}}

{{range .FuncDefs}}
func (m *mock{{.ServiceName}}{{$.TypeArguments}}) {{.Name}}({{.Signature}}) {{.Return}} {
	return {{if .Return}}m.options.func{{.Name}}({{.Args}}){{end}}
}
{{end}}

func {{if .Export}}N{{else}}n{{end}}ewMock{{.ServiceName}}{{.TypeParamList}}(opt ...mock{{.ServiceName}}Option{{.TypeArguments}}) {{if .Qualify}}{{.Package}}.{{end}}{{if and .Qualify .PrefixPackage }}{{.InterfaceName}}{{else}}{{.ServiceName}}{{end}}{{.TypeArguments}} {
	opts := {{if eq .TypeParamList ""}}defaultMock{{.ServiceName}}Options{{else}}newDefaultMock{{.ServiceName}}Options{{.TypeArguments}}(){{end}}
	for _, o := range opt {
		o(&opts)
	}
	return &mock{{.ServiceName}}{{.TypeArguments}}{
		options: opts,
	}
}`

const Struct = `
type {{if .Export}}M{{else}}m{{end}}ock{{.ServiceName}}{{.TypeParamList}} struct {
	{{range .FuncDefs}}{{.Name}}Func  func({{.Signature}}) {{.Return}}
	{{end}}
}

{{range .FuncDefs}}
func (m *{{if $.Export}}M{{else}}m{{end}}ock{{.ServiceName}}{{$.TypeArguments}}) {{.Name}}({{.Signature}}) {{.Return}} {
	if m.{{.Name}}Func != nil {
		{{if .Return}}return m.{{.Name}}Func({{.Args}}){{else}}m.{{.Name}}Func({{.Args}}){{end}}
	}
	{{if .Return}}return {{.ReturnValues}}{{end -}}
}
{{end}}`
