//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package block

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/unitutil"
	"github.com/jaypipes/ghw/pkg/util"
)

// DriveType describes the general category of drive device
type DriveType int

const (
	// DriveTypeUnknown means we could not determine the drive type of the disk
	DriveTypeUnknown DriveType = iota
	// DriveTypeHDD indicates a hard disk drive
	DriveTypeHDD
	// DriveTypeFDD indicates a floppy disk drive
	DriveTypeFDD
	// DriveTypeODD indicates an optical disk drive
	DriveTypeODD
	// DriveTypeSSD indicates a solid-state drive
	DriveTypeSSD
	// DriveTypeVirtual indicates a virtual drive i.e. loop devices
	DriveTypeVirtual
)

const (
	// DEPRECATED: Please use DriveTypeUnknown
	DRIVE_TYPE_UNKNOWN = DriveTypeUnknown
	// DEPRECATED: Please use DriveTypeHDD
	DRIVE_TYPE_HDD = DriveTypeHDD
	// DEPRECATED: Please use DriveTypeFDD
	DRIVE_TYPE_FDD = DriveTypeFDD
	// DEPRECATED: Please use DriveTypeODD
	DRIVE_TYPE_ODD = DriveTypeODD
	// DEPRECATED: Please use DriveTypeSSD
	DRIVE_TYPE_SSD = DriveTypeSSD
	// DEPRECATED: Please use DriveTypeVirtual
	DRIVE_TYPE_VIRTUAL = DriveTypeVirtual
)

var (
	driveTypeString = map[DriveType]string{
		DriveTypeUnknown: "Unknown",
		DriveTypeHDD:     "HDD",
		DriveTypeFDD:     "FDD",
		DriveTypeODD:     "ODD",
		DriveTypeSSD:     "SSD",
		DriveTypeVirtual: "virtual",
	}

	// NOTE(fromani): the keys are all lowercase and do not match
	// the keys in the opposite table `driveTypeString`.
	// This is done because of the choice we made in
	// DriveType::MarshalJSON.
	// We use this table only in UnmarshalJSON, so it should be OK.
	stringDriveType = map[string]DriveType{
		"unknown": DriveTypeUnknown,
		"hdd":     DriveTypeHDD,
		"fdd":     DriveTypeFDD,
		"odd":     DriveTypeODD,
		"ssd":     DriveTypeSSD,
		"virtual": DriveTypeVirtual,
	}
)

func (dt DriveType) String() string {
	return driveTypeString[dt]
}

// NOTE(jaypipes): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (dt DriveType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.ToLower(dt.String()))), nil
}

func (dt *DriveType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	key := strings.ToLower(s)
	val, ok := stringDriveType[key]
	if !ok {
		return fmt.Errorf("unknown drive type: %q", key)
	}
	*dt = val
	return nil
}

// StorageController is a category of block storage controller/driver. It
// represents more of the physical hardware interface than the storage
// protocol, which represents more of the software interface.
//
// See discussion on https://github.com/jaypipes/ghw/issues/117
type StorageController int

const (
	// StorageControllerUnknown indicates we could not determine the storage
	// controller for the disk
	StorageControllerUnknown StorageController = iota
	// StorageControllerIDE indicates a Integrated Drive Electronics (IDE)
	// controller
	StorageControllerIDE
	// StorageControllerSCSI indicates a  Small computer system interface
	// (SCSI) controller
	StorageControllerSCSI
	// StorageControllerNVMe indicates a Non-volatile Memory Express (NVMe)
	// controller
	StorageControllerNVMe
	// StorageControllerVirtIO indicates a virtualized storage
	// controller/driver
	StorageControllerVirtIO
	// StorageControllerMMC indicates a Multi-media controller (used for mobile
	// phone storage devices)
	StorageControllerMMC
	// StorageControllerLoop indicates a loopback storage controller
	StorageControllerLoop
)

const (
	// DEPRECATED: Please use StorageControllerUnknown
	STORAGE_CONTROLLER_UNKNOWN = StorageControllerUnknown
	// DEPRECATED: Please use StorageControllerIDE
	STORAGE_CONTROLLER_IDE = StorageControllerIDE
	// DEPRECATED: Please use StorageControllerSCSI
	STORAGE_CONTROLLER_SCSI = StorageControllerSCSI
	// DEPRECATED: Please use StorageControllerNVMe
	STORAGE_CONTROLLER_NVME = StorageControllerNVMe
	// DEPRECATED: Please use StorageControllerVirtIO
	STORAGE_CONTROLLER_VIRTIO = StorageControllerVirtIO
	// DEPRECATED: Please use StorageControllerMMC
	STORAGE_CONTROLLER_MMC = StorageControllerMMC
	// DEPRECATED: Please use StorageControllerLoop
	STORAGE_CONTROLLER_LOOP = StorageControllerLoop
)

