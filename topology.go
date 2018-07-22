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

type Architecture int

const (
	SMP Architecture = iota
	NUMA
)

// A TopologyNode is an abstract construct representing a collection of
// processors and various levels of memory cache that those processors share.
// In a NUMA architecture, there are multiple NUMA nodes, abstracted here as
// multiple TopologyNode structs. In an SMP architecture, a single TopologyNode
// will be available in the TopologyInfo struct and this single struct can be
// used to describe the levels of memory caching available to the single
// physical processor package's physical processor cores
type TopologyNode struct {
	Id     uint32
	Cores  []*ProcessorCore
	Caches []*MemoryCache
}

func (n *TopologyNode) String() string {
	return fmt.Sprintf(
		"node #%d (%d cores)",
		n.Id,
		len(n.Cores),
	)
}

type TopologyInfo struct {
	Architecture Architecture
	Nodes        []*TopologyNode
}

func Topology() (*TopologyInfo, error) {
	info := &TopologyInfo{}
	err := topologyFillInfo(info)
	for _, node := range info.Nodes {
		sort.Sort(SortByMemoryCacheLevelTypeFirstProcessor(node.Caches))
	}
	if err != nil {
		return nil, err
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
