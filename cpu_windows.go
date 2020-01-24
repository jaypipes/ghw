// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

type Win32_Processor struct {
	//LoadPercentage            *uint16
	//Family                    uint16
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
	//ProcessorID               *string
	//Stepping                  *string
	//MaxClockSpeed             uint32
}

func (ctx *context) cpuFillInfo(info *CPUInfo) error {
	// Getting info from WMI
	var win32descriptions []Win32_Processor
	q := wmi.CreateQuery(&win32descriptions, "")
	if err := wmi.Query(q, &win32descriptions); err != nil {
		return err
	}
	// Converting into standard structures
	info.Processors = ctx.processorsGet(win32descriptions)
	var totCores uint32
	var totThreads uint32
	for _, p := range info.Processors {
		totCores += p.NumCores
		totThreads += p.NumThreads
	}
	info.TotalCores = totCores
	info.TotalThreads = totThreads
	return nil
}

func (ctx *context) processorsGet(win32descriptions []Win32_Processor) []*Processor {
	var procs []*Processor
	// Converting into standard structures
	for index, description := range win32descriptions {
		p := &Processor{
			Id:         index,            // TODO: how to get a decent "Physical ID" to use ?
			Model:      description.Name, // description.ProcessorID description.Manufacturer description.Name
			Vendor:     description.Manufacturer,
			NumCores:   description.NumberOfCores,
			NumThreads: description.NumberOfLogicalProcessors,
			// TODO: find a way to get these informations since it seems wmi command doesn't offer them
			//Cores:
			//Capabilities:
		}
		procs = append(procs, p)
	}
	return procs
}
