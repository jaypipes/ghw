package ghw

type MemoryInfo struct {
    TotalPhysicalBytes uint64
    TotalUsageBytes uint64
    SupportedPageSizes []uint64
}

func NewMemoryInfo() (*MemoryInfo, error) {
    info := &MemoryInfo{}
    err := memFillInfo(info)
    if err != nil {
        return nil, err
    }
    return info, nil
}
