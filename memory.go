package ghw

import (
    "fmt"
    "math"
)

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

func (i *MemoryInfo) String() string {
    tpbs := "unknown"
    if i.TotalPhysicalBytes > 0 {
        tpb := i.TotalPhysicalBytes
        unit, unitStr := unitWithString(tpb)
        tpb = int64(math.Ceil(float64(i.TotalPhysicalBytes) / float64(unit)))
        tpbs = fmt.Sprintf("%d%s", tpb, unitStr)
    }
    return fmt.Sprintf("memory (%s physical)", tpbs)
}
