package writer

import (
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
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("write with underlying type", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() foo.Foo
}
`,
				out: `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() foo.Foo
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() foo.Foo {
		return 0
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncGet(f func() foo.Foo) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockTestInterface) Get() foo.Foo {
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

			out, err := New(md, WriteOpts{
				Underlying: []string{"foo.Foo=int"},
			}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("write with underlying type of qualified type", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() Foo
}
`,
				out: `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() test.Foo
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() test.Foo {
		return 0
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncGet(f func() test.Foo) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockTestInterface) Get() test.Foo {
	return m.options.funcGet()
}

func newMockTestInterface(opt ...mockTestInterfaceOption) test.TestInterface {
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

			out, err := New(md, WriteOpts{
				Qualify:    true,
				Underlying: []string{"test.Foo=int"},
			}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("write with underlying type of local type", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() Foo
}
`,
				out: `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() Foo
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() Foo {
		return 0
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncGet(f func() Foo) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockTestInterface) Get() Foo {
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

			out, err := New(md, WriteOpts{
				Qualify:    false,
				Underlying: []string{"Foo=int"},
			}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("disambiguate", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface interface {
	Get() int
}
`,
				out: `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() int
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() int {
		return 0
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncTestInterfaceGet(f func() int) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockTestInterface) Get() int {
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

			out, err := New(md, WriteOpts{
				Disambiguate: true,
			}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})
}
