package ghw

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Architecture int

const (
	SMP Architecture = iota
	NUMA
)

type NodeId uint32

type Node struct {
	Id     NodeId
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

type MemoryCacheType int

const (
	UNIFIED MemoryCacheType = iota
	INSTRUCTION
	DATA
)

type ByCacheLevel []*MemoryCache

func (a ByCacheLevel) Len() int      { return len(a) }
func (a ByCacheLevel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCacheLevel) Less(i, j int) bool {
	if a[i].Level < a[j].Level {
		return true
	} else if a[i].Level == a[j].Level {
		return a[i].Type < a[j].Type
	}
	return false
}

type MemoryCache struct {
	Level     uint8
	Type      MemoryCacheType
	SizeBytes uint64
	// The set of logical processors (hardware threads) that have access to the
	// cache
	LogicalProcessors []ProcessorId
}

func (c *MemoryCache) String() string {
	sizeKb := c.SizeBytes / uint64(KB)
	typeStr := ""
	if c.Type == INSTRUCTION {
		typeStr = "i"
	} else if c.Type == DATA {
		typeStr = "d"
	}
	cacheIdStr := fmt.Sprintf("L%d%s", c.Level, typeStr)
	processorMapStr := ""
	if c.LogicalProcessors != nil {
		lpStrings := make([]string, len(c.LogicalProcessors))
		for x, lpid := range c.LogicalProcessors {
			lpStrings[x] = strconv.Itoa(int(lpid))
		}
		processorMapStr = " shared with logical processors: " + strings.Join(lpStrings, ",")
	}
	return fmt.Sprintf(
		"%s cache (%d KB)%s",
		cacheIdStr,
		sizeKb,
		processorMapStr,
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
		sort.Sort(ByCacheLevel(node.Caches))
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
