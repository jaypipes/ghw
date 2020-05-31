// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

func topologyFillInfo(ctx *context.Context, info *TopologyInfo) error {
	info.Nodes = topologyNodes(ctx)
	if len(info.Nodes) == 1 {
		info.Architecture = ARCHITECTURE_SMP
	} else {
		info.Architecture = ARCHITECTURE_NUMA
	}
	return nil
}

func topologyNodes(ctx *context.Context) []*TopologyNode {
	paths := linuxpath.New(ctx)
	nodes := make([]*TopologyNode, 0)

	files, err := ioutil.ReadDir(paths.SysDevicesSystemNode)
	if err != nil {
		util.Warn("failed to determine nodes: %s\n", err)
		return nodes
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "node") {
			continue
		}
		node := &TopologyNode{}
		nodeID, err := strconv.Atoi(filename[4:])
		if err != nil {
			util.Warn("failed to determine node ID: %s\n", err)
			return nodes
		}
		node.ID = nodeID
		cores, err := coresForNode(ctx, nodeID)
		if err != nil {
			util.Warn("failed to determine cores for node: %s\n", err)
			return nodes
		}
		node.Cores = cores
		caches, err := cachesForNode(ctx, nodeID)
		if err != nil {
			util.Warn("failed to determine caches for node: %s\n", err)
			return nodes
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes
}
