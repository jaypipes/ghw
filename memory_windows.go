//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

const wmqlMemory = "SELECT FreePhysicalMemory, FreeSpaceInPagingFiles, FreeVirtualMemory, TotalSwapSpaceSize, TotalVirtualMemorySize, TotalVisibleMemorySize FROM Win32_OperatingSystem"

type win32OperatingSystem struct {
	FreePhysicalMemory     uint64
	FreeSpaceInPagingFiles uint64
	FreeVirtualMemory      uint64
	TotalSwapSpaceSize     uint64
	TotalVirtualMemorySize uint64
	TotalVisibleMemorySize uint64
}

func (ctx *context) memFillInfo(info *MemoryInfo) error {
	// Getting info from WMI
	var win32OSDescriptions []win32OperatingSystem
	if err := wmi.Query(wmqlMemory, &win32OSDescriptions); err != nil {
		return err
	}
	// Handling physical memory total/free size (as seen by OS)
	var totalUsableBytes uint64
	var totalPhysicalBytes uint64
	for _, description := range win32OSDescriptions {
		totalUsableBytes += description.FreePhysicalMemory
		totalPhysicalBytes += description.TotalVisibleMemorySize
	}
	info.TotalUsableBytes = int64(totalUsableBytes)
	info.TotalPhysicalBytes = int64(totalPhysicalBytes)
	return nil
}
