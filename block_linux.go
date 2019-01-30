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
	sectorSize = 512
)

var regexNVMeDev = regexp.MustCompile(`^nvme\d+n\d+$`)
var regexNVMePart = regexp.MustCompile(`^(nvme\d+n\d+)p\d+$`)

func (ctx *context) blockFillInfo(info *BlockInfo) error {
	info.Disks = ctx.disks()
	var tpb uint64
	for _, d := range info.Disks {
		tpb += d.SizeBytes
	}
	info.TotalPhysicalBytes = tpb
	return nil
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
	ctx := contextFromEnv()
	return ctx.diskPhysicalBlockSizeBytes(disk)
}

func (ctx *context) diskPhysicalBlockSizeBytes(disk string) uint64 {
	// We can find the sector size in Linux by looking at the
	// /sys/block/$DEVICE/queue/physical_block_size file in sysfs
	path := filepath.Join(ctx.pathSysBlock(), disk, "queue", "physical_block_size")
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

// DiskSizeBytes has been deprecated in 0.2. Please use the Disk.SizeBytes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskSizeBytes(disk string) uint64 {
	msg := `
The DiskSizeBytes() function has been DEPRECATED and will be
removed in the 1.0 release of ghw. Please use the Disk.SizeBytes attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskSizeBytes(disk)
}

func (ctx *context) diskSizeBytes(disk string) uint64 {
	// We can find the number of 512-byte sectors by examining the contents of
	// /sys/block/$DEVICE/size and calculate the physical bytes accordingly.
	path := filepath.Join(ctx.pathSysBlock(), disk, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * sectorSize
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
	ctx := contextFromEnv()
	return ctx.diskNUMANodeID(disk)
}

func (ctx *context) diskNUMANodeID(disk string) int {
	link, err := os.Readlink(filepath.Join(ctx.pathSysBlock(), disk))
	if err != nil {
		return -1
	}
	for partial := link; strings.HasPrefix(partial, "../devices/"); partial = filepath.Base(partial) {
		if nodeContents, err := ioutil.ReadFile(filepath.Join(ctx.pathSysBlock(), partial, "numa_node")); err != nil {
			if nodeInt, err := strconv.Atoi(string(nodeContents)); err != nil {
				return nodeInt
			}
		}
	}
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
	ctx := contextFromEnv()
	return ctx.diskVendor(disk)
}

func (ctx *context) diskVendor(disk string) string {
	// In Linux, the vendor for a disk device is found in the
	// /sys/block/$DEVICE/device/vendor file in sysfs
	path := filepath.Join(ctx.pathSysBlock(), disk, "device", "vendor")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return UNKNOWN
	}
	return strings.TrimSpace(string(contents))
}

func (ctx *context) udevInfo(disk string) (map[string]string, error) {
	// Get device major:minor numbers
	devNo, err := ioutil.ReadFile(filepath.Join(ctx.pathSysBlock(), disk, "dev"))
	if err != nil {
		return nil, err
	}

	// Look up block device in udev runtime database
	udevID := "b" + strings.TrimSpace(string(devNo))
	udevBytes, err := ioutil.ReadFile(filepath.Join(ctx.pathRunUdevData(), udevID))
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

// DiskModel has been deprecated in 0.2. Please use the Disk.Model attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskModel(disk string) string {
	msg := `
The DiskModel() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.Model attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskModel(disk)
}

func (ctx *context) diskModel(disk string) string {
	info, err := ctx.udevInfo(disk)
	if err != nil {
		return UNKNOWN
	}

	if model, ok := info["ID_MODEL"]; ok {
		return model
	}
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
	ctx := contextFromEnv()
	return ctx.diskSerialNumber(disk)
}

func (ctx *context) diskSerialNumber(disk string) string {
	info, err := ctx.udevInfo(disk)
	if err != nil {
		return UNKNOWN
	}

	// There are two serial number keys, ID_SERIAL and ID_SERIAL_SHORT
	// The non-_SHORT version often duplicates vendor information collected elsewhere, so use _SHORT.
	if serial, ok := info["ID_SERIAL_SHORT"]; ok {
		return serial
	}
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
	ctx := contextFromEnv()
	return ctx.diskBusPath(disk)
}

