package ghw

type HostInfo struct {
    Memory *MemoryInfo
    Block *BlockInfo
    CPU *CPUInfo
    Topology *TopologyInfo
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
    return info, nil
}
