package writer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

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
	foo
	Get() string
}
type foo interface {
	Do() error
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
			assert.NotEmpty(t, md.Components)

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

	t.Run("inherited methods", func(t *testing.T) {
		tmpdir := t.TempDir()
		testgo := `
package test
type TestInterface interface {
	foo
	Get() string
}`
		err := os.MkdirAll(tmpdir+"/test", 0755)
		require.Nil(t, err)
		err = os.MkdirAll(tmpdir+"/bar", 0755)
		require.Nil(t, err)

		err = os.WriteFile(tmpdir+"/test/test.go", []byte(testgo), 0644)
		require.Nil(t, err)

		foogo := `
package test
type foo interface {
	Do() error
}`
		err = os.WriteFile(tmpdir+"/test/foo.go", []byte(foogo), 0644)
		require.Nil(t, err)

		md, err := gomock.Parse(tmpdir+"/test/test.go", nil, "")
		require.Nil(t, err)

		want := `
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet func() string
	funcDo  func() error
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() string {
		return ""
	},
	funcDo: func() error {
		return nil
	},
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)

func withFuncGet(f func() string) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func withFuncDo(f func() error) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcDo = f
	}
}

func (m *mockTestInterface) Get() string {
	return m.options.funcGet()
}

func (m *mockTestInterface) Do() error {
	return m.options.funcDo()
}

func newMockTestInterface(opt ...mockTestInterfaceOption) TestInterface {
	opts := defaultMockTestInterfaceOptions
	for _, o := range opt {
		o(&opts)
	}
	return &mockTestInterface{
		options: opts,
	}
}`
		out, err := New(md, WriteOpts{}).Write()
		assert.Nil(t, err)
		assert.Equal(t, want, string(out))
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

	t.Run("prefix", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package foo
type Interface interface {
	Get() int
}
`,
				out: `
type mockFooInterface struct {
	options mockFooInterfaceOptions
}

type mockFooInterfaceOptions struct {
	funcGet func() int
}

var defaultMockFooInterfaceOptions = mockFooInterfaceOptions{
	funcGet: func() int {
		return 0
	},
}

type mockFooInterfaceOption func(*mockFooInterfaceOptions)

func WithFuncFooInterfaceGet(f func() int) mockFooInterfaceOption {
	return func(o *mockFooInterfaceOptions) {
		o.funcGet = f
	}
}

func (m *mockFooInterface) Get() int {
	return m.options.funcGet()
}

func NewMockFooInterface(opt ...mockFooInterfaceOption) foo.Interface {
	opts := defaultMockFooInterfaceOptions
	for _, o := range opt {
		o(&opts)
	}
	return &mockFooInterface{
		options: opts,
	}
}`,
			},
		}

		for _, c := range cases {
			md, err := gomock.Parse("", c.in, "")
			assert.Nil(t, err)

			out, err := New(md, WriteOpts{
				Qualify:       true,
				Export:        true,
				Disambiguate:  true,
				PrefixPackage: true,
			}).Write()
			assert.Nil(t, err)
			assert.Equal(t, c.out, string(out))
		}
	})

	t.Run("generic interface", func(t *testing.T) {
		cases := []struct {
			in  string
			out string
		}{
			{
				in: `
package test
type TestInterface[T any, R ~int] interface {
	Get() R
	Foo(v T)
}
`,
				out: `
type mockTestInterface[T any, R ~int] struct {
	options mockTestInterfaceOptions[T, R]
}

type mockTestInterfaceOptions[T any, R ~int] struct {
	funcGet func() R
	funcFoo func(v T)
}

func newDefaultMockTestInterfaceOptions[T any, R ~int]() mockTestInterfaceOptions[T, R] {
	return mockTestInterfaceOptions[T, R]{
		funcGet: func() R {
			return *new(R)
		},
		funcFoo: func(v T) {
			return
		},
	}
}

type mockTestInterfaceOption[T any, R ~int] func(*mockTestInterfaceOptions[T, R])

func withFuncGet[T any, R ~int](f func() R) mockTestInterfaceOption[T, R] {
	return func(o *mockTestInterfaceOptions[T, R]) {
		o.funcGet = f
	}
}

func withFuncFoo[T any, R ~int](f func(v T)) mockTestInterfaceOption[T, R] {
	return func(o *mockTestInterfaceOptions[T, R]) {
		o.funcFoo = f
	}
}

func (m *mockTestInterface[T, R]) Get() R {
	return m.options.funcGet()
}

func (m *mockTestInterface[T, R]) Foo(v T) {
	return
}

func newMockTestInterface[T any, R ~int](opt ...mockTestInterfaceOption[T, R]) TestInterface[T, R] {
	opts := newDefaultMockTestInterfaceOptions[T, R]()
	for _, o := range opt {
		o(&opts)
	}
	return &mockTestInterface[T, R]{
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
}
