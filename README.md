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

    $ -f <file> -- Go source file 
    $ -o <file> -- output file
    $ -q        -- whether to qualify the type names with the package
    $ -i <name> -- if the source contains multiple interfaces, specify which one to mock
    
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
	
	funcSet: func(v string)  {
		return 
	},
	
}

type mockTestInterfaceOption func(*mockTestInterfaceOptions)


func withFuncGet(f func() string) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcGet = f
	}
}

func withFuncSet(f func(v string) ) mockTestInterfaceOption {
	return func(o *mockTestInterfaceOptions) {
		o.funcSet = f
	}
}



func (m *mockTestInterface) Get() string {
	return m.options.funcGet()
}

func (m *mockTestInterface) Set(v string)  {
	return m.options.funcSet(v)
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

## Current version

0.0.6

## Authors

* **Gabriele Vaccari** - *Initial work* - [Vibridi](https://github.com/vibridi/)

Currently there are no other contributors

## TODOs

None (for now)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
