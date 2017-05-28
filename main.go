package ghw

type HostInfo struct {
    Memory *MemoryInfo
}

func NewInfo() (*HostInfo, error) {
    info := &HostInfo{}
    mem, err := NewMemoryInfo()
    if err != nil {
        return nil, err
    }
    info.Memory = mem
    return info, nil
}