var (
	storageControllerString = map[StorageController]string{
		StorageControllerUnknown: "Unknown",
		StorageControllerIDE:     "IDE",
		StorageControllerSCSI:    "SCSI",
		StorageControllerNVMe:    "NVMe",
		StorageControllerVirtIO:  "virtio",
		StorageControllerMMC:     "MMC",
		StorageControllerLoop:    "loop",
	}

	// NOTE(fromani): the keys are all lowercase and do not match
	// the keys in the opposite table `storageControllerString`.
	// This is done/ because of the choice we made in
	// StorageController::MarshalJSON.
	// We use this table only in UnmarshalJSON, so it should be OK.
	stringStorageController = map[string]StorageController{
		"unknown": StorageControllerUnknown,
		"ide":     StorageControllerIDE,
		"scsi":    StorageControllerSCSI,
		"nvme":    StorageControllerNVMe,
		"virtio":  StorageControllerVirtIO,
		"mmc":     StorageControllerMMC,
		"loop":    StorageControllerLoop,
	}
)

func (sc StorageController) String() string {
	return storageControllerString[sc]
}

func (sc *StorageController) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	key := strings.ToLower(s)
	val, ok := stringStorageController[key]
	if !ok {
		return fmt.Errorf("unknown storage controller: %q", key)
	}
	*sc = val
	return nil
}

// NOTE(jaypipes): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (sc StorageController) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.ToLower(sc.String()))), nil
}

// Disk describes a single disk drive on the host system. Disk drives provide
// raw block storage resources.
type Disk struct {
	// Name contains a short name for the disk, e.g. `sda`
	Name string `json:"name"`
	// SizeBytes contains the total amount of storage, in bytes, for this disk
	SizeBytes uint64 `json:"size_bytes"`
	// PhysicalBlockSizeBytes is the size, in bytes, of the physical blocks in
	// this disk. This is typically the minimum amount of data that can be
	// written to a disk in a single write operation.
	PhysicalBlockSizeBytes uint64 `json:"physical_block_size_bytes"`
	// DriveType is the category of disk drive for this disk.
	DriveType DriveType `json:"drive_type"`
	// IsRemovable indicates if the disk drive is removable.
	IsRemovable bool `json:"removable"`
	// StorageController is the category of storage controller used by the
	// disk.
	StorageController StorageController `json:"storage_controller"`
	// BusPath is the filepath to the bus for this disk.
	BusPath string `json:"bus_path"`
	// NUMANodeID contains the numeric index (0-based) of the NUMA Node this
	// disk is affined to, or -1 if the host system is non-NUMA.
	// TODO(jaypipes): Convert this to a TopologyNode struct pointer and then
	// add to serialized output as "numa_node,omitempty"
	NUMANodeID int `json:"-"`
	// Vendor is the manufacturer of the disk.
	Vendor string `json:"vendor"`
	// Model is the model number of the disk.
	Model string `json:"model"`
	// SerialNumber is the serial number of the disk.
	SerialNumber string `json:"serial_number"`
	// WWN is the World-wide Name of the disk.
	// See: https://en.wikipedia.org/wiki/World_Wide_Name
	WWN string `json:"wwn"`
	// WWNNoExtension is the World-wide Name of the disk with any vendor
	// extensions excluded.
	// See: https://en.wikipedia.org/wiki/World_Wide_Name
	WWNNoExtension string `json:"wwnNoExtension"`
	// Partitions contains an array of pointers to `Partition` structs, one for
	// each partition on the disk.
	Partitions []*Partition `json:"partitions"`
	// TODO(jaypipes): Add PCI field for accessing PCI device information
	// PCI *PCIDevice `json:"pci"`
}

