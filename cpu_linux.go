// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PathProcCpuinfo = "/proc/cpuinfo"
)

func cpuFillInfo(info *CPUInfo) error {
	info.Processors = Processors()
	var totCores uint32
	var totThreads uint32
	for _, p := range info.Processors {
		totCores += p.NumCores
		totThreads += p.NumThreads
	}
	info.TotalCores = totCores
	info.TotalThreads = totThreads
	return nil
}

func Processors() []*Processor {
	procs := make([]*Processor, 0)

	r, err := os.Open(PathProcCpuinfo)
	if err != nil {
		return nil
	}
	defer r.Close()

	// An array of maps of attributes describing the logical processor
	procAttrs := make([]map[string]string, 0)
	curProcAttrs := make(map[string]string, 0)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// Output of /proc/cpuinfo has a blank newline to separate logical
			// processors, so here we collect up all the attributes we've
			// collected for this logical processor block
			procAttrs = append(procAttrs, curProcAttrs)
			// Reset the current set of processor attributes...
			curProcAttrs = make(map[string]string, 0)
			continue
		}
		parts := strings.Split(line, ":")
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		curProcAttrs[key] = value
	}

	// Build a set of physical processor IDs which represent the physical
	// package of the CPU
	setPhysicalIds := make(map[uint32]bool, 0)
	for _, attrs := range procAttrs {
		pid, err := strconv.Atoi(attrs["physical id"])
		if err != nil {
			continue
		}
		setPhysicalIds[uint32(pid)] = true
	}

	for pid, _ := range setPhysicalIds {
		p := &Processor{
			Id: pid,
		}
		// The indexes into the array of attribute maps for each logical
		// processor within the physical processor
		lps := make([]int, 0)
		for x, _ := range procAttrs {
			lppid, err := strconv.Atoi(procAttrs[x]["physical id"])
			if err != nil {
				continue
			}
			if pid == uint32(lppid) {
				lps = append(lps, x)
			}
		}
		first := procAttrs[lps[0]]
		p.Model = first["model name"]
		p.Vendor = first["vendor_id"]
		numCores, err := strconv.Atoi(first["cpu cores"])
		if err != nil {
			continue
		}
		p.NumCores = uint32(numCores)
		numThreads, err := strconv.Atoi(first["siblings"])
		if err != nil {
			continue
		}
		p.NumThreads = uint32(numThreads)

		// The flags field is a space-separated list of CPU capabilities
		p.Capabilities = strings.Split(first["flags"], " ")

		cores := make([]*ProcessorCore, 0)
		for _, lpidx := range lps {
			lpid, err := strconv.Atoi(procAttrs[lpidx]["processor"])
			if err != nil {
				continue
			}
			coreId, err := strconv.Atoi(procAttrs[lpidx]["core id"])
			if err != nil {
				continue
			}
			var core *ProcessorCore
			for _, c := range cores {
				if c.Id == uint32(coreId) {
					c.LogicalProcessors = append(
						c.LogicalProcessors,
						uint32(lpid),
					)
					c.NumThreads = uint32(len(c.LogicalProcessors))
					core = c
				}
			}
			if core == nil {
				coreLps := make([]uint32, 1)
				coreLps[0] = uint32(lpid)
				core = &ProcessorCore{
					Id:                uint32(coreId),
					Index:             len(cores),
					NumThreads:        1,
					LogicalProcessors: coreLps,
				}
				cores = append(cores, core)
			}
		}
		p.Cores = cores
		procs = append(procs, p)
	}
	return procs
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

	findCoreById := func(cid uint32) *ProcessorCore {
		for _, c := range cores {
			if c.Id == cid {
				return c
			}
		}

		c := &ProcessorCore{
			Id:                cid,
			Index:             len(cores),
			LogicalProcessors: make([]uint32, 0),
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
		coreId := uint32(coreIdInt)

		core := findCoreById(coreId)
		core.LogicalProcessors = append(
			core.LogicalProcessors,
			uint32(lpId),
		)
	}

	for _, c := range cores {
		c.NumThreads = uint32(len(c.LogicalProcessors))
	}

	return cores, nil
}
