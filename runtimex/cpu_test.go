package runtimex

import (
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGoMaxProcPercent(t *testing.T) {
	assert := A.New(t)
	assert.NotNil(SetGoMaxProcPercent(-1))
	assert.NotNil(SetGoMaxProcPercent(101))
	assert.Nil(SetGoMaxProcPercent(55))
}