// Partition describes a logical division of a Disk.
type Partition struct {
	// Disk is a pointer to the `Disk` struct that houses this partition.
	Disk *Disk `json:"-"`
	// Name is the system name given to the partition, e.g. "sda1".
	Name string `json:"name"`
	// Label is the human-readable label given to the partition. On Linux, this
	// is derived from the `ID_PART_ENTRY_NAME` udev entry.
	Label string `json:"label"`
	// MountPoint is the path where this partition is mounted.
	MountPoint string `json:"mount_point"`
	// SizeBytes contains the total amount of storage, in bytes, this partition
	// can consume.
	SizeBytes uint64 `json:"size_bytes"`
	// Type contains the type of the partition.
	Type string `json:"type"`
	// IsReadOnly indicates if the partition is marked read-only.
	IsReadOnly bool `json:"read_only"`
	// UUID is the universally-unique identifier (UUID) for the partition.
	// This will be volume UUID on Darwin, PartUUID on linux, empty on Windows.
	UUID string `json:"uuid"`
	// FilesystemLabel is the label of the filesystem contained on the
	// partition. On Linux, this is derived from the `ID_FS_NAME` udev entry.
	FilesystemLabel string `json:"filesystem_label"`
}

// Info describes all disk drives and partitions in the host system.
type Info struct {
	ctx *context.Context
	// TotalSizeBytes contains the total amount of storage, in bytes, on the
	// host system.
	TotalSizeBytes uint64 `json:"total_size_bytes"`
	// DEPRECATED: Please use TotalSizeBytes
	TotalPhysicalBytes uint64 `json:"-"`
	// Disks contains an array of pointers to `Disk` structs, one for each disk
	// drive on the host system.
	Disks []*Disk `json:"disks"`
	// Partitions contains an array of pointers to `Partition` structs, one for
	// each partition on any disk drive on the host system.
	Partitions []*Partition `json:"-"`
}

// New returns a pointer to an Info struct that describes the block storage
// resources of the host system.
func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}
	return info, nil
}

// String returns a short string indicating important information about the
// block storage on the host system.
func (i *Info) String() string {
	tpbs := util.UNKNOWN
	if i.TotalPhysicalBytes > 0 {
		tpb := i.TotalPhysicalBytes
		unit, unitStr := unitutil.AmountString(int64(tpb))
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

// String returns a short string indicating important information about the
// disk.
func (d *Disk) String() string {
	sizeStr := util.UNKNOWN
	if d.SizeBytes > 0 {
		size := d.SizeBytes
		unit, unitStr := unitutil.AmountString(int64(size))
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
	if d.Model != util.UNKNOWN {
		model = " model=" + d.Model
	}
	serial := ""
	if d.SerialNumber != util.UNKNOWN {
		serial = " serial=" + d.SerialNumber
	}
	wwn := ""
	if d.WWN != util.UNKNOWN {
		wwn = " WWN=" + d.WWN
	}
	removable := ""
	if d.IsRemovable {
		removable = " removable=true"
	}
	return fmt.Sprintf(
		"%s %s (%s) %s [@%s%s]%s",
		d.Name,
		d.DriveType.String(),
		sizeStr,
		d.StorageController.String(),
		d.BusPath,
		atNode,
		util.ConcatStrings(
			vendor,
			model,
			serial,
			wwn,
			removable,
		),
	)
}

// String returns a short string indicating important information about the
// partition.
func (p *Partition) String() string {
	typeStr := ""
	if p.Type != "" {
		typeStr = fmt.Sprintf("[%s]", p.Type)
	}
	mountStr := ""
	if p.MountPoint != "" {
		mountStr = fmt.Sprintf(" mounted@%s", p.MountPoint)
	}
	sizeStr := util.UNKNOWN
	if p.SizeBytes > 0 {
		size := p.SizeBytes
		unit, unitStr := unitutil.AmountString(int64(size))
		size = uint64(math.Ceil(float64(size) / float64(unit)))
		sizeStr = fmt.Sprintf("%d%s", size, unitStr)
	}
	return fmt.Sprintf(
		"%s (%s) %s%s",
		p.Name,
		sizeStr,
		typeStr,
		mountStr,
	)
}

// simple private struct used to encapsulate block information in a top-level
// "block" YAML/JSON map/object key
type blockPrinter struct {
	Info *Info `json:"block" yaml:"block"`
}

// YAMLString returns a string with the block information formatted as YAML
// under a top-level "block:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, blockPrinter{i})
}

// JSONString returns a string with the block information formatted as JSON
// under a top-level "block:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, blockPrinter{i}, indent)
}
