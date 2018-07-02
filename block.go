//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"math"
)

type Disk struct {
	Name            string
	SizeBytes       uint64
	SectorSizeBytes uint64
	BusType         string
	Vendor          string
	SerialNumber    string
	Partitions      []*Partition
}

type Partition struct {
	Disk       *Disk
	Name       string
	Label      string
	MountPoint string
	SizeBytes  uint64
	Type       string
	IsReadOnly bool
}

type BlockInfo struct {
	TotalPhysicalBytes uint64
	Disks              []*Disk
	Partitions         []*Partition
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
		tpb = uint64(math.Ceil(float64(tpb) / float64(unit)))
		tpbs = fmt.Sprintf("%d%s", tpb, unitStr)
	}
	dplural := "disks"
	if len(i.Disks) == 1 {
		dplural = "disk"
	}
	return fmt.Sprintf("block storage (%d %s, %s physical storage)",
		len(i.Disks), dplural, tpbs)
}

func (d *Disk) String() string {
	vendor := ""
	if d.Vendor != "" {
		vendor = "  " + d.Vendor
	}
	serial := ""
	if d.SerialNumber != "" {
		serial = " - SN #" + d.SerialNumber
	}
	sizeStr := "unknown"
	if d.SizeBytes > 0 {
		size := d.SizeBytes
		unit, unitStr := unitWithString(int64(size))
		size = uint64(math.Ceil(float64(size) / float64(unit)))
		sizeStr = fmt.Sprintf("%d%s", size, unitStr)
	}
	return fmt.Sprintf(
		"/dev/%s (%s) [%s]%s%s",
		d.Name,
		sizeStr,
		d.BusType,
		vendor,
		serial,
	)
}

func (p *Partition) String() string {
	typeStr := ""
	if p.Type != "" {
		typeStr = fmt.Sprintf("[%s]", p.Type)
	}
	mountStr := ""
	if p.MountPoint != "" {
		mountStr = fmt.Sprintf(" mounted@%s", p.MountPoint)
	}
	sizeStr := "unknown"
	if p.SizeBytes > 0 {
		size := p.SizeBytes
		unit, unitStr := unitWithString(int64(size))
		size = uint64(math.Ceil(float64(size) / float64(unit)))
		sizeStr = fmt.Sprintf("%d%s", size, unitStr)
	}
	return fmt.Sprintf(
		"/dev/%s (%s) %s%s",
		p.Name,
		sizeStr,
		typeStr,
		mountStr,
	)
}
