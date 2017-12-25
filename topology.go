package ghw

import (
	"fmt"
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
	return fmt.Sprintf(
		"topology %s (%d nodes)",
		archStr,
		len(i.Nodes),
	)
}
