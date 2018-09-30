package error

import (
	"github.com/pkg/errors"
)

const (
	cantOpenFile          = "failed to open file"
	notGoSource           = "source is not a Go file"
	writeError            = "could not write output to destination file"
	noSource              = "no source specified"
	noTypeDeclarations    = "source contains no type declarations"
	interfaceNotSpecified = "source contains multiple types but no target was specified"
	notFound              = "no suitable type declaration was found in source"
	noMethods             = "source interface declares no methods"
)

var FileError = GoMockError{cantOpenFile}
var NotGoSource = GoMockError{notGoSource}
var WriteError = GoMockError{writeError}
var NoSource = GoMockError{noSource}
var NoTypeDeclarations = GoMockError{noTypeDeclarations}
var InterfaceNotSpecified = GoMockError{interfaceNotSpecified}
var InterfaceNotFound = GoMockError{notFound}
var NoMethods = GoMockError{noMethods}

type GoMockError struct {
	error string
}

func (e GoMockError) Error() string {
	return e.error
}

func (e GoMockError) Wrap(err error) error {
	return errors.Wrap(err, e.error)
}
