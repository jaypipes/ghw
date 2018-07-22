// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PATH_DEVICES_SYSTEM_NODE = "/sys/devices/system/node/"
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

	files, err := ioutil.ReadDir(PATH_DEVICES_SYSTEM_NODE)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "node") {
			continue
		}
		node := &Node{}
		nodeId, err := strconv.Atoi(filename[4:])
		if err != nil {
			return nil, err
		}
		node.Id = uint32(nodeId)
		cores, err := coresForNode(node.Id)
		if err != nil {
			return nil, err
		}
		node.Cores = cores
		caches, err := cachesForNode(node.Id)
		if err != nil {
			return nil, err
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func coresForNode(nodeId uint32) ([]*ProcessorCore, error) {
	// The /sys/devices/system/node/nodeX directory contains a subdirectory
	// called 'cpuX' for each logical processor assigned to the node. Each of
	// those subdirectories contains a topology subdirectory which has a
	// core_id file that indicates the 0-based identifier of the physical core
	// the logical processor (hardware thread) is on.
	path := filepath.Join(
		PATH_DEVICES_SYSTEM_NODE,
		fmt.Sprintf("node%d", nodeId),
	)
	cores := make([]*ProcessorCore, 0)

	findCoreById := func(cid ProcessorId) *ProcessorCore {
		for _, c := range cores {
			if c.Id == cid {
				return c
			}
		}

		c := &ProcessorCore{
			Id:                cid,
			Index:             len(cores),
			LogicalProcessors: make([]ProcessorId, 0),
		}
		cores = append(cores, c)
		return c
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "cpu") {
			continue
		}
		if filename == "cpumap" || filename == "cpulist" {
			// There are two files in the node directory that start with 'cpu'
			// but are not subdirectories ('cpulist' and 'cpumap'). Ignore
			// these files.
			continue
		}
		// Grab the logical processor ID by cutting the integer from the
		// /sys/devices/system/node/nodeX/cpuX filename
		cpuPath := filepath.Join(path, filename)
		lpId, _ := strconv.Atoi(filename[3:])
		coreIdPath := filepath.Join(cpuPath, "topology", "core_id")
		coreIdContents, err := ioutil.ReadFile(coreIdPath)
		if err != nil {
			continue
		}
		// coreIdContents is a []byte with the last byte being a newline rune
		coreIdStr := string(coreIdContents[:len(coreIdContents)-1])
		coreIdInt, _ := strconv.Atoi(coreIdStr)
		coreId := ProcessorId(coreIdInt)

		core := findCoreById(coreId)
		core.LogicalProcessors = append(
			core.LogicalProcessors,
			ProcessorId(lpId),
		)
	}

	for _, c := range cores {
		c.NumThreads = uint32(len(c.LogicalProcessors))
	}

	return cores, nil
}
