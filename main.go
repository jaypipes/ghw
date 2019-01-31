//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "fmt"

const (
	UNKNOWN = "unknown"
)

// Concrete merged set of configuration switches that act as an execution
// context when calling internal discovery methods
type context struct {
	chroot string
}

type HostInfo struct {
	Memory   *MemoryInfo   `json:"memory"`
	Block    *BlockInfo    `json:"block"`
	CPU      *CPUInfo      `json:"cpu"`
	Topology *TopologyInfo `json:"topology"`
	Network  *NetworkInfo  `json:"network"`
	GPU      *GPUInfo      `json:"gpu"`
}

// Host returns a pointer to a HostInfo struct that contains fields with
// information about the host system's CPU, memory, network devices, etc
func Host(opts ...*WithOption) (*HostInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	mem := &MemoryInfo{}
	if err := ctx.memFillInfo(mem); err != nil {
		return nil, err
	}
	block := &BlockInfo{}
	if err := ctx.blockFillInfo(block); err != nil {
		return nil, err
	}
	cpu := &CPUInfo{}
	if err := ctx.cpuFillInfo(cpu); err != nil {
		return nil, err
	}
	topology := &TopologyInfo{}
	if err := ctx.topologyFillInfo(topology); err != nil {
		return nil, err
	}
	net := &NetworkInfo{}
	if err := ctx.netFillInfo(net); err != nil {
		return nil, err
	}
	gpu := &GPUInfo{}
	if err := ctx.gpuFillInfo(gpu); err != nil {
		return nil, err
	}
	return &HostInfo{
		CPU:      cpu,
		Memory:   mem,
		Block:    block,
		Topology: topology,
		Network:  net,
		GPU:      gpu,
	}, nil
}

func (info *HostInfo) String() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n",
		info.Block.String(),
		info.CPU.String(),
		info.GPU.String(),
		info.Memory.String(),
		info.Network.String(),
		info.Topology.String(),
	)
}

// YAMLString returns a string with the host information formatted as YAML
// under a top-level "host:" key
func (i *HostInfo) YAMLString() string {
	return safeYAML(i)
}

// JSONString returns a string with the host information formatted as JSON
// under a top-level "host:" key
func (i *HostInfo) JSONString(indent bool) string {
	return safeJSON(i, indent)
}
