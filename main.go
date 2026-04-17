package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vibridi/gomock/v3/parser"
	"github.com/vibridi/gomock/v3/version"
	"github.com/vibridi/gomock/v3/writer"
	"github.com/vibridi/gomock/v3/writer/template"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "gomock"
	app.Usage = "simple interface mocking tool"
	app.UsageText = "gomock { help | [options] filename }"
	app.UseShortOptionHandling = true
	app.Version = version.Version()

	var (
		sourceFile    string
		destination   string
		target        string
		noQualify     bool
		export        bool
		unnamedsig    bool
		structStyle   bool
		mockName      string
		underlying    cli.StringSlice
		disambiguate  bool
		prefixPackage bool
	)

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "f",
			Usage:       "Read input from `FILE`. Must be valid Go code",
			Destination: &sourceFile,
		},
		&cli.StringFlag{
			Name:        "o",
			Usage:       "Write output to `FILE`",
			Value:       "",
			Destination: &destination,
		},
		&cli.StringFlag{
			Name:        "i",
			Usage:       "Mock the interface named `IDENTIFIER`",
			Value:       "",
			Destination: &target,
		},
		&cli.BoolFlag{
			Name:        "x",
			Usage:       "Export 'with' and 'new' functions",
			Destination: &export,
		},
		&cli.BoolFlag{
			Name:        "u",
			Usage:       "Output func signatures with unnamed parameters where possible",
			Destination: &unnamedsig,
		},
		&cli.BoolFlag{
			Name:        "d",
			Usage:       "Disambiguate 'withFunc' identifiers with service name, e.g. withFuncMyServiceGet()",
			Destination: &disambiguate,
		},
		&cli.BoolFlag{
			Name:        "p",
			Usage:       "Merge the package name and the mock name in function identifiers, e.g. foo.Client gives NewMockFooClient",
			Destination: &prefixPackage,
		},
		&cli.BoolFlag{
			Name:        "local",
			Usage:       "Don't qualify types with the package name",
			Destination: &noQualify,
		},
		&cli.BoolFlag{
			Name:        "struct",
			Usage:       "Prints the output mock in struct style (default: options style)",
			Destination: &structStyle,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "Use `NAME` in output types instead of the name of the mocked interface",
			Value:       "",
			Destination: &mockName,
		},
		&cli.StringSliceFlag{
			Name:        "utype",
			Usage:       "Maps a type to its underlying type. `MAPPING` must in the format 'type=underlying'",
			Value:       nil,
			Destination: &underlying,
		},
	}

	app.Action = func(c *cli.Context) error {
		if mockName != "" && prefixPackage {
			return fmt.Errorf("option conflict: specify only one of --name and -p")
		}

		if sourceFile == "" {
			sourceFile = c.Args().Get(0)
		}
		_, _ = fmt.Fprintf(os.Stderr, "parsing %s\n", sourceFile)

		if !strings.HasSuffix(sourceFile, ".go") {
			return errors.New("source is not a Go file")
		}

		f, err := filepath.Abs(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		md, err := parser.Parse(f, nil, target)
		if err != nil {
			return err
		}

		out, err := template.Exec(
			md,
			template.Opts{
				Qualify:          !noQualify,
				Export:           export,
				UnnamedSignature: unnamedsig,
				StructStyle:      structStyle,
				Disambiguate:     disambiguate,
				MockName:         mockName,
				Underlying:       underlying.Value(),
				PrefixPackage:    prefixPackage,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		if destination == "" {
			fmt.Println(string(out))
			return nil
		}

		if err := writer.File(destination, "", out); err != nil {
			return fmt.Errorf("failed to write destination file: %w", err)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
