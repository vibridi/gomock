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

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "gomock"
	app.Version = version.Version()

	var srcFile string
	var dst string
	var tgt string
	var qualify bool

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "Read go code from `FILE`",
			Destination: &srcFile,
		},
		cli.StringFlag{
			Name:        "o",
			Usage:       "Output mock code to `FILE`",
			Value:       "",
			Destination: &dst,
		},
		cli.StringFlag{
			Name:        "i",
			Usage:       "Mock the interface with `IDENTIFIER`",
			Value:       "",
			Destination: &tgt,
		},
		cli.BoolFlag{
			Name:        "q",
			Usage:       "Qualify types with the package name",
			Destination: &qualify,
		},
	}

	app.Action = func(c *cli.Context) error {
		if srcFile == "" {
			srcFile = c.Args().Get(0)
		}
		fmt.Errorf("parsing %s\n", srcFile)

		if !strings.HasSuffix(srcFile, ".go") {
			return throws.NotGoSource
		}

		f, err := filepath.Abs(srcFile)
		if err != nil {
			return throws.FileError
		}

		md, err := parser.Parse(f, nil, tgt)
		if err != nil {
			return err
		}

		out, err := writer.Write(md, qualify)
		if err != nil {
			return throws.WriteError
		}

		if dst == "" {
			fmt.Println(out)
			return nil
		}

		if err := ioutil.WriteFile(dst, []byte(out), 0644); err != nil {
			return throws.WriteError.Wrap(err)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Errorf("error: %s\n", err.Error())
	}
}
