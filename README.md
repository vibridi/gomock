# gomock

Automatically generates mock code from a Go interface.  


### Installation

    $ go get -v github.com/vibridi/gomock
    $ cd $GOPATH/github.com/vibridi/gomock
    $ make install

### Requirements:
   
Go 1.11

Optionally: `goimports`, `dep`


## Usage

Type `gomock help` for detailed usage tips.

In short, it supports the following flags:

    -f FILE        Read go code from FILE
    -o FILE        Output mock code to FILE
    -i IDENTIFIER  Mock the interface with IDENTIFIER
    -q             Qualify types with the package name
    -x             Export 'with' and 'new' functions
    
## Example

To try out the tool after cloning the repo:

    $ make build
    $ ./build/gomock -f _example/_example.go

It will print out the generated mock code:


```
type mockTestInterface struct {
	options mockTestInterfaceOptions
}

type mockTestInterfaceOptions struct {
	funcGet  func() string
	funcSet  func(v string) 
	
}

var defaultMockTestInterfaceOptions = mockTestInterfaceOptions{
	funcGet: func() string {
		return ""
	},
	funcSet: func(string)  {
		return 
	},
	
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)


func withFuncGet(f func() string) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func withFuncSet(f func(string) ) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcSet = f
	}
}



func (m *mockTestInterface) Get() string {
	return m.options.funcGet()
}

func (m *mockTestInterface) Set(v string)  {
	return 
}


func newMockTestInterface(opt ...mockTestInterfaceOption) TestInterface {
	opts := defaultMockTestInterfaceOptions
	for _, o := range opt {
		o(&opts)
	}
	return &mockTestInterface{
		options: opts,
	}
}
```

Then you can use the generated code in your unit tests:

```
myMock := newMockTestInterface(
    withFuncGet(f func() string {
        return "test-value"
    }),
)
myMock.Get() // "test-value"

objectThatUsesTestInterface := NewObject(myMock)
// ...

```

## Authors

* **vibridi** - *Initial work* - [Vibridi](https://github.com/vibridi/)

Currently there are no other contributors

## TODOs

* Make unnamed parameters optional in default and with* functions
* Remove extra space between signature and `{` when the function has no return types

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
