package psutil

import (
	A "github.com/stretchr/testify/assert"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestVCPUCount(t *testing.T) {
	assert := A.New(t)
	assert.Equal(runtime.NumCPU(), VCPUCount())

}

func TestCPUCount(t *testing.T) {
	assert := A.New(t)

	switch runtime.GOOS {
	case "darwin", "linux", "freebsd", "openbsd":
		cmd := exec.Command("sysctl", "-n", "hw.physicalcpu")
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("failed run driver module")
			return
		}

		physicalCPUs := strings.TrimSpace(string(output))
		assert.Equal(physicalCPUs, strconv.Itoa(CPUCount()))
	default:
	}
}

func TestCPUInfo(t *testing.T) {
	assert := A.New(t)

	state := CPUInfo()
	assert.NotEqual(len(state), 0)
}

func TestCPUPercent(t *testing.T) {
	assert := A.New(t)

	assert.NotNil(TotalCPUPercent)
	assert.NotNil(PerCPUPercent)
}
