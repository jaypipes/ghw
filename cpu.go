//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

// ProcessorCore describes a physical host processor core. A processor core is
// a separate processing unit within some types of central processing units
// (CPU).
type ProcessorCore struct {
	// TODO(jaypipes): Deprecated in 0.2, remove in 1.0
	Id                int
	ID                int
	Index             int
	NumThreads        uint32
	LogicalProcessors []int
}

func (c *ProcessorCore) String() string {
	return fmt.Sprintf(
		"processor core #%d (%d threads), logical processors %v",
		c.Index,
		c.NumThreads,
		c.LogicalProcessors,
	)
}

// Processor describes a physical host central processing unit (CPU).
type Processor struct {
	// TODO(jaypipes): Deprecated in 0.2, remove in 1.0
	Id           int
	ID           int
	NumCores     uint32
	NumThreads   uint32
	Vendor       string
	Model        string
	Capabilities []string
	Cores        []*ProcessorCore
}

// HasCapability returns true if the `ghw.Processor` has the supplied cpuid
// capability, false otherwise. Example of cpuid capabilities would be 'vmx' or
// 'sse4_2'. To see a list of potential cpuid capabilitiies, see the section on
// CPUID feature bits in the following article:
//
// https://en.wikipedia.org/wiki/CPUID
func (p *Processor) HasCapability(find string) bool {
	for _, c := range p.Capabilities {
		if c == find {
			return true
		}
	}
	return false
}

func (p *Processor) String() string {
	ncs := "cores"
	if p.NumCores == 1 {
		ncs = "core"
	}
	nts := "threads"
	if p.NumThreads == 1 {
		nts = "thread"
	}
	return fmt.Sprintf(
		"physical package #%d (%d %s, %d hardware %s)",
		p.ID,
		p.NumCores,
		ncs,
		p.NumThreads,
		nts,
	)
}

// CPUInfo describes all central processing unit (CPU) functionality on a host.
// Returned by the `ghw.CPU()` function.
type CPUInfo struct {
	TotalCores   uint32
	TotalThreads uint32
	Processors   []*Processor
}

// CPU returns a struct containing information about the host's CPU resources.
func CPU(opts ...*WithOption) (*CPUInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &CPUInfo{}
	if err := ctx.cpuFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *CPUInfo) String() string {
	nps := "packages"
	if len(i.Processors) == 1 {
		nps = "package"
	}
	ncs := "cores"
	if i.TotalCores == 1 {
		ncs = "core"
	}
	nts := "threads"
	if i.TotalThreads == 1 {
		nts = "thread"
	}
	return fmt.Sprintf(
		"cpu (%d physical %s, %d %s, %d hardware %s)",
		len(i.Processors),
		nps,
		i.TotalCores,
		ncs,
		i.TotalThreads,
		nts,
	)
}
