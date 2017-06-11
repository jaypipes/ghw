// +build linux

package ghw

import (
    "io/ioutil"
    "strconv"
    "strings"
)

const (
    PathSysDevicesSystemNode = "/sys/devices/system/node/"
)

func topologyFillInfo(info *TopologyInfo) error {
    nodes, err := Nodes()
    if err != nil {
        return err
    }
    info.Nodes = nodes
    if len(info.Nodes) == 1 {
        info.Architecture = SMP
    } else {
        info.Architecture = NUMA
    }
    return nil
}

func Nodes() ([]*Node, error) {
    nodes := make([]*Node, 0)

    files, err := ioutil.ReadDir(PathSysDevicesSystemNode)
    if err != nil {
        return nil, err
    }
    for _, file := range files {
        filename := file.Name()
        if ! strings.HasPrefix(filename, "node") {
            continue
        }
        node := &Node{}
        nodeId, err := strconv.Atoi(filename[4:])
        if err != nil {
            return nil, err
        }
        node.Id = NodeId(nodeId)
        nodes = append(nodes, node)
    }
    return nodes, nil
}
