//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

type Win32_OperatingSystem struct {
	FreePhysicalMemory     uint64
	FreeSpaceInPagingFiles uint64
	FreeVirtualMemory      uint64
	TotalSwapSpaceSize     uint64
	TotalVirtualMemorySize uint64
	TotalVisibleMemorySize uint64
}

type Win32_PhysicalMemory struct {
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
	// Getting info from WMI "Win32_OperatingSystem"
	var win32OSDescriptions []Win32_OperatingSystem
	q1 := wmi.CreateQuery(&win32OSDescriptions, "")
	if err := wmi.Query(q1, &win32OSDescriptions); err != nil {
		return err
	}
	// Getting info from WMI "Win32_PhysicalMemory"
	var win32MemDescriptions []Win32_PhysicalMemory
	q2 := wmi.CreateQuery(&win32MemDescriptions, "")
	if err := wmi.Query(q2, &win32MemDescriptions); err != nil {
		return err
	}
	// Converting into standard structures
	// Handling physical memory banks
	info.Banks = make([]*MemoryBank, 0, len(win32MemDescriptions))
	for _, description := range win32MemDescriptions {
		info.Banks = append(info.Banks, &MemoryBank{
			Name:         description.Description,
			Label:        description.BankLabel,
			Location:     description.DeviceLocator,
			SerialNumber: description.SerialNumber,
			SizeBytes:    int64(description.Capacity),
			Vendor:       description.Manufacturer,
		})
		//totalPhysicalBytes += description.Capacity
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
