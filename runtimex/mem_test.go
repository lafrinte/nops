package runtimex

import (
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGoMemoryLimitPercent(t *testing.T) {
	assert := A.New(t)
	assert.NotNil(SetGoMemoryLimitPercent(-1))
	assert.NotNil(SetGoMemoryLimitPercent(101))
	assert.Nil(SetGoMemoryLimitPercent(55))
}

func TestSetGoMemoryLimit(t *testing.T) {
	assert := A.New(t)
	assert.NotNil(SetGoMemoryLimit("1111BXA"))
	assert.NotNil(SetGoMemoryLimit("111MB1"))
	assert.NotNil(SetGoMemoryLimit("1EB"))

	assert.Nil(SetGoMemoryLimit("100000000"))
	assert.Nil(SetGoMemoryLimit("100000000B"))
	assert.Nil(SetGoMemoryLimit("100000KB"))
	assert.Nil(SetGoMemoryLimit("100MB"))
	assert.Nil(SetGoMemoryLimit("1GB"))
	assert.Nil(SetGoMemoryLimit("1TB"))
	assert.Nil(SetGoMemoryLimit("1PB"))
}
