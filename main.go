package ghw

type HostInfo struct {
    Memory *MemoryInfo
    Block *BlockInfo
    CPU *CPUInfo
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
    return info, nil
}
