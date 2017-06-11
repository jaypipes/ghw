package ghw

import (
    "fmt"
)

type Architecture int

const (
    SMP Architecture = iota
    NUMA
)

type NodeId uint32

type Node struct {
    Id NodeId
}

func (n *Node) String() string {
    return fmt.Sprintf(
        "node #%d",
        n.Id,
    )
}

type TopologyInfo struct {
    Architecture Architecture
    Nodes []*Node
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
