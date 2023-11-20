package psutil

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

func VCPUCount() int {
	c, _ := cpu.Counts(true)
	return c
}

func CPUCount() int {
	c, _ := cpu.Counts(false)
	return c
}

func CPUInfo() []cpu.InfoStat {
	i, _ := cpu.Info()
	return i
}

func TotalCPUPercent() []float64 {
	p, _ := cpu.Percent(time.Second*3, false)
	return p
}

func PerCPUPercent() []float64 {
	p, _ := cpu.Percent(time.Second*3, true)
	return p
}
