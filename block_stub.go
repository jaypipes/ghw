// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func blockFillInfo(info *BlockInfo) error {
	return errors.New("blockFillInfo not implemented on " + runtime.GOOS)
}

func DiskPhysicalBlockSizeBytes(disk string) uint64 {
	return 0
}

func DiskSizeBytes(disk string) uint64 {
	return 0
}

func DiskNUMANodeID(disk string) int {
	return -1
}

func DiskVendor(disk string) string {
	return UNKNOWN
}

func DiskModel(disk string) string {
	return UNKNOWN
}

func DiskSerialNumber(disk string) string {
	return UNKNOWN
}

func DiskBusPath(disk string) string {
	return UNKNOWN
}

func DiskWWN(disk string) string {
	return UNKNOWN
}

func DiskPartitions(disk string) []*Partition {
	return nil
}

func Disks() []*Disk {
	return nil
}

func PartitionSizeBytes(part string) uint64 {
	return 0
}

func PartitionInfo(part string) (string, string, bool) {
	// full name, short name, read-only
	return "", "", true
}

func PartitionMountPoint(part string) string {
	mp, _, _ := PartitionInfo(part)
	return mp
}

func PartitionType(part string) string {
	_, pt, _ := PartitionInfo(part)
	return pt
}

func PartitionIsReadOnly(part string) bool {
	_, _, ro := PartitionInfo(part)
	return ro
}
