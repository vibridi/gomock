package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	VERSION = "v1.0.0"
	GOVERSION = "1.26"
	assert.Equal(t, "v1.0.0 (1.26)", Version())
}
