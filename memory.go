package ghw

type MemoryInfo struct {
    TotalPhysicalBytes int64
    TotalUsageBytes int64
    // An array of sizes, in bytes, of memory pages supported by the host
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
