package psutil

import (
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMem(t *testing.T) {
	assert := A.New(t)

	assert.True(GetTotalMem() > 0)
	assert.True(GetMemUsed() > 0)

	p := GetMemUsedPercent()
	assert.True(p <= 100 && p >= 0)

	assert.True(GetSwapTotal() > 0)
	assert.True(GetSwapUsed() > 0)

	p = GetSwapUsedPercent()
	assert.True(p < 100 && p >= 0)
}
