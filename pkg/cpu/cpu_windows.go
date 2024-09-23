//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cpu

import (
	"github.com/StackExchange/wmi"
)

const wmqlProcessor = "SELECT Manufacturer, Name, NumberOfLogicalProcessors, NumberOfCores FROM Win32_Processor"

type win32Processor struct {
	Manufacturer              *string
	Name                      *string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
}

func (i *Info) load() error {
	// Getting info from WMI
	var win32descriptions []win32Processor
	if err := wmi.Query(wmqlProcessor, &win32descriptions); err != nil {
		return err
	}
	// Converting into standard structures
	i.Processors = processorsGet(win32descriptions)
	var totCores uint32
	var totThreads uint32
	for _, p := range i.Processors {
		totCores += p.TotalCores
		totThreads += p.TotalHardwareThreads
	}
	i.TotalCores = totCores
	i.TotalHardwareThreads = totThreads
	// TODO(jaypipes): Remove TotalThreads by v1.0
	i.TotalThreads = totThreads
	return nil
}

func processorsGet(win32descriptions []win32Processor) []*Processor {
	var procs []*Processor
	// Converting into standard structures
	for index, description := range win32descriptions {
		p := &Processor{
			ID:         index,
			Model:      *description.Name,
			Vendor:     *description.Manufacturer,
			TotalCores: description.NumberOfCores,
			// TODO(jaypipes): Remove NumCores before v1.0
			NumCores:             description.NumberOfCores,
			TotalHardwareThreads: description.NumberOfLogicalProcessors,
			// TODO(jaypipes): Remove NumThreads before v1.0
			NumThreads: description.NumberOfLogicalProcessors,
		}
		procs = append(procs, p)
	}
	return procs
}
