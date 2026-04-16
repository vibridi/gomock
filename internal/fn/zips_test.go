package fn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZips(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		s := Zips(nil, nil, "")
		assert.Len(t, s, 0)
	})

	t.Run("both empty", func(t *testing.T) {
		s := Zips([]string{}, []string{}, "")
		assert.Len(t, s, 0)
	})

	t.Run("first arg longer", func(t *testing.T) {
		s := Zips([]string{"a", "b"}, []string{"1"}, "")
		assert.Equal(t, []string{"a1", "b"}, s)
	})

	t.Run("second arg longer", func(t *testing.T) {
		s := Zips([]string{"a"}, []string{"1", "2"}, "")
		assert.Equal(t, []string{"a1", "2"}, s)
	})

	t.Run("same length", func(t *testing.T) {
		s := Zips([]string{"a", "b"}, []string{"1", "2"}, "---")
		assert.Equal(t, []string{"a---1", "b---2"}, s)
	})
}
