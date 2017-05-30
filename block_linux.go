// +build linux

package ghw

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "strconv"
)

const (
    PathSysBlock = "/sys/block"
    PathDevDiskById = "/dev/disk/by-id"
)

func blockFillInfo(info *BlockInfo) error {
    info.Disks = Disks()
    var tpb uint64
    for _, d := range info.Disks {
        tpb += d.SizeBytes
    }
    info.TotalPhysicalBytes = tpb
    return nil
}

func DiskSectorSize(disk string) uint64 {
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
    ss := DiskSectorSize(disk)
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
        parts := strings.Split(lname, "-")
        return parts[1]
    }
    return "unknown"
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
        // Hard drives start with an 's' or an 'h' (for SCSI and IDE) followed
        // by a 'd'
        if !((dname[0] == 's' || dname[0] == 'h') && dname[1] == 'd') {
            continue
        }

        busType := "SCSI"
        if dname[0] == 'h' {
            busType = "IDE"
        }
        size := DiskSizeBytes(dname)
        ss := DiskSectorSize(dname)
        vendor := DiskVendor(dname)
        serialNo := DiskSerialNumber(dname)

        d := &Disk{
            Name: dname,
            SizeBytes: size,
            SectorSize: ss,
            BusType: busType,
            Vendor: vendor,
            SerialNumber: serialNo,
        }

        disks = append(disks, d)
    }

    return disks
}
