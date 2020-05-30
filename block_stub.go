// +build !linux,!darwin,!windows
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"

	"github.com/jaypipes/ghw/pkg/context"
)

func blockFillInfo(ctx *context.Context, info *BlockInfo) error {
	return errors.New("blockFillInfo not implemented on " + runtime.GOOS)
}

// DiskPhysicalBlockSizeBytes has been deprecated in 0.2. Please use the
// Disk.PhysicalBlockSizeBytes attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskPhysicalBlockSizeBytes(disk string) uint64 {
	msg := `
The DiskPhysicalBlockSizeBytes() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.PhysicalBlockSizeBytes
attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskPhysicalBlockSizeBytes(ctx, disk)
}

func diskPhysicalBlockSizeBytes(ctx *context.Context, disk string) uint64 {
	return 0
}

// DiskSizeBytes has been deprecated in 0.2. Please use the Disk.SizeBytes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskSizeBytes(disk string) uint64 {
	msg := `
The DiskSizeBytes() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.SizeBytes attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskSizeBytes(ctx, disk)
}

func diskSizeBytes(ctx *context.Context, disk string) uint64 {
	return 0
}

// DiskNUMANodeID has been deprecated in 0.2. Please use the Disk.NUMANodeID
// attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskNUMANodeID(disk string) int {
	msg := `
The DiskNUMANodeID() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.NUMANodeID attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskNUMANodeID(ctx, disk)
}

func diskNUMANodeID(ctx *context.Context, disk string) int {
	return -1
}

// DiskVendor has been deprecated in 0.2. Please use the Disk.Vendor attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskVendor(disk string) string {
	msg := `
The DiskVendor() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.Vendor attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskVendor(ctx, disk)
}

func diskVendor(ctx *context.Context, disk string) string {
	return UNKNOWN
}

// DiskModel has been deprecated in 0.2. Please use the Disk.Model attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskModel(disk string) string {
	msg := `
The DiskModel() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.Model attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskModel(ctx, disk)
}

func diskModel(ctx *context.Context, disk string) string {
	return UNKNOWN
}

// DiskSerialNumber has been deprecated in 0.2. Please use the Disk.SerialNumber attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskSerialNumber(disk string) string {
	msg := `
The DiskSerialNumber() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the Disk.SerialNumber attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskSerialNumber(ctx, disk)
}

func diskSerialNumber(ctx *context.Context, disk string) string {
	return UNKNOWN
}

// DiskBusPath has been deprecated in 0.2. Please use the Disk.BusPath attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskBusPath(disk string) string {
	msg := `
The DiskBusPath() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.BusPath attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskBusPath(ctx, disk)
}

func diskBusPath(ctx *context.Context, disk string) string {
	return UNKNOWN
}

// DiskWWN has been deprecated in 0.2. Please use the Disk.WWN attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskWWN(disk string) string {
	msg := `
The DiskWWN() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.WWN attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskWWN(ctx, disk)
}

func diskWWN(ctx *context.Context, disk string) string {
	return UNKNOWN
}

// DiskPartitions has been deprecated in 0.2. Please use the Disk.Partitions attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskPartitions(disk string) []*Partition {
	msg := `
The DiskPartitions() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the Disk.Partitions attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return diskPartitions(ctx, disk)
}

func diskPartitions(ctx *context.Context, disk string) []*Partition {
	return nil
}

// Disks has been deprecated in 0.2. Please use the BlockInfo.Disks attribute.
// TODO(jaypipes): Remove in 1.0.
func Disks() []*Disk {
	msg := `
The Disks() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the BlockInfo.Disks attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return disks(ctx)
}

func disks(ctx *context.Context) []*Disk {
	return nil
}

// PartitionSizeBytes has been deprecated in 0.2. Please use the
// Partition.SizeBytes attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionSizeBytes(part string) uint64 {
	msg := `
The PartitionSizeBytes() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.SizeBytes attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return partitionSizeBytes(ctx, part)
}

func partitionSizeBytes(ctx *context.Context, part string) uint64 {
	return 0
}

// PartitionInfo has been deprecated in 0.2. Please use the Partition struct.
// TODO(jaypipes): Remove in 1.0.
func PartitionInfo(part string) (string, string, bool) {
	msg := `
The PartitionInfo() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition struct.
`
	warn(msg)
	ctx := context.FromEnv()
	return partitionInfo(ctx, part)
}

// Given a full or short partition name, returns the mount point, the type of
// the partition and whether it's readonly
func partitionInfo(ctx *context.Context, part string) (string, string, bool) {
	// full name, short name, read-only
	return "", "", true
}

// PartitionMountPoint has been deprecated in 0.2. Please use the
// Partition.MountPoint attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionMountPoint(part string) string {
	msg := `
The PartitionMountPoint() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.MountPoint attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return partitionMountPoint(ctx, part)
}

func (ctx *context) partitionMountPoint(ctx *context.Context, part string) string {
	mp, _, _ := partitionInfo(ctx, part)
	return mp
}

// PartitionType has been deprecated in 0.2. Please use the
// Partition.Type attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionType(part string) string {
	msg := `
The PartitionType() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.Type attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return partitionType(ctx, part)
}

func (ctx *context) partitionType(ctx *context.Context, part string) string {
	_, pt, _ := partitionInfo(ctx, part)
	return pt
}

// PartitionIsReadOnly has been deprecated in 0.2. Please use the
// Partition.IsReadOnly attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionIsReadOnly(part string) bool {
	msg := `
The PartitionIsReadOnly() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.IsReadOnly attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return partitionIsReadOnly(ctx, part)
}

func partitionIsReadOnly(ctx *context.Context, part string) bool {
	_, _, ro := partitionInfo(ctx, part)
	return ro
}
