//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/bios"
	"github.com/jaypipes/ghw/pkg/block"
	"github.com/jaypipes/ghw/pkg/chassis"
	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/memory"
	"github.com/jaypipes/ghw/pkg/net"
)

const (
	UNKNOWN = "unknown"
)

// HostInfo is a wrapper struct containing information about the host system's
// memory, block storage, CPU, etc
type HostInfo struct {
	Memory    *memory.Info   `json:"memory"`
	Block     *block.Info    `json:"block"`
	CPU       *cpu.Info      `json:"cpu"`
	Topology  *TopologyInfo  `json:"topology"`
	Network   *NetworkInfo   `json:"network"`
	GPU       *GPUInfo       `json:"gpu"`
	Chassis   *ChassisInfo   `json:"chassis"`
	BIOS      *bios.Info     `json:"bios"`
	Baseboard *BaseboardInfo `json:"baseboard"`
	Product   *ProductInfo   `json:"product"`
}

// Host returns a pointer to a HostInfo struct that contains fields with
// information about the host system's CPU, memory, network devices, etc
func Host(opts ...*WithOption) (*HostInfo, error) {
	ctx := context.New(opts...)
	memInfo, err := memory.New(opts...)
	if err != nil {
		return nil, err
	}
	blockInfo, err := block.New(opts...)
	if err != nil {
		return nil, err
	}
	cpuInfo, err := cpu.New(opts...)
	if err != nil {
		return nil, err
	}
	topology := &TopologyInfo{}
	if err := topologyFillInfo(ctx, topology); err != nil {
		return nil, err
	}
	netInfo, err := net.New(opts...)
	if err != nil {
		return nil, err
	}
	gpu := &GPUInfo{}
	if err := gpuFillInfo(ctx, gpu); err != nil {
		return nil, err
	}
	chassisInfo, err := chassis.New(opts...)
	if err != nil {
		return nil, err
	}
	biosInfo, err := bios.New(opts...)
	if err != nil {
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
		CPU:       cpuInfo,
		Memory:    memInfo,
		Block:     blockInfo,
		Topology:  topology,
		Network:   netInfo,
		GPU:       gpu,
		Chassis:   chassisInfo,
		BIOS:      biosInfo,
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
	return marshal.SafeYAML(i)
}

// JSONString returns a string with the host information formatted as JSON
// under a top-level "host:" key
func (i *HostInfo) JSONString(indent bool) string {
	return marshal.SafeJSON(i, indent)
}
