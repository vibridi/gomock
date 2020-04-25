package writer

import (
	"bytes"
	"go/ast"
	"text/template"

	"github.com/vibridi/gomock/writer/templates"

	"github.com/vibridi/gomock/parser"
)

type WriteOpts struct {
	Qualify          bool
	Export           bool
	UnnamedSignature bool
	StructStyle      bool
}

type writer struct {
	data *parser.MockData
	opts WriteOpts
}

func New(data *parser.MockData, opts WriteOpts) *writer {
	return &writer{data, opts}
}

func (w *writer) Write() (string, error) {
	if w.data.Len() == 0 {
		return "", nil
	}

	d := w.buildTemplateData()

	mockTemplate := templates.Options
	if w.opts.StructStyle {
		mockTemplate = templates.Struct
	}

	var buf bytes.Buffer
	t := template.Must(template.New("mock").Parse(mockTemplate))
	if err := t.Execute(&buf, d); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (w *writer) buildTemplateData() *templates.Data {
	d := &templates.Data{
		Qualify:     w.opts.Qualify,
		Export:      w.opts.Export,
		Package:     w.data.PackageName,
		ServiceName: w.data.InterfaceName,
		UnnamedSig:  w.opts.UnnamedSignature,
	}

	funcDefs := make([]*templates.FuncDef, 0, len(w.data.MethodFields))

	for _, field := range w.data.MethodFields {
		funcDefs = append(funcDefs, d.ToFuncDef(field))
	}

	for _, field := range w.data.Components {
		local := w.data.InheritedMethodFields[field.Type.(*ast.Ident).Name]
		for _, lm := range local {
			funcDefs = append(funcDefs, d.ToFuncDef(lm))
		}
	}

	for _, field := range w.data.ExternalComponents {
		imported := w.data.InheritedMethodFields[field.Type.(*ast.SelectorExpr).Sel.Name]
		for _, im := range imported {
			funcDefs = append(funcDefs, d.ToFuncDef(im))
		}
	}

	d.FuncDefs = funcDefs
	return d
}
