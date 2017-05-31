package ghw

import (
    "fmt"
)

type ProcessorId uint32
type CoreId uint32

type CoreMap map[CoreId][]ProcessorId

type Processor struct {
    Id ProcessorId
    NumCores uint32
    NumThreads uint32
    Vendor string
    Model string
    CoreMap CoreMap
}

type CPUInfo struct {
    TotalCores uint32
    TotalThreads uint32
    Processors []*Processor
}

func CPU() (*CPUInfo, error) {
    info := &CPUInfo{}
    err := cpuFillInfo(info)
    if err != nil {
        return nil, err
    }
    return info, nil
}

func (i *CPUInfo) String() string {
    return fmt.Sprintf(
        "cpu (%d physical packages, %d cores, %d hardware threads)",
        len(i.Processors),
        i.TotalCores,
        i.TotalThreads,
    )
}
