package str

import (
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestToBytes(t *testing.T) {
	assert := A.New(t)

	s := "1234567890"

	assert.Equal([]byte(s), ToBytes(s))
}

func TestToString(t *testing.T) {
	assert := A.New(t)

	s := "1234567890"
	b := []byte(s)

	assert.Equal(ToString(b), s)
}
