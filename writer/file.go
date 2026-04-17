package writer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/vibridi/gomock/v3/writer/template"
)

func File(destination string, pkg string, text []byte) error {
	dst, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	fi, err := dst.Stat()
	if err != nil {
		return err
	}

	pos := 0
	if fi.Size() > 0 {
		b, err := io.ReadAll(dst)
		if err != nil {
			return err
		}
		pos, err = getWritePos(b)
		if err != nil {
			return err
		}
	}

	if pos > 0 {
		dst.Truncate(int64(pos))
		dst.Seek(0, io.SeekEnd)
	} else {
		s := strings.Split(destination, "/")
		if pkg == "" {
			if len(s) > 1 {
				pkg = s[len(s)-2]
			} else {
				pkg = "main"
			}
		}

		dst.WriteString("package " + pkg)
		dst.WriteString("\n\n")
	}

	_, err = dst.WriteString(template.Notice)
	_, err = dst.Write(text)
	return err
}

func getWritePos(src []byte) (int, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return 0, fmt.Errorf("parse file error: %w", err)
	}

	var pos token.Pos
	for _, cmt := range f.Comments {
		if strings.TrimSpace(cmt.Text()) == strings.Trim(template.Notice, "/ \n") {
			pos = cmt.Pos()
		}
	}
	if pos == 0 {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if ok && gen.Tok == token.IMPORT {
				pos = max(pos, gen.End())
			}
		}
	}
	if pos == 0 && f.Package.IsValid() {
		pos = f.Name.End()
	}
	return fset.Position(pos).Offset, nil
}