func (ctx *context) diskBusPath(disk string) string {
	info, err := ctx.udevInfo(disk)
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

// DiskWWN has been deprecated in 0.2. Please use the Disk.WWN attribute.
// TODO(jaypipes): Remove in 1.0.
func DiskWWN(disk string) string {
	msg := `
The DiskWWN() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the Disk.WWN attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.diskWWN(disk)
}

func (ctx *context) diskWWN(disk string) string {
	info, err := ctx.udevInfo(disk)
	if err != nil {
		return UNKNOWN
	}

	// Trying ID_WWN_WITH_EXTENSION and falling back to ID_WWN is the same logic lsblk uses
	if wwn, ok := info["ID_WWN_WITH_EXTENSION"]; ok {
		return wwn
	}
	if wwn, ok := info["ID_WWN"]; ok {
		return wwn
	}
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
	ctx := contextFromEnv()
	return ctx.diskPartitions(disk)
}

func (ctx *context) diskPartitions(disk string) []*Partition {
	out := make([]*Partition, 0)
	path := filepath.Join(ctx.pathSysBlock(), disk)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, file := range files {
		fname := file.Name()
		if !strings.HasPrefix(fname, disk) {
			continue
		}
		size := ctx.partitionSizeBytes(fname)
		mp, pt, ro := ctx.partitionInfo(fname)
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

// Disks has been deprecated in 0.2. Please use the BlockInfo.Disks attribute.
// TODO(jaypipes): Remove in 1.0.
func Disks() []*Disk {
	msg := `
The Disks() function has been DEPRECATED and will be removed in the
1.0 release of ghw. Please use the BlockInfo.Disks attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.disks()
}

func (ctx *context) disks() []*Disk {
	// In Linux, we could use the fdisk, lshw or blockdev commands to list disk
	// information, however all of these utilities require root privileges to
	// run. We can get all of this information by examining the /sys/block
	// and /sys/class/block files
	disks := make([]*Disk, 0)
	files, err := ioutil.ReadDir(ctx.pathSysBlock())
	if err != nil {
		return nil
	}
	for _, file := range files {
		dname := file.Name()

		busType := BUS_TYPE_UNKNOWN
		if strings.HasPrefix(dname, "sd") {
			busType = BUS_TYPE_SCSI
		} else if strings.HasPrefix(dname, "hd") {
			busType = BUS_TYPE_IDE
		} else if strings.HasPrefix(dname, "vd") {
			busType = BUS_TYPE_VIRTIO
		} else if regexNVMeDev.MatchString(dname) {
			busType = BUS_TYPE_NVME
		}
		if busType == BUS_TYPE_UNKNOWN {
			continue
		}

		size := ctx.diskSizeBytes(dname)
		pbs := ctx.diskPhysicalBlockSizeBytes(dname)
		busPath := ctx.diskBusPath(dname)
		node := ctx.diskNUMANodeID(dname)
		vendor := ctx.diskVendor(dname)
		model := ctx.diskModel(dname)
		serialNo := ctx.diskSerialNumber(dname)
		wwn := ctx.diskWWN(dname)

		d := &Disk{
			Name:                   dname,
			SizeBytes:              size,
			PhysicalBlockSizeBytes: pbs,
			BusType:                busType,
			BusPath:                busPath,
			NUMANodeID:             node,
			Vendor:                 vendor,
			Model:                  model,
			SerialNumber:           serialNo,
			WWN:                    wwn,
		}

		parts := ctx.diskPartitions(dname)
		// Map this Disk object into the Partition...
		for _, part := range parts {
			part.Disk = d
		}
		d.Partitions = parts

		disks = append(disks, d)
	}

	return disks
}

// PartitionSizeBytes has been deprecated in 0.2. Please use the
// Partition.SizeBytes attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionSizeBytes(part string) uint64 {
	msg := `
The PartitionSizeBytes() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.SizeBytes attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionSizeBytes(part)
}

func (ctx *context) partitionSizeBytes(part string) uint64 {
	// Allow calling PartitionSize with either the full partition name
	// "/dev/sda1" or just "sda1"
	part = strings.TrimPrefix(part, "/dev")
	disk := part[0:3]
	if m := regexNVMePart.FindStringSubmatch(part); len(m) > 0 {
		disk = m[1]
	}
	path := filepath.Join(ctx.pathSysBlock(), disk, part, "size")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0
	}
	return uint64(i) * sectorSize
}

// PartitionInfo has been deprecated in 0.2. Please use the Partition struct.
// TODO(jaypipes): Remove in 1.0.
func PartitionInfo(part string) (string, string, bool) {
	msg := `
The PartitionInfo() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition struct.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionInfo(part)
}

// Given a full or short partition name, returns the mount point, the type of
// the partition and whether it's readonly
func (ctx *context) partitionInfo(part string) (string, string, bool) {
	// Allow calling PartitionInfo with either the full partition name
	// "/dev/sda1" or just "sda1"
	if !strings.HasPrefix(part, "/dev") {
		part = "/dev/" + part
	}

	// /etc/mtab entries for mounted partitions look like this:
	// /dev/sda6 / ext4 rw,relatime,errors=remount-ro,data=ordered 0 0
	var r io.ReadCloser
	r, err := os.Open(ctx.pathEtcMtab())
	if err != nil {
		return "", "", true
	}
	defer safeClose(r)

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

// PartitionMountPoint has been deprecated in 0.2. Please use the
// Partition.MountPoint attribute.  TODO(jaypipes): Remove in 1.0.
func PartitionMountPoint(part string) string {
	msg := `
The PartitionMountPoint() function has been DEPRECATED and will be removed in
the 1.0 release of ghw. Please use the Partition.MountPoint attribute.
`
	warn(msg)
	ctx := contextFromEnv()
	return ctx.partitionMountPoint(part)
}

func (ctx *context) partitionMountPoint(part string) string {
	mp, _, _ := ctx.partitionInfo(part)
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
	ctx := contextFromEnv()
	return ctx.partitionType(part)
}

func (ctx *context) partitionType(part string) string {
	_, pt, _ := ctx.partitionInfo(part)
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
	ctx := contextFromEnv()
	return ctx.partitionIsReadOnly(part)
}

func (ctx *context) partitionIsReadOnly(part string) bool {
	_, _, ro := ctx.partitionInfo(part)
	return ro
}
