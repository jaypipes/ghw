// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func topologyFillInfo(info *TopologyInfo) error {
	nodes, err := topologyNodes()
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

// TopologyNodes has been deprecated in 0.2. Please use the TopologyInfo.Nodes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func TopologyNodes() ([]*TopologyNode, error) {
	msg := `
The TopologyNodes() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the TopologyInfo.Nodes attribute.
`
	warn(msg)
	return topologyNodes()
}

func topologyNodes() ([]*TopologyNode, error) {
	nodes := make([]*TopologyNode, 0)

	files, err := ioutil.ReadDir(pathSysDevicesSystemNode())
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "node") {
			continue
		}
		node := &TopologyNode{}
		nodeID, err := strconv.Atoi(filename[4:])
		if err != nil {
			return nil, err
		}
		node.ID = nodeID
		cores, err := coresForNode(nodeID)
		if err != nil {
			return nil, err
		}
		node.Cores = cores
		caches, err := cachesForNode(nodeID)
		if err != nil {
			return nil, err
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes, nil
}
