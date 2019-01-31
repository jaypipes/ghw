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
	Id                int    `json:"-"`
	ID                int    `json:"id"`
	Index             int    `json:"index"`
	NumThreads        uint32 `json:"total_threads"`
	LogicalProcessors []int  `json:"logical_processors"`
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
	Id           int              `json:"-"`
	ID           int              `json:"id"`
	NumCores     uint32           `json:"total_cores"`
	NumThreads   uint32           `json:"total_threads"`
	Vendor       string           `json:"vendor"`
	Model        string           `json:"model"`
	Capabilities []string         `json:"capabilities"`
	Cores        []*ProcessorCore `json:"cores"`
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
	TotalCores   uint32 `json:"total_cores"`
	TotalThreads uint32 `json:"total_threads"`

	Processors []*Processor `json:"processors"`
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

// simple private struct used to encapsulate cpu information in a top-level
// "cpu" YAML/JSON map/object key
type cpuPrinter struct {
	Info *CPUInfo `json:"cpu"`
}

// YAMLString returns a string with the cpu information formatted as YAML
// under a top-level "cpu:" key
func (i *CPUInfo) YAMLString() string {
	return safeYAML(cpuPrinter{i})
}

// JSONString returns a string with the cpu information formatted as JSON
// under a top-level "cpu:" key
func (i *CPUInfo) JSONString(indent bool) string {
	return safeJSON(cpuPrinter{i}, indent)
}
