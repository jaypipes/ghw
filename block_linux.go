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
	LINUX_SECTOR_SIZE = 512
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

func DiskPhysicalBlockSizeBytes(disk string) uint64 {
	// We can find the sector size in Linux by looking at the
	// /sys/block/$DEVICE/queue/physical_block_size file in sysfs
	path := filepath.Join(pathSysBlock(), disk, "queue", "physical_block_size")
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
	path := filepath.Join(pathSysBlock(), disk, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * LINUX_SECTOR_SIZE
}

func DiskVendor(disk string) string {
	// In Linux, the vendor for a disk device is found in the
	// /sys/block/$DEVICE/device/vendor file in sysfs
	path := filepath.Join(pathSysBlock(), disk, "device", "vendor")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return UNKNOWN
	}
	return strings.TrimSpace(string(contents))
}

func udevInfo(disk string) (map[string]string, error) {
	// Get device major:minor numbers
	devNo, err := ioutil.ReadFile(filepath.Join(pathSysBlock(), disk, "dev"))
	if err != nil {
		return nil, err
	}

	// Look up block device in udev runtime database
	udevId := "b" + strings.TrimSpace(string(devNo))
	udevBytes, err := ioutil.ReadFile(filepath.Join(pathRunUdevData(), udevId))
	if err != nil {
		return nil, err
	}

	udevInfo := make(map[string]string)
	for _, udevLine := range strings.Split(string(udevBytes), "\n") {
		if strings.HasPrefix(udevLine, "E:") {
			if s := strings.SplitN(udevLine[2:], "=", 2); len(s) == 2 {
				udevInfo[s[0]] = s[1]
			}
		}
	}
	return udevInfo, nil
}

func DiskSerialNumber(disk string) string {
	info, err := udevInfo(disk)
	if err != nil {
		return UNKNOWN
	}

	// There are two serial number keys, ID_SERIAL and ID_SERIAL_SHORT
	// The non-_SHORT version often duplicates vendor information collected elsewhere, so use _SHORT.
	if path, ok := info["ID_SERIAL_SHORT"]; ok {
		return path
	}
	return UNKNOWN
}

func DiskBusPath(disk string) string {
	info, err := udevInfo(disk)
	if err != nil {
		return UNKNOWN
	}

	// There are two path keys, ID_PATH and ID_PATH_TAG.
	// The difference seems to be _TAG has funky characters converted to underscores.
	if path, ok := info["ID_PATH"]; ok {
		return path
	}
	return UNKNOWN
}

func DiskPartitions(disk string) []*Partition {
	out := make([]*Partition, 0)
	path := filepath.Join(pathSysBlock(), disk)
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
	files, err := ioutil.ReadDir(pathSysBlock())
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
		pbs := DiskPhysicalBlockSizeBytes(dname)
		busPath := DiskBusPath(dname)
		vendor := DiskVendor(dname)
		serialNo := DiskSerialNumber(dname)

		d := &Disk{
			Name:                   dname,
			SizeBytes:              size,
			PhysicalBlockSizeBytes: pbs,
			BusType:                busType,
			BusPath:                busPath,
			Vendor:                 vendor,
			SerialNumber:           serialNo,
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
	path := filepath.Join(pathSysBlock(), disk, part, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * LINUX_SECTOR_SIZE
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
	r, err := os.Open(pathEtcMtab())
	if err != nil {
		return "", "", true
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		entry := parseMtabEntry(line)
		if entry == nil || entry.Partition != part {
			continue
		}
		ro := true
		for _, opt := range entry.Options {
			if opt == "rw" {
				ro = false
				break
			}
		}

		return entry.Mountpoint, entry.FilesystemType, ro
	}
	return "", "", true
}

type mtabEntry struct {
	Partition      string
	Mountpoint     string
	FilesystemType string
	Options        []string
}

func parseMtabEntry(line string) *mtabEntry {
	// /etc/mtab entries for mounted partitions look like this:
	// /dev/sda6 / ext4 rw,relatime,errors=remount-ro,data=ordered 0 0
	if line[0] != '/' {
		return nil
	}
	fields := strings.Fields(line)

	if len(fields) < 4 {
		return nil
	}

	// We do some special parsing of the mountpoint, which may contain space,
	// tab and newline characters, encoded into the mtab entry line using their
	// octal-to-string representations. From the GNU mtab man pages:
	//
	//   "Therefore these characters are encoded in the files and the getmntent
	//   function takes care of the decoding while reading the entries back in.
	//   '\040' is used to encode a space character, '\011' to encode a tab
	//   character, '\012' to encode a newline character, and '\\' to encode a
	//   backslash."
	mp := fields[1]
	r := strings.NewReplacer(
		"\\011", "\t", "\\012", "\n", "\\040", " ", "\\\\", "\\",
	)
	mp = r.Replace(mp)

	res := &mtabEntry{
		Partition:      fields[0],
		Mountpoint:     mp,
		FilesystemType: fields[2],
	}
	opts := strings.Split(fields[3], ",")
	res.Options = opts
	return res
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
