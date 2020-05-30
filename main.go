//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
)

const (
	UNKNOWN = "unknown"
)

// HostInfo is a wrapper struct containing information about the host system's
// memory, block storage, CPU, etc
type HostInfo struct {
	Memory    *MemoryInfo    `json:"memory"`
	Block     *BlockInfo     `json:"block"`
	CPU       *CPUInfo       `json:"cpu"`
	Topology  *TopologyInfo  `json:"topology"`
	Network   *NetworkInfo   `json:"network"`
	GPU       *GPUInfo       `json:"gpu"`
	Chassis   *ChassisInfo   `json:"chassis"`
	BIOS      *BIOSInfo      `json:"bios"`
	Baseboard *BaseboardInfo `json:"baseboard"`
	Product   *ProductInfo   `json:"product"`
}

// Host returns a pointer to a HostInfo struct that contains fields with
// information about the host system's CPU, memory, network devices, etc
func Host(opts ...*WithOption) (*HostInfo, error) {
	ctx := context.New(opts...)
	mem := &MemoryInfo{}
	if err := memFillInfo(ctx, mem); err != nil {
		return nil, err
	}
	block := &BlockInfo{}
	if err := blockFillInfo(ctx, block); err != nil {
		return nil, err
	}
	cpu := &CPUInfo{}
	if err := cpuFillInfo(ctx, cpu); err != nil {
		return nil, err
	}
	topology := &TopologyInfo{}
	if err := topologyFillInfo(ctx, topology); err != nil {
		return nil, err
	}
	net := &NetworkInfo{}
	if err := netFillInfo(ctx, net); err != nil {
		return nil, err
	}
	gpu := &GPUInfo{}
	if err := gpuFillInfo(ctx, gpu); err != nil {
		return nil, err
	}
	chassis := &ChassisInfo{}
	if err := chassisFillInfo(ctx, chassis); err != nil {
		return nil, err
	}
	bios := &BIOSInfo{}
	if err := biosFillInfo(ctx, bios); err != nil {
		return nil, err
	}
	baseboard := &BaseboardInfo{}
	if err := baseboardFillInfo(ctx, baseboard); err != nil {
		return nil, err
	}
	product := &ProductInfo{}
	if err := productFillInfo(ctx, product); err != nil {
		return nil, err
	}
	return &HostInfo{
		CPU:       cpu,
		Memory:    mem,
		Block:     block,
		Topology:  topology,
		Network:   net,
		GPU:       gpu,
		Chassis:   chassis,
		BIOS:      bios,
		Baseboard: baseboard,
		Product:   product,
	}, nil
}

// String returns a newline-separated output of the HostInfo's component
// structs' String-ified output
func (info *HostInfo) String() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		info.Block.String(),
		info.CPU.String(),
		info.GPU.String(),
		info.Memory.String(),
		info.Network.String(),
		info.Topology.String(),
		info.Chassis.String(),
		info.BIOS.String(),
		info.Baseboard.String(),
		info.Product.String(),
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
