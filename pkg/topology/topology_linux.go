// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package topology

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/memory"
	"github.com/jaypipes/ghw/pkg/util"
)

func (i *Info) load() error {
	i.Nodes = topologyNodes(i.ctx)
	if len(i.Nodes) == 1 {
		i.Architecture = ARCHITECTURE_SMP
	} else {
		i.Architecture = ARCHITECTURE_NUMA
	}
	return nil
}

func topologyNodes(ctx *context.Context) []*Node {
	paths := linuxpath.New(ctx)
	nodes := make([]*Node, 0)

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
		node := &Node{}
		nodeID, err := strconv.Atoi(filename[4:])
		if err != nil {
			util.Warn("failed to determine node ID: %s\n", err)
			return nodes
		}
		node.ID = nodeID
		cores, err := cpu.CoresForNode(ctx, nodeID)
		if err != nil {
			util.Warn("failed to determine cores for node: %s\n", err)
			return nodes
		}
		node.Cores = cores
		caches, err := memory.CachesForNode(ctx, nodeID)
		if err != nil {
			util.Warn("failed to determine caches for node: %s\n", err)
			return nodes
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes
}
