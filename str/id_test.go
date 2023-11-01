package str

import (
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestID(t *testing.T) {
	assert := A.New(t)

	i := ID()
	ii := ID()

	assert.NotEqual(i.String(), ii.String())
}
