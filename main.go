package ghw

type HwInfo struct {
    Memory *MemoryInfo
}

func NewInfo() (*HwInfo, error) {
    info := &HwInfo{}
    mem, err := NewMemoryInfo()
    if err != nil {
        return nil, err
    }
    info.Memory = mem
    return info, nil
}
