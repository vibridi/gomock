package template

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"strings"
	"text/template"
	"unicode"

	"github.com/vibridi/gomock/v3/parser"
)

type Opts struct {
	Qualify          bool
	Export           bool
	UnnamedSignature bool
	StructStyle      bool
	Disambiguate     bool
	MockName         string
	Underlying       []string
	ImportAliases    []string
	PrefixPackage    bool
}

// Executes the template according to the write options
func Exec(mock *parser.MockData, opts Opts) ([]byte, error) {
	if mock.Len() == 0 {
		return nil, nil
	}

	d, err := buildData(mock, opts)
	if err != nil {
		return nil, err
	}

	mockTemplate := Options
	if opts.StructStyle {
		mockTemplate = Struct
	}

	var buf bytes.Buffer
	t := template.Must(template.New("mock").Parse(mockTemplate))
	if err := t.Execute(&buf, d); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func buildData(mock *parser.MockData, opts Opts) (*data, error) {
	d := &data{
		Qualify:       opts.Qualify,
		Export:        opts.Export,
		Disambiguate:  opts.Disambiguate,
		Package:       mock.PackageName,
		ServiceName:   mock.InterfaceName,
		InterfaceName: mock.InterfaceName,
		UnnamedSig:    opts.UnnamedSignature,
		Underlying:    make(map[string]string, len(opts.Underlying)),
		Aliases:       make(map[string]string, len(opts.ImportAliases)),
		PrefixPackage: opts.PrefixPackage,
		// computed
		FuncDefs:      nil,
		TypeArguments: "",
		TypeParamList: "",
	}
	// Override the service name with the one supplied by the user, if any
	if opts.MockName != "" {
		d.ServiceName = opts.MockName
	}
	if opts.PrefixPackage {
		r := []rune(mock.PackageName)
		r[0] = unicode.ToUpper(r[0])
		prefix := string(r)
		d.ServiceName = prefix + mock.InterfaceName
	}

	for _, utype := range opts.Underlying {
		t, u, ok := strings.Cut(utype, "=")
		if !ok {
			return nil, fmt.Errorf("invalid underlying type option: %s", utype)
		}
		d.Underlying[t] = u
	}

	for _, alias := range opts.ImportAliases {
		p, a, ok := strings.Cut(alias, "=")
		if !ok {
			return nil, fmt.Errorf("invalid alias option: %s", alias)
		}
		d.Aliases[p] = a
	}

	d.AddTypeParameters(mock.TypeParamFields)

	for _, field := range mock.MethodFields {
		d.AppendFuncDef(field)
	}

	for _, field := range mock.Components {
		local := mock.InheritedMethodFields[field.Type.(*ast.Ident).Name]
		for _, lm := range local {
			d.AppendFuncDef(lm)
		}
	}

	for _, field := range mock.ExternalComponents {
		imported := mock.InheritedMethodFields[field.Type.(*ast.SelectorExpr).Sel.Name]
		for _, im := range imported {
			d.AppendFuncDef(im)
		}
	}
	return d, nil
}
