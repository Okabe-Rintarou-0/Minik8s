package utility

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func GetCpuAndMemoryUsage() (cpu, mem float64) {
	cpu = 0.0
	mem = 0.0
	cpuCh := make(chan float64, 1)
	memCh := make(chan float64, 1)
	go func() {
		cpuCh <- CpuUsageUnix()
	}()
	go func() {
		memCh <- MemoryUsageUnix()
	}()
	cpu = <-cpuCh
	mem = <-memCh
	return
}

func CpuUsageUnix() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func MemoryUsageUnix() float64 {
	memory, _ := mem.VirtualMemory()
	return memory.UsedPercent
}
