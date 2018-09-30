package version

import "fmt"

// VERSION make build sets this automaticaly
var VERSION string

// GOVERSION make build sets this automaticaly
var GOVERSION string

// Version return string representation of this version
func Version() string {
	return fmt.Sprintf("%s (%s)", VERSION, GOVERSION)
}
