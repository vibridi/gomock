package templates

import "strings"

type FuncDef struct {
	ServiceName  string
	Name         string
	Signature    string
	Return       string
	Args         string
	ReturnValues string
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
