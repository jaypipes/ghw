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

type Node struct {
	Id     uint32
	Cores  []*ProcessorCore
	Caches []*MemoryCache
}

func (n *Node) String() string {
	return fmt.Sprintf(
		"node #%d (%d cores)",
		n.Id,
		len(n.Cores),
	)
}

type TopologyInfo struct {
	Architecture Architecture
	Nodes        []*Node
}

func Topology() (*TopologyInfo, error) {
	info := &TopologyInfo{}
	err := topologyFillInfo(info)
	for _, node := range info.Nodes {
		sort.Sort(SortByMemoryCacheLevel(node.Caches))
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
