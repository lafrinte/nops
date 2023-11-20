package psutil

import (
	A "github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestPlatform(t *testing.T) {
	assert := A.New(t)

	assert.NotEqual(GetOS(), "")
	assert.NotEqual(GetHostname(), "")
	assert.NotNil(GetPlatform(), "")
	assert.NotEqual(GetPlatformVersion(), "")

	assert.False(IsSunOs())
	assert.False(IsSmartOs())
	assert.False(IsOpenBsd())
	assert.False(IsFreeBsd())
	assert.False(IsAix())

	switch runtime.GOOS {
	case "darwin":
		assert.True(IsDarwin())
		assert.False(IsWindows())
		assert.False(IsLinux())
	case "linux":
		assert.False(IsDarwin())
		assert.False(IsWindows())
		assert.True(IsLinux())
	case "windows":
		assert.False(IsDarwin())
		assert.True(IsWindows())
		assert.False(IsLinux())
	}

	switch runtime.GOARCH {
	case "amd64":
		assert.True(IsX86())
		assert.False(IsArch64())
	case "arm":
		assert.False(IsX86())
		assert.True(IsArch64())
	default:
	}

}
