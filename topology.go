//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"sort"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/memory"
)

// TopologyNode is an abstract construct representing a collection of
// processors and various levels of memory cache that those processors share.
// In a NUMA architecture, there are multiple NUMA nodes, abstracted here as
// multiple TopologyNode structs. In an SMP architecture, a single TopologyNode
// will be available in the TopologyInfo struct and this single struct can be
// used to describe the levels of memory caching available to the single
// physical processor package's physical processor cores
type TopologyNode struct {
	ID     int                  `json:"id"`
	Cores  []*cpu.ProcessorCore `json:"cores"`
	Caches []*memory.Cache      `json:"caches"`
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
	Architecture Architecture    `json:"architecture"`
	Nodes        []*TopologyNode `json:"nodes"`
}

// Topology returns a TopologyInfo struct that describes the system topology of
// the host hardware
func Topology(opts ...*WithOption) (*TopologyInfo, error) {
	ctx := context.New(opts...)
	info := &TopologyInfo{}
	if err := topologyFillInfo(ctx, info); err != nil {
		return nil, err
	}
	for _, node := range info.Nodes {
		sort.Sort(memory.SortByCacheLevelTypeFirstProcessor(node.Caches))
	}
	return info, nil
}

func (i *TopologyInfo) String() string {
	archStr := "SMP"
	if i.Architecture == ARCHITECTURE_NUMA {
		archStr = "NUMA"
	}
	res := fmt.Sprintf(
		"topology %s (%d nodes)",
		archStr,
		len(i.Nodes),
	)
	return res
}

// simple private struct used to encapsulate topology information in a
// top-level "topology" YAML/JSON map/object key
type topologyPrinter struct {
	Info *TopologyInfo `json:"topology"`
}

// YAMLString returns a string with the topology information formatted as YAML
// under a top-level "topology:" key
func (i *TopologyInfo) YAMLString() string {
	return marshal.SafeYAML(topologyPrinter{i})
}

// JSONString returns a string with the topology information formatted as JSON
// under a top-level "topology:" key
func (i *TopologyInfo) JSONString(indent bool) string {
	return marshal.SafeJSON(topologyPrinter{i}, indent)
}
