package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("option conflict", func(t *testing.T) {
		err := run([]string{"gomock", "-p", "--name", "Foo"})
		assert.Equal(t, "option conflict: specify only one of --name and -p", err.Error())
	})

	t.Run("source not go", func(t *testing.T) {
		err := run([]string{"gomock", "-f", "Foo"})
		assert.Equal(t, "source is not a Go file", err.Error())
	})

	t.Run("source fallback not go", func(t *testing.T) {
		err := run([]string{"gomock", "Foo"})
		assert.Equal(t, "source is not a Go file", err.Error())
	})

	t.Run("parsing failed", func(t *testing.T) {
		tmpdir := t.TempDir()
		tmpfile := tmpdir + "/foo.go"
		err := os.WriteFile(tmpfile, []byte("blah blah"), 0644)
		require.Nil(t, err)

		err = run([]string{"gomock", "-f", tmpfile})
		assert.Contains(t, err.Error(), "cannot parse source")
	})

	t.Run("writing failed", func(t *testing.T) {
		tmpdir := t.TempDir()
		tmpfile := tmpdir + "/foo.go"
		err := os.WriteFile(tmpfile, []byte("package foo\n\ntype Foo interface {\nDo() error\n}"), 0644)
		require.Nil(t, err)

		err = run([]string{"gomock", "-f", tmpfile, "--utype", "foo/bar"})
		assert.Contains(t, err.Error(), "failed to write output")
	})

	t.Run("write to stdout", func(t *testing.T) {
		tmpdir := t.TempDir()
		tmpfile := tmpdir + "/foo.go"
		err := os.WriteFile(tmpfile, []byte("package foo\n\ntype Foo interface {\nDo() error\n}"), 0644)
		require.Nil(t, err)

		stdout := os.Stdout
		r, w, err := os.Pipe()
		require.Nil(t, err)

		os.Stdout = w
		defer func() {
			os.Stdout = stdout
		}()

		err = run([]string{"gomock", "-f", tmpfile})
		require.Nil(t, err)
		_ = w.Close()

		out := readstr(r)
		assert.Contains(t, out, "options mockFooOptions")
	})

	t.Run("write to file", func(t *testing.T) {
		tmpdir := t.TempDir()
		tmpfile := tmpdir + "/foo.go"
		outfile := tmpdir + "/out.go"
		err := os.WriteFile(tmpfile, []byte("package foo\n\ntype Foo interface {\nDo() error\n}"), 0644)
		require.Nil(t, err)

		err = run([]string{"gomock", "-f", tmpfile, "-o", outfile})
		require.Nil(t, err)

		f, err := os.Open(outfile)
		require.Nil(t, err)
		defer f.Close()

		out := readstr(f)
		assert.Contains(t, out, "options mockFooOptions")
	})

}

func readstr(f *os.File) string {
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}
