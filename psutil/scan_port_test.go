package psutil

import (
	"fmt"
	A "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewScanPort(t *testing.T) {
	assert := A.New(t)

	_, err := NewScanPort("aaa.aaa.aaa.aaa", []int{0, 1, 100}, 1000, TCP)
	assert.NotNil(err)

	if err != nil {
		assert.Equal(err.Error(), "aaa.aaa.aaa.aaa is not valid ip address")
	}

	_, err = NewScanPort("127.0.0.1", []int{0, 1, 2, 3}, 1000, "g")
	assert.NotNil(err)

	if err != nil {
		assert.Equal(err.Error(), fmt.Sprintf("optional value: (%s|%s), current G", TCP, UDP))
	}

	s, err := NewScanPort("127.0.0.1", []int{0, 1, 1, 100, 200, 999999}, 1000, TCP)
	assert.Nil(err)
	assert.Equal(s.port, []int{1, 100, 200})
}

func TestScan(t *testing.T) {
	assert := A.New(t)

	size := (1 << 14) - 1

	ports := make([]int, size)
	for i := 0; i < size; i++ {
		ports[i] = i
	}

	s, err := NewScanPort("127.0.0.1", ports, 100, TCP)
	assert.Nil(err)

	s.Scan()
}
