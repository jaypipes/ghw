//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

const (
	UNKNOWN = "unknown"
)

type HostInfo struct {
	Memory   *MemoryInfo
	Block    *BlockInfo
	CPU      *CPUInfo
	Topology *TopologyInfo
	Network  *NetworkInfo
	GPU      *GPUInfo
}

// Host returns a pointer to a HostInfo struct that contains fields with
// information about the host system's CPU, memory, network devices, etc
func Host() (*HostInfo, error) {
	mem, err := Memory()
	if err != nil {
		return nil, err
	}
	block, err := Block()
	if err != nil {
		return nil, err
	}
	cpu, err := CPU()
	if err != nil {
		return nil, err
	}
	topology, err := Topology()
	if err != nil {
		return nil, err
	}
	net, err := Network()
	if err != nil {
		return nil, err
	}
	gpu, err := GPU()
	if err != nil {
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
