package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testInterface = `
package test
type TestInterface interface { 
	Get(a string) 
}`
)

func TestParser(t *testing.T) {
	t.Run("can use", func(t *testing.T) {
		md, err := Parse("", testInterface, "")
		assert.NotNil(t, md)
		assert.Nil(t, err)
		assert.Equal(t, "test", md.PackageName)
		assert.Equal(t, "TestInterface", md.InterfaceName)
		assert.Len(t, md.MethodFields, 1)
	})
}
