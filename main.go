package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	throws "github.com/vibridi/gomock/error"
	"github.com/vibridi/gomock/parser"
	"github.com/vibridi/gomock/version"
	"github.com/vibridi/gomock/writer"

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
		sourceFile  string
		destination string
		target      string
		noQualify   bool
		export      bool
		unnamedsig  bool
		structStyle bool
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
			Name:        "local",
			Usage:       "Don't qualify types with the package name",
			Destination: &noQualify,
		},
		&cli.BoolFlag{
			Name:        "struct",
			Usage:       "Prints the output mock in struct style (default: options style)",
			Destination: &structStyle,
		},
	}

	app.Action = func(c *cli.Context) error {
		if sourceFile == "" {
			sourceFile = c.Args().Get(0)
		}
		_, _ = fmt.Fprintf(os.Stderr, "parsing %s\n", sourceFile)

		if !strings.HasSuffix(sourceFile, ".go") {
			return throws.NotGoSource
		}

		f, err := filepath.Abs(sourceFile)
		if err != nil {
			return throws.FileError
		}

		md, err := parser.Parse(f, nil, target)
		if err != nil {
			return err
		}

		w := writer.New(
			md,
			writer.WriteOpts{
				Qualify:          !noQualify,
				Export:           export,
				UnnamedSignature: unnamedsig,
				StructStyle:      structStyle,
			},
		)
		out, err := w.Write()
		if err != nil {
			return throws.WriteError
		}

		if destination == "" {
			fmt.Println(string(out))
			return nil
		}
		if err := ioutil.WriteFile(destination, out, 0644); err != nil {
			return throws.WriteError.Wrap(err)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
