package writer

import (
	"fmt"
	"testing"

	gomock "github.com/vibridi/gomock/v3/parser"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	t.Run("write options template", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() string
}
`,
				out: `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() string
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() string {
		return ""
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncGet(f func() string) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockTestInterface) Get() string {
	return m.options.funcGet()
}

func newMockTestInterface(opt ...mockTestInterfaceOption) TestInterface {
	opts := defaultMockTestInterfaceOptions
	for _, o := range opt {
		o(&opts)
	}
	return &mockTestInterface{
		options: opts,
	}
}`,
			},
		}

		for _, c := range cases {
			md, err := gomock.Parse("", c.in, "")
			assert.Nil(t, err)

			out, err := New(md, WriteOpts{}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("write struct template", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() string
}
`,
				out: `
type mockTestInterface struct {
	GetFunc func() string
}

func (m *mockTestInterface) Get() string {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return ""
}
`,
			},
		}

		for _, c := range cases {
			md, err := gomock.Parse("", c.in, "")
			assert.Nil(t, err)

			out, err := New(md, WriteOpts{StructStyle: true}).Write()
			fmt.Println(err)
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("override name", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() string
}
`,
				out: `
type mockOverridden struct {
	GetFunc func() string
}

func (m *mockOverridden) Get() string {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return ""
}
`,
			},
		}

		for _, c := range cases {
			md, err := gomock.Parse("", c.in, "")
			assert.Nil(t, err)

			out, err := New(md, WriteOpts{
				StructStyle: true,
				MockName:    "Overridden",
			}).Write()
			fmt.Println(err)
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})
}
