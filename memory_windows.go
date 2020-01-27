//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

type Win32_PhysicalMemory struct {
	BankLabel     string
	Capacity      int64
	DataWidth     int16
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
	TotalWidth    int16
}

type Win32_ComputerSystem struct {
	TotalPhysicalMemory int64
}

/* https://docs.microsoft.com/it-it/windows/win32/cimwin32prov/win32-pagefile
type Win32_PageFile struct {
	FileSize    uint64
	MaximumSize uint64
}*/

func (ctx *context) memFillInfo(info *MemoryInfo) error {
	// Getting info from WMI
	var win32MemDescriptions []Win32_PhysicalMemory
	q1 := wmi.CreateQuery(&win32MemDescriptions, "")
	if err := wmi.Query(q1, &win32MemDescriptions); err != nil {
		return err
	}
	// Converting into standard structures
	info.Banks = make([]*MemoryBank, 0, len(win32MemDescriptions))
	var totalUsableBytes int64
	var totalPhysicalBytes int64
	//var supportedPageSizes []uint64
	for _, description := range win32MemDescriptions {
		info.Banks = append(info.Banks, &MemoryBank{
			Name:         description.Description,
			Label:        description.BankLabel,
			Location:     description.DeviceLocator,
			SerialNumber: description.SerialNumber,
			SizeBytes:    description.Capacity,
			Vendor:       description.Manufacturer,
		})
		//totalUsableBytes += description.Capacity
		totalPhysicalBytes += description.Capacity
	}

	// Getting info from WMI
	var win32SysDescriptions []Win32_ComputerSystem
	q2 := wmi.CreateQuery(&win32SysDescriptions, "")
	if err := wmi.Query(q2, &win32SysDescriptions); err != nil {
		return err
	}
	// Converting into standard structures
	for _, description := range win32SysDescriptions {
		totalUsableBytes += description.TotalPhysicalMemory
	}
	info.TotalUsableBytes = totalUsableBytes
	info.TotalPhysicalBytes = totalPhysicalBytes
	// TODO: find a way to collect these informations
	info.SupportedPageSizes = make([]uint64, 0)

	return nil
}
