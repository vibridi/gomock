package error

const (
	noTypeDeclarations    = "source contains no type declarations"
	interfaceNotSpecified = "source contains multiple types but no target was specified"
	notFound              = "no suitable type declaration was found in source"
	typeIsStruct          = "found only struct type in source"
	noMethods             = "source interface declares no methods"
)

var NoTypeDeclarations = ParserError{noTypeDeclarations}

var InterfaceNotSpecified = ParserError{interfaceNotSpecified}

var InterfaceNotFound = ParserError{notFound}

var NoMethods = ParserError{noMethods}

type ParserError struct {
	error string
}

func (e ParserError) Error() string {
	return e.error
}
