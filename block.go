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

// Disk describes a single disk drive on the host system. Disk drives provide
// raw block storage resources.
type Disk struct {
	Name                   string
	SizeBytes              uint64
	PhysicalBlockSizeBytes uint64
	BusType                string
	BusPath                string
	NUMANodeID             int
	Vendor                 string
	Model                  string
	SerialNumber           string
	WWN                    string
	Partitions             []*Partition
}

// Partition describes a logical division of a Disk.
type Partition struct {
	Disk       *Disk
	Name       string
	Label      string
	MountPoint string
	SizeBytes  uint64
	Type       string
	IsReadOnly bool
}

// BlockInfo describes all disk drives and partitions in the host system.
type BlockInfo struct {
	TotalPhysicalBytes uint64
	Disks              []*Disk
	Partitions         []*Partition
}

// Block returns a BlockInfo struct that describes the block storage resources
// of the host system.
func Block(opts ...*WithOption) (*BlockInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &BlockInfo{}
	if err := ctx.blockFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *BlockInfo) String() string {
	tpbs := UNKNOWN
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
	sizeStr := UNKNOWN
	if d.SizeBytes > 0 {
		size := d.SizeBytes
		unit, unitStr := unitWithString(int64(size))
		size = uint64(math.Ceil(float64(size) / float64(unit)))
		sizeStr = fmt.Sprintf("%d%s", size, unitStr)
	}
	atNode := ""
	if d.NUMANodeID >= 0 {
		atNode = fmt.Sprintf(" (node #%d)", d.NUMANodeID)
	}
	vendor := ""
	if d.Vendor != "" {
		vendor = " vendor=" + d.Vendor
	}
	model := ""
	if d.Model != UNKNOWN {
		model = " model=" + d.Model
	}
	serial := ""
	if d.SerialNumber != UNKNOWN {
		serial = " serial=" + d.SerialNumber
	}
	wwn := ""
	if d.WWN != UNKNOWN {
		wwn = " WWN=" + d.WWN
	}
	return fmt.Sprintf(
		"/dev/%s (%s) [%s @ %s%s]%s%s%s%s",
		d.Name,
		sizeStr,
		d.BusType,
		d.BusPath,
		atNode,
		vendor,
		model,
		serial,
		wwn,
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
	sizeStr := UNKNOWN
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
