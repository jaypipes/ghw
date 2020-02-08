//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

const wqlOperatingSystem = "SELECT FreePhysicalMemory, FreeSpaceInPagingFiles, FreeVirtualMemory, TotalSwapSpaceSize, TotalVirtualMemorySize, TotalVisibleMemorySize FROM Win32_OperatingSystem"

type win32OperatingSystem struct {
	FreePhysicalMemory     uint64
	FreeSpaceInPagingFiles uint64
	FreeVirtualMemory      uint64
	TotalSwapSpaceSize     uint64
	TotalVirtualMemorySize uint64
	TotalVisibleMemorySize uint64
}

const wqlPhysicalMemory = "SELECT BankLabel, Capacity, DataWidth, Description, DeviceLocator, Manufacturer, Model, Name, PartNumber, PositionInRow, SerialNumber, Speed, Tag, TotalWidth FROM Win32_PhysicalMemory"

type win32PhysicalMemory struct {
	BankLabel     string
	Capacity      uint64
	DataWidth     uint16
	Description   string
	DeviceLocator string
	Manufacturer  string
	Model         string
	Name          string
	PartNumber    string
	PositionInRow uint32
	SerialNumber  string
	Speed         uint32
	Tag           string
	TotalWidth    uint16
}

func (ctx *context) memFillInfo(info *MemoryInfo) error {
	// Getting info from WMI
	var win32OSDescriptions []win32OperatingSystem
	if err := wmi.Query(wqlOperatingSystem, &win32OSDescriptions); err != nil {
		return err
	}
	var win32MemDescriptions []win32PhysicalMemory
	if err := wmi.Query(wqlPhysicalMemory, &win32MemDescriptions); err != nil {
		return err
	}
	// Converting into standard structures
	// Handling physical memory modules
	info.Modules = make([]*MemoryModule, 0, len(win32MemDescriptions))
	for _, description := range win32MemDescriptions {
		info.Modules = append(info.Modules, &MemoryModule{
			Label:        description.BankLabel,
			Location:     description.DeviceLocator,
			SerialNumber: description.SerialNumber,
			SizeBytes:    int64(description.Capacity),
			Vendor:       description.Manufacturer,
		})
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
