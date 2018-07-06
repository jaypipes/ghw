// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	PathMtab        = "/etc/mtab"
	PathSysBlock    = "/sys/block"
	PathDevDiskById = "/dev/disk/by-id"
)

var RegexNVMeDev = regexp.MustCompile(`^nvme\d+n\d+$`)
var RegexNVMePart = regexp.MustCompile(`^(nvme\d+n\d+)p\d+$`)

func blockFillInfo(info *BlockInfo) error {
	info.Disks = Disks()
	var tpb uint64
	for _, d := range info.Disks {
		tpb += d.SizeBytes
	}
	info.TotalPhysicalBytes = tpb
	return nil
}

func DiskSectorSizeBytes(disk string) uint64 {
	// We can find the sector size in Linux by looking at the
	// /sys/block/$DEVICE/queue/physical_block_size file in sysfs
	path := filepath.Join(PathSysBlock, disk, "queue", "physical_block_size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i)
}

func DiskSizeBytes(disk string) uint64 {
	// We can find the number of 512-byte sectors by examining the contents of
	// /sys/block/$DEVICE/size and calculate the physical bytes accordingly.
	path := filepath.Join(PathSysBlock, disk, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	ss := DiskSectorSizeBytes(disk)
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * ss
}

func DiskVendor(disk string) string {
	// In Linux, the vendor for a disk device is found in the
	// /sys/block/$DEVICE/device/vendor file in sysfs
	path := filepath.Join(PathSysBlock, disk, "device", "vendor")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(contents))
}

func DiskSerialNumber(disk string) string {
	// Finding the serial number of a disk without root privileges in Linux is
	// a little tricky. The /dev/disk/by-id directory contains a bunch of
	// symbolic links to disk devices and partitions. The serial number is
	// embedded as part of the symbolic link. For example, on my system, the
	// primary SCSI disk (/dev/sda) is represented as a symbolic link named
	// /dev/disk/by-id/scsi-3600508e000000000f8253aac9a1abd0c. The serial
	// number is 3600508e000000000f8253aac9a1abd0c.
	//
	// Some SATA drives (or rather, disk drive vendors) use inconsistent ways
	// of putting the serial numbers of the disks in this symbolic link name.
	// For example, here are two SATA drive identifiers (examples come from
	// @antylama on GH Issue #19):
	//
	// /dev/disk/by-id/ata-AXIOMTEK_Corp.-FSA032G300MW5T-H_BCA11704240020001
	//
	// in the above identifier, "BCA11704240020001" is the drive serial number.
	// The vendor name along with what appears to be a vendor model name
	// (FSA032G300MW5T-H) are also included in the symbolic link name.
	//
	// /dev/disk/by-id/ata-WDC_WD10JFCX-68N6GN0_WD-WX31A76R3KFS
	//
	// in the above identifier, the serial number of the disk is actually
	// WD-WX31A76R3KFS, not WX31A76R3KFS. Go figure...
	path := filepath.Join(PathDevDiskById)
	links, err := ioutil.ReadDir(path)
	if err != nil {
		return "unknown"
	}
	for _, link := range links {
		lname := link.Name()
		lpath := filepath.Join(PathDevDiskById, lname)
		dest, err := os.Readlink(lpath)
		if err != nil {
			continue
		}
		dest = filepath.Base(dest)
		if dest != disk {
			continue
		}
		pos := strings.LastIndex(lname, "_")
		if pos < 0 {
			pos = strings.Index(lname, "-")
		}
		if pos >= 0 {
			return lname[pos+1:]
		}
	}
	return "unknown"
}

func DiskPartitions(disk string) []*Partition {
	out := make([]*Partition, 0)
	path := filepath.Join(PathSysBlock, disk)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, file := range files {
		fname := file.Name()
		if !strings.HasPrefix(fname, disk) {
			continue
		}
		size := PartitionSizeBytes(fname)
		mp, pt, ro := PartitionInfo(fname)
		p := &Partition{
			Name:       fname,
			SizeBytes:  size,
			MountPoint: mp,
			Type:       pt,
			IsReadOnly: ro,
		}
		out = append(out, p)
	}
	return out
}

func Disks() []*Disk {
	// In Linux, we could use the fdisk, lshw or blockdev commands to list disk
	// information, however all of these utilities require root privileges to
	// run. We can get all of this information by examining the /sys/block
	// and /sys/class/block files
	disks := make([]*Disk, 0)
	files, err := ioutil.ReadDir(PathSysBlock)
	if err != nil {
		return nil
	}
	for _, file := range files {
		dname := file.Name()

		var busType string
		if strings.HasPrefix(dname, "sd") {
			busType = "SCSI"
		} else if strings.HasPrefix(dname, "hd") {
			busType = "IDE"
		} else if RegexNVMeDev.MatchString(dname) {
			busType = "NVMe"
		}
		if busType == "" {
			continue
		}

		size := DiskSizeBytes(dname)
		ss := DiskSectorSizeBytes(dname)
		vendor := DiskVendor(dname)
		serialNo := DiskSerialNumber(dname)

		d := &Disk{
			Name:            dname,
			SizeBytes:       size,
			SectorSizeBytes: ss,
			BusType:         busType,
			Vendor:          vendor,
			SerialNumber:    serialNo,
		}

		parts := DiskPartitions(dname)
		// Map this Disk object into the Partition...
		for _, part := range parts {
			part.Disk = d
		}
		d.Partitions = parts

		disks = append(disks, d)
	}

	return disks
}

func PartitionSizeBytes(part string) uint64 {
	// Allow calling PartitionSize with either the full partition name
	// "/dev/sda1" or just "sda1"
	if strings.HasPrefix(part, "/dev") {
		part = part[4:len(part)]
	}
	disk := part[0:3]
	if m := RegexNVMePart.FindStringSubmatch(part); len(m) > 0 {
		disk = m[1]
	}
	path := filepath.Join(PathSysBlock, disk, part, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	ss := DiskSectorSizeBytes(disk)
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * ss
}

// Given a full or short partition name, returns the mount point, the type of
// the partition and whether it's readonly
func PartitionInfo(part string) (string, string, bool) {
	// Allow calling PartitionInfo with either the full partition name
	// "/dev/sda1" or just "sda1"
	if !strings.HasPrefix(part, "/dev") {
		part = "/dev/" + part
	}

	// /etc/mtab entries for mounted partitions look like this:
	// /dev/sda6 / ext4 rw,relatime,errors=remount-ro,data=ordered 0 0
	var r io.ReadCloser
	r, err := os.Open(PathMtab)
	if err != nil {
		return "", "", true
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] != '/' {
			continue
		}
		fields := strings.Fields(line)
		if fields[0] != part {
			continue
		}
		opts := strings.Split(fields[3], ",")
		ro := true
		for _, opt := range opts {
			if opt == "rw" {
				ro = false
				break
			}
		}

		return fields[1], fields[2], ro
	}
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
