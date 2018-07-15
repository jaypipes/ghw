//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

type HostInfo struct {
	Memory   *MemoryInfo
	Block    *BlockInfo
	CPU      *CPUInfo
	Topology *TopologyInfo
	Network  *NetworkInfo
	GPU      *GPUInfo
}

func Host() (*HostInfo, error) {
	info := &HostInfo{}
	mem, err := Memory()
	if err != nil {
		return nil, err
	}
	info.Memory = mem
	block, err := Block()
	if err != nil {
		return nil, err
	}
	info.Block = block
	cpu, err := CPU()
	if err != nil {
		return nil, err
	}
	info.CPU = cpu
	topology, err := Topology()
	if err != nil {
		return nil, err
	}
	info.Topology = topology
	net, err := Network()
	if err != nil {
		return nil, err
	}
	info.Network = net
	gpu, err := GPU()
	if err != nil {
		return nil, err
	}
	info.GPU = gpu
	return info, nil
}
