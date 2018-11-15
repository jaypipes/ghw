//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"sort"
)

// Architecture describes the overall hardware architecture. It can be either
// Symmetric Multi-Processor (SMP) or Non-Uniform Memory Access (NUMA)
type Architecture int

const (
	// SMP is a Symmetric Multi-Processor system
	SMP Architecture = iota
	// NUMA is a Non-Uniform Memory Access system
	NUMA
)

// TopologyNode is an abstract construct representing a collection of
// processors and various levels of memory cache that those processors share.
// In a NUMA architecture, there are multiple NUMA nodes, abstracted here as
// multiple TopologyNode structs. In an SMP architecture, a single TopologyNode
// will be available in the TopologyInfo struct and this single struct can be
// used to describe the levels of memory caching available to the single
// physical processor package's physical processor cores
type TopologyNode struct {
	// TODO(jaypipes): Deprecated in 0.2, remove in 1.0
	Id     int
	ID     int
	Cores  []*ProcessorCore
	Caches []*MemoryCache
}

func (n *TopologyNode) String() string {
	return fmt.Sprintf(
		"node #%d (%d cores)",
		n.ID,
		len(n.Cores),
	)
}

// TopologyInfo describes the system topology for the host hardware
type TopologyInfo struct {
	Architecture Architecture
	Nodes        []*TopologyNode
}

// Topology returns a TopologyInfo struct that describes the system topology of
// the host hardware
func Topology(opts ...*WithOption) (*TopologyInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &TopologyInfo{}
	if err := ctx.topologyFillInfo(info); err != nil {
		return nil, err
	}
	for _, node := range info.Nodes {
		sort.Sort(SortByMemoryCacheLevelTypeFirstProcessor(node.Caches))
	}
	return info, nil
}

func (i *TopologyInfo) String() string {
	archStr := "SMP"
	if i.Architecture == NUMA {
		archStr = "NUMA"
	}
	res := fmt.Sprintf(
		"topology %s (%d nodes)",
		archStr,
		len(i.Nodes),
	)
	return res
}
