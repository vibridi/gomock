package writer

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vibridi/gomock/v3/writer/template"
)

func TestWriteFile(t *testing.T) {
	tmpdir := t.TempDir()
	tmpfile := tmpdir + "/foo.go"
	data := []byte("func foo() {}")

	t.Run("file does not exist", func(t *testing.T) {
		err := File(tmpfile, "foo", data)
		require.Nil(t, err)
		b, err := os.ReadFile(tmpfile)
		require.Nil(t, err)
		assert.Equal(t, "package foo\n\nfunc foo() {}", string(b))
	})

	t.Run("uses notice comment when present", func(t *testing.T) {
		src := `
package blah

// not deleted

` + strings.TrimSpace(template.Notice) + `

func main() {}
`
		err := os.WriteFile(tmpfile, []byte(src), 0644)
		require.Nil(t, err)
		err = File(tmpfile, "", data)
		require.Nil(t, err)
		b, err := os.ReadFile(tmpfile)
		require.Nil(t, err)

		want := `
package blah

// not deleted

func foo() {}`
		assert.Equal(t, want, string(b))
	})

	t.Run("end of import block", func(t *testing.T) {
		src := `
package blah

import (
	"fmt"
	"os"
)

func main() {}
`
		err := os.WriteFile(tmpfile, []byte(src), 0644)
		require.Nil(t, err)
		err = File(tmpfile, "", append([]byte("\n\n"), data...))
		require.Nil(t, err)
		b, err := os.ReadFile(tmpfile)
		require.Nil(t, err)

		want := `
package blah

import (
	"fmt"
	"os"
)

func foo() {}`
		assert.Equal(t, want, string(b))
	})

	t.Run("end of last import block", func(t *testing.T) {
		src := `
package test

import (
	"fmt"
	"os"
)
import (
	"context"
	"time"
)

func main() {}
`
		err := os.WriteFile(tmpfile, []byte(src), 0644)
		require.Nil(t, err)
		err = File(tmpfile, "", append([]byte("\n\n"), data...))
		require.Nil(t, err)
		b, err := os.ReadFile(tmpfile)
		require.Nil(t, err)

		want := `
package test

import (
	"fmt"
	"os"
)
import (
	"context"
	"time"
)

func foo() {}`
		assert.Equal(t, want, string(b))
	})

	t.Run("no imports no comment", func(t *testing.T) {
		src := `
package main

func main() {}
`
		err := os.WriteFile(tmpfile, []byte(src), 0644)
		require.Nil(t, err)
		err = File(tmpfile, "", append([]byte("\n\n"), data...))
		require.Nil(t, err)
		b, err := os.ReadFile(tmpfile)
		require.Nil(t, err)

		want := `
package main

func foo() {}`
		assert.Equal(t, want, string(b))
	})
}
