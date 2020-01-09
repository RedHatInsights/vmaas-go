package utils

import (
	"fmt"
	"runtime"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nAlloc (allocated heap objects) = %v MB", bToMB(m.Alloc))
	fmt.Printf("\nTotalAlloc (cummulative) = %v MB", bToMB(m.TotalAlloc))
	fmt.Printf("\nSys = %v MB", bToMB(m.Sys))
}

func bToMB(bytes uint64) uint64 {
	return bytes / 1e6
}
