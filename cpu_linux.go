// +build linux

package ghw

import (
	"bufio"
	"os"
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
	setPhysicalIds := make(map[ProcessorId]bool, 0)
	for _, attrs := range procAttrs {
		pid, err := strconv.Atoi(attrs["physical id"])
		if err != nil {
			continue
		}
		setPhysicalIds[ProcessorId(pid)] = true
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
			if pid == ProcessorId(lppid) {
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
				if c.Id == ProcessorId(coreId) {
					c.LogicalProcessors = append(
						c.LogicalProcessors,
						ProcessorId(lpid),
					)
					c.NumThreads = uint32(len(c.LogicalProcessors))
					core = c
				}
			}
			if core == nil {
				coreLps := make([]ProcessorId, 1)
				coreLps[0] = ProcessorId(lpid)
				core = &ProcessorCore{
					Id:                ProcessorId(coreId),
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
