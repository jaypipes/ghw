package ghw

import (
    "fmt"
    "math"
)

type Disk struct {
    Name string
    SizeBytes uint64
    SectorSize uint64
    BusType string
    Vendor string
    SerialNumber string
    Partitions []*Partition
}

type Partition struct {
    Disk *Disk
    Name string
    Uuid string
    Label string
    MountPoint string
    SizeBytes uint64
    Type string
    IsReadOnly bool
}

type BlockInfo struct {
    TotalPhysicalBytes uint64
    Disks []*Disk
    Partitions []*Partition
}

func Block() (*BlockInfo, error) {
    info := &BlockInfo{}
    err := blockFillInfo(info)
    if err != nil {
        return nil, err
    }
    return info, nil
}

func (i *BlockInfo) String() string {
    tpbs := "unknown"
    if i.TotalPhysicalBytes > 0 {
        tpb := i.TotalPhysicalBytes
        unit, unitStr := unitWithString(int64(tpb))
        tpb = uint64(math.Ceil(float64(i.TotalPhysicalBytes) / float64(unit)))
        tpbs = fmt.Sprintf("%d%s", tpb, unitStr)
    }
    return fmt.Sprintf("block storage (%s physical)", tpbs)
}

