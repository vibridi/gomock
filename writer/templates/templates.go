package templates

const Options = `
type mock{{.ServiceName}} struct {
	options mock{{.ServiceName}}Options
}

type mock{{.ServiceName}}Options struct {
	{{range .FuncDefs}}func{{.Name}}  func({{.Signature}}) {{.Return}}
	{{end}}
}

var defaultMock{{.ServiceName}}Options = mock{{.ServiceName}}Options{
	{{range .FuncDefs}}func{{.Name}}: func({{.Signature}}) {{.Return}} {
		return {{.ReturnValues}}
	},
	{{end}}
}

type mock{{.ServiceName}}Option func(*mock{{.ServiceName}}Options)

{{range .FuncDefs}}
func {{if $.Export}}W{{else}}w{{end}}ithFunc{{if $.Disambiguate}}{{.ServiceName}}{{.Name}}{{else}}{{.Name}}{{end}}(f func({{.Signature}}) {{.Return}}) mock{{.ServiceName}}Option {
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

const Struct = `
type {{if .Export}}M{{else}}m{{end}}ock{{.ServiceName}} struct {
	{{range .FuncDefs}}{{.Name}}Func  func({{.Signature}}) {{.Return}}
	{{end}}
}

{{range .FuncDefs}}
func (m *{{if $.Export}}M{{else}}m{{end}}ock{{.ServiceName}}) {{.Name}}({{.Signature}}) {{.Return}} {
	if m.{{.Name}}Func != nil {
		{{if .Return}}return m.{{.Name}}Func({{.Args}}){{else}}m.{{.Name}}Func({{.Args}}){{end}}
	}
	{{if .Return}}return {{.ReturnValues}}{{end}}
}
{{end}}`
