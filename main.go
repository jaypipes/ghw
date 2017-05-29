package ghw

type HostInfo struct {
    Memory *MemoryInfo
}

func Host() (*HostInfo, error) {
    info := &HostInfo{}
    mem, err := Memory()
    if err != nil {
        return nil, err
    }
    info.Memory = mem
    return info, nil
}
