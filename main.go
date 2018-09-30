package main

import (
	"fmt"
	throws "gomock/error"
	"gomock/parser"
	"gomock/version"
	"gomock/writer"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "gomock"
	app.Version = version.Version()

	var src string
	var dst string
	var tgt string
	var qualify bool

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "Read go code from `FILE`",
			Destination: &src,
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

		fmt.Printf("parsing %s\n", src)

		if !strings.HasSuffix(src, ".go") {
			return throws.NotGoSource
		}

		f, err := filepath.Abs(src)
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
		fmt.Printf("error: %s\n", err.Error())
	}
}
