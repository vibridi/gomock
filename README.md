# gomock

Automatically generates mock code from a Go interface.  

### Status

This project is actively maintained. Commits are infrequent because it reached 
a certain degree of stability and no additional features are needed at this time.

### Installation

    $ go install github.com/vibridi/gomock/v3@latest

### Requirements:
   
Go 1.17


## Usage

Type `gomock help` for detailed usage tips. Mainly: `gomock { help | [options] filename }`

In short, it supports the following flags:
   

- `-f FILE` allows to specify the input file where your interface is declared. If not provided, it's assumed 
the input file is the first argument after other options. 
- `-o FILE` if set, tells the program to write the output to `FILE`. Otherwise it just prints to stdout.
You can always capture the output with a pipe. E.g. if you are on MacOS, you could do `gomock -f myfile.go | pbcopy`
- `-i IDENTIFIER` if the input file contains more than one interface declaration, you can use the `-i` flag to tell the program which one to parse.
If not set, the program defaults to the first encountered interface. 
- `-x` if set, static functions are exported (usually those whose name begins with `with` and `new`)
- `-u` if set, allows to output default functions and `With*` functions with unnamed arguments. 
- `--local` if set, doesn't qualify output mock types with the package name. It qualifies them by default.
The default behavior is to always output named arguments, as some IDEs reference them in code completion.
- `--struct` if set, prints the output in struct style, instead of options style (see below for further details).
- `--name NAME` allows to override the interface name used in output types with `NAME`.
- `--help, -h` prints a help message.
- `--version, -v` prints the version number.  

### Breaking changes from version 2.x

- The option `-q` is removed. It's assumed that mocked types are always qualified with their package name. 
The option `--local` can be used instead to opt out (i.e. to not qualify them).  
- The option `--struct-type` has been renamed to `--struct`. It has the same effect.      
    
## Features    
    
This tool is able to resolve composed interfaces, however all declarations must live 
in the same directory or sub-directories relative to the main file. To see this in action, run `make example-compose`.

    
## Example (options style)

To try out the tool after cloning the repo:

    $ make build
    $ ./build/gomock -f _example/_example.go

It will print out the generated mock code, with the options pattern:


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

## Example (struct style)

To print the mock code in struct style, run:

    $ make build
    $ ./build/gomock -f _example/_example.go --struct-style
    
And it will print:

```
type mockTestInterface struct {
	GetFunc  func() string
	SetFunc  func(v string) 
}


func (m *mockTestInterface) Get() string {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	return ""
}

func (m *mockTestInterface) Set(v string)  {
	if m.SetFunc != nil {
		m.SetFunc(v)
	}
}
```

## Authors

* **[Gabriele V.](https://github.com/vibridi/)** - *Initial work and maintenance*

Currently there are no other contributors

## TODOs

* Remove extra space between signature and `{` when the function has no return types

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
