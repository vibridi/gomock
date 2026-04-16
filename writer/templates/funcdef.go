package templates

import "strings"

// FuncDef represents the information needed to output mocks based on the interface's methods
type FuncDef struct {
	ServiceName  string // Name of the interface that is being mocked (can be ovverridden by some options)
	Name         string // Identifier of this function
	Signature    string // Full parameter list of this function excluding brackets
	Return       string // Full return parameter list of this function including brackets
	Args         string // List of function arguments
	ReturnValues string // List of values that can appear in this function's return statement
}

func (fd FuncDef) String() string {
	s := fd.Name + "(" + fd.Signature + ") " + fd.Return
	return strings.TrimSpace(s)
}

type ParamName struct {
	string
	IsVararg bool
}

func (pn ParamName) Expand() string {
	if pn.IsVararg {
		return pn.string + "..."
	}
	return pn.string
}
