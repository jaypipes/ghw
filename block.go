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
	Name                   string  `json:"name"`
	SizeBytes              uint64  `json:"size_bytes"`
	PhysicalBlockSizeBytes uint64  `json:"physical_block_size_bytes"`
	BusType                BusType `json:"bus_type"`
	BusPath                string  `json:"bus_path"`
	// TODO(jaypipes): Convert this to a TopologyNode struct pointer and then
	// add to serialized output as "numa_node,omitempty"
	NUMANodeID   int          `json:"-"`
	Vendor       string       `json:"vendor"`
	Model        string       `json:"model"`
	SerialNumber string       `json:"serial_number"`
	WWN          string       `json:"wwn"`
	Partitions   []*Partition `json:"partitions"`
}

// Partition describes a logical division of a Disk.
type Partition struct {
	Disk       *Disk  `json:"-"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	MountPoint string `json:"mount_point"`
	SizeBytes  uint64 `json:"size_bytes"`
	Type       string `json:"type"`
	IsReadOnly bool   `json:"read_only"`
}

// BlockInfo describes all disk drives and partitions in the host system.
type BlockInfo struct {
	// TODO(jaypipes): Deprecate this field and replace with TotalSizeBytes
	TotalPhysicalBytes uint64       `json:"total_size_bytes"`
	Disks              []*Disk      `json:"disks"`
	Partitions         []*Partition `json:"-"`
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
		d.BusType.String(),
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

// simple private struct used to encapsulate block information in a top-level
// "block" YAML/JSON map/object key
type blockPrinter struct {
	Info *BlockInfo `json:"block" yaml:"block"`
}

// YAMLString returns a string with the block information formatted as YAML
// under a top-level "block:" key
func (i *BlockInfo) YAMLString() string {
	return safeYAML(blockPrinter{i})
}

// JSONString returns a string with the block information formatted as JSON
// under a top-level "block:" key
func (i *BlockInfo) JSONString(indent bool) string {
	return safeJSON(blockPrinter{i}, indent)
}
