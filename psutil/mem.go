package psutil

import (
	"github.com/shirou/gopsutil/mem"
)

func GetTotalMem() uint64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	return m.Total
}

func GetMemUsed() uint64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	return m.Used
}

func GetMemUsedPercent() float64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	return m.UsedPercent
}

func GetSwapTotal() uint64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	return m.SwapTotal
}

func GetSwapUsed() uint64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	return m.SwapTotal - m.SwapFree
}

func GetSwapUsedPercent() float64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}

	if m.SwapTotal == 0 {
		return 0
	}

	return float64(m.SwapTotal-m.SwapFree) * 100 / float64(m.SwapTotal)
}
