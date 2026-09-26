// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package block

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"howett.net/plist"

	"github.com/jaypipes/ghw/internal/config"
)

type diskOrPartitionPlistNode struct {
	Content          string
	DeviceIdentifier string
	DiskUUID         string
	VolumeName       string
	VolumeUUID       string
	Size             int64
	MountPoint       string
	Partitions       []diskOrPartitionPlistNode
	APFSVolumes      []diskOrPartitionPlistNode
}

type diskUtilListPlist struct {
	AllDisks              []string
	AllDisksAndPartitions []diskOrPartitionPlistNode
	VolumesFromDisks      []string
	WholeDisks            []string
}

type diskUtilInfoPlist struct {
	AESHardware                                 bool   // true
	Bootable                                    bool   // true
	BooterDeviceIdentifier                      string // disk1s2
	BusProtocol                                 string // PCI-Express
	CanBeMadeBootable                           bool   // false
	CanBeMadeBootableRequiresDestroy            bool   // false
	Content                                     string // some-uuid-foo-bar
	DeviceBlockSize                             int64  // 4096
	DeviceIdentifier                            string // disk1s1
	DeviceNode                                  string // /dev/disk1s1
	DeviceTreePath                              string // IODeviceTree:/PCI0@0/RP17@1B/ANS2@0/AppleANS2Controller
	DiskUUID                                    string // some-uuid-foo-bar
	Ejectable                                   bool   // false
	EjectableMediaAutomaticUnderSoftwareControl bool   // false
	EjectableOnly                               bool   // false
	FilesystemName                              string // APFS
	FilesystemType                              string // apfs
	FilesystemUserVisibleName                   string // APFS
	FreeSpace                                   int64  // 343975677952
	GlobalPermissionsEnabled                    bool   // true
	IOKitSize                                   int64  // 499963174912
	IORegistryEntryName                         string // Macintosh HD
	Internal                                    bool   // true
	MediaName                                   string //
	MediaType                                   string // Generic
	MountPoint                                  string // /
	ParentWholeDisk                             string // disk1
	PartitionMapPartition                       bool   // false
	RAIDMaster                                  bool   // false
	RAIDSlice                                   bool   // false
	RecoveryDeviceIdentifier                    string // disk1s3
	Removable                                   bool   // false
	RemovableMedia                              bool   // false
	RemovableMediaOrExternalDevice              bool   // false
	SMARTStatus                                 string // Verified
	Size                                        int64  // 499963174912
	SolidState                                  bool   // true
	SupportsGlobalPermissionsDisable            bool   // true
	SystemImage                                 bool   // false
	TotalSize                                   int64  // 499963174912
	VolumeAllocationBlockSize                   int64  // 4096
	VolumeName                                  string // Macintosh HD
	VolumeSize                                  int64  // 499963174912
	VolumeUUID                                  string // some-uuid-foo-bar
	WholeDisk                                   bool   // false
	Writable                                    bool   // true
	WritableMedia                               bool   // true
	WritableVolume                              bool   // true
	// also has a SMARTDeviceSpecificKeysMayVaryNotGuaranteed dict with various info
	// NOTE: VolumeUUID sometimes == DiskUUID, but not always. So far Content is always a different UUID.
}

type ioregPlist struct {
	// there's a lot more than just this...
	EntryID      int64  `plist:"IORegistryEntryID"`
	ModelNumber  string `plist:"Model Number"`
	SerialNumber string `plist:"Serial Number"`
	VendorName   string `plist:"Vendor Name"`
}

// ioregTreeNode is used only to walk the IODeviceTree registry plane by
// path, to uniquely identify a device that may share its node name with
// others elsewhere in the tree. It intentionally doesn't carry Model/Serial
// etc: those properties aren't reliably present on this plane's view of an
// entry, only on ioreg's default plane (which is what ioregPlist above is
// populated from).
type ioregTreeNode struct {
	EntryID  int64           `plist:"IORegistryEntryID"`
	Name     string          `plist:"IORegistryEntryName"`
	Location string          `plist:"IORegistryEntryLocation"`
	Children []ioregTreeNode `plist:"IORegistryEntryChildren"`
}

// segment returns this node's IODeviceTree path component, in the same
// "name" or "name@location" form diskutil reports in DeviceTreePath.
func (n ioregTreeNode) segment() string {
	if n.Location == "" {
		return n.Name
	}
	return n.Name + "@" + n.Location
}

func getDiskUtilListPlist() (*diskUtilListPlist, error) {
	out, err := exec.Command("diskutil", "list", "-plist").Output()
	if err != nil {
		return nil, fmt.Errorf("diskutil list failed: %w", err)
	}

	var data diskUtilListPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("diskutil list plist unmarshal failed: %w", err)
	}

	return &data, nil
}

func getDiskUtilInfoPlist(device string) (*diskUtilInfoPlist, error) {
	out, err := exec.Command("diskutil", "info", "-plist", device).Output()
	if err != nil {
		return nil, fmt.Errorf("diskutil info for %q failed: %w", device, err)
	}

	var data diskUtilInfoPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("diskutil info plist unmarshal for %q failed: %w", device, err)
	}

	return &data, nil
}

// getIODeviceTree fetches the entire IODeviceTree registry plane in one
// call, for use with findEntryID to resolve a device tree path to a unique
// IORegistryEntryID.
func getIODeviceTree() (*ioregTreeNode, error) {
	out, err := exec.Command("ioreg", "-a", "-p", "IODeviceTree").Output()
	if err != nil {
		return nil, fmt.Errorf("ioreg IODeviceTree query failed: %w", err)
	}

	var root ioregTreeNode
	if _, err := plist.Unmarshal(out, &root); err != nil {
		return nil, fmt.Errorf("ioreg IODeviceTree unmarshal failed: %w", err)
	}

	return &root, nil
}

// findEntryID walks the given IODeviceTree plane tree following the
// segments of ioDeviceTreePath (e.g.
// "IODeviceTree:/arm-io@10F00000/ans@77400000/iop-ans-nub/AppleANS3NVMeController")
// to find the exact node, and returns its IORegistryEntryID. That ID is a
// property of the underlying kernel object and so is the same no matter
// which registry plane it's looked up through, letting getIoregPlist
// disambiguate devices that share a node name elsewhere in the tree.
//
// Real hardware nodes are conventionally rooted under a node named
// "device-tree" in the IODeviceTree plane (a convention that predates
// Apple Silicon), which is not itself part of DeviceTreePath.
func findEntryID(root *ioregTreeNode, ioDeviceTreePath string) (int64, error) {
	segments := strings.Split(strings.TrimPrefix(ioDeviceTreePath, "IODeviceTree:/"), "/")

	current := root
	found := false
	for i := range current.Children {
		if current.Children[i].Name == "device-tree" {
			current = &current.Children[i]
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("could not find device-tree root node")
	}

	for _, segment := range segments {
		found = false
		for i := range current.Children {
			if current.Children[i].segment() == segment {
				current = &current.Children[i]
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf(
				"could not find node %q while walking device tree path %q",
				segment, ioDeviceTreePath,
			)
		}
	}

	return current.EntryID, nil
}

func getIoregPlist(tree *ioregTreeNode, ioDeviceTreePath string) (*ioregPlist, error) {
	if strings.TrimPrefix(ioDeviceTreePath, "IODeviceTree:/") == "" {
		// Synthesized/virtual disks (e.g. an APFS container with no
		// corresponding physical device) report an empty or root-only
		// DeviceTreePath. There's no ioreg node to look up.
		return nil, nil
	}

	wantID, err := findEntryID(tree, ioDeviceTreePath)
	if err != nil {
		return nil, fmt.Errorf("locate device tree node for %q: %w", ioDeviceTreePath, err)
	}

	name := path.Base(ioDeviceTreePath)

	args := []string{
		"ioreg",
		"-a",      // use XML output
		"-d", "1", // limit device tree output depth to root node
		"-r",       // root device tree at matched node
		"-n", name, // match by name
	}
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return nil, fmt.Errorf("ioreg query for %q failed: %w", ioDeviceTreePath, err)
	}
	if out == nil || len(out) == 0 {
		return nil, nil
	}

	var data []ioregPlist
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("ioreg unmarshal for %q failed: %w", ioDeviceTreePath, err)
	}

	// ioreg -n matches by node name, not by the full device tree path, so
	// more than one node can come back when multiple devices share a name
	// (e.g. multiple identically-named storage controllers, as seen with
	// AppleVirtIOStorageDevice under macOS virtualization). Use the
	// IORegistryEntryID found via findEntryID to pick the right one.
	for _, entry := range data {
		if entry.EntryID == wantID {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf(
		"none of %d ioreg matches for %q had IORegistryEntryID %d",
		len(data), ioDeviceTreePath, wantID,
	)
}

func makePartition(disk, s diskOrPartitionPlistNode, isAPFS bool) (*Partition, error) {
	if s.Size < 0 {
		return nil, fmt.Errorf("invalid size %q of partition %q", s.Size, s.DeviceIdentifier)
	}

	var partType string
	if isAPFS {
		partType = "APFS Volume"
	} else {
		partType = s.Content
	}

	info, err := getDiskUtilInfoPlist(s.DeviceIdentifier)
	if err != nil {
		return nil, err
	}

	return &Partition{
		Disk:       nil, // filled in later
		Name:       s.DeviceIdentifier,
		Label:      s.VolumeName,
		MountPoint: s.MountPoint,
		SizeBytes:  uint64(s.Size),
		Type:       partType,
		IsReadOnly: !info.WritableVolume,
		UUID:       s.VolumeUUID,
	}, nil
}

// driveTypeFromPlist looks at the supplied property list struct and attempts to
// determine the disk type
func driveTypeFromPlist(infoPlist *diskUtilInfoPlist) DriveType {
	dt := DriveTypeHDD
	if infoPlist.SolidState {
		dt = DriveTypeSSD
	}
	// TODO(jaypipes): Figure out how to determine floppy and/or CD/optical
	// drive type on Mac
	return dt
}

// storageControllerFromPlist looks at the supplied property list struct and
// attempts to determine the storage controller in use for the device
func storageControllerFromPlist(infoPlist *diskUtilInfoPlist) StorageController {
	sc := StorageControllerSCSI
	switch {
	case strings.HasSuffix(infoPlist.DeviceTreePath, "IONVMeController"):
		sc = StorageControllerNVMe
	case infoPlist.BusProtocol == "Apple Fabric":
		// Apple Silicon's internal NVMe controller is attached via the SoC's
		// own integrated storage fabric rather than IONVMeController, and its
		// device tree node name (e.g. AppleANS2Controller on T2 Macs,
		// AppleANS3NVMeController on M-series Macs) varies across chip
		// generations, so match on BusProtocol instead, which macOS reports
		// consistently as "Apple Fabric" for it.
		sc = StorageControllerNVMe
	}
	// TODO(jaypipes): I don't know if Mac even supports IDE controllers and
	// the "virtio" controller is libvirt-specific
	return sc
}

func (info *Info) load(ctx context.Context) error {
	if !config.ToolsEnabled(ctx) {
		return fmt.Errorf("DisableTools=true on darwin disables block support entirely.")
	}

	listPlist, err := getDiskUtilListPlist()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	tree, err := getIODeviceTree()
	if err != nil {
		return err
	}

	var tsb uint64
	info.Disks = make([]*Disk, 0, len(listPlist.AllDisksAndPartitions))
	info.Partitions = []*Partition{}

	for _, disk := range listPlist.AllDisksAndPartitions {
		if disk.Size < 0 {
			return fmt.Errorf("invalid size %q of disk %q", disk.Size, disk.DeviceIdentifier)
		}

		infoPlist, err := getDiskUtilInfoPlist(disk.DeviceIdentifier)
		if err != nil {
			return err
		}
		if infoPlist.DeviceBlockSize < 0 {
			return fmt.Errorf("invalid block size %q of disk %q", infoPlist.DeviceBlockSize, disk.DeviceIdentifier)
		}

		busPath := strings.TrimPrefix(infoPlist.DeviceTreePath, "IODeviceTree:")

		ioregPlist, err := getIoregPlist(tree, infoPlist.DeviceTreePath)
		if err != nil {
			return err
		}
		if ioregPlist == nil {
			continue
		}

		// The NUMA node & WWN don't seem to be reported by any tools available by default in macOS.
		diskReport := &Disk{
			Name:                   disk.DeviceIdentifier,
			SizeBytes:              uint64(disk.Size),
			PhysicalBlockSizeBytes: uint64(infoPlist.DeviceBlockSize),
			DriveType:              driveTypeFromPlist(infoPlist),
			IsRemovable:            infoPlist.Removable,
			StorageController:      storageControllerFromPlist(infoPlist),
			BusPath:                busPath,
			NUMANodeID:             -1,
			Vendor:                 ioregPlist.VendorName,
			Model:                  ioregPlist.ModelNumber,
			SerialNumber:           ioregPlist.SerialNumber,
			WWN:                    "",
			WWNNoExtension:         "",
			Partitions:             make([]*Partition, 0, len(disk.Partitions)+len(disk.APFSVolumes)),
		}

		for _, partition := range disk.Partitions {
			part, err := makePartition(disk, partition, false)
			if err != nil {
				return err
			}
			part.Disk = diskReport
			diskReport.Partitions = append(diskReport.Partitions, part)
		}
		for _, volume := range disk.APFSVolumes {
			part, err := makePartition(disk, volume, true)
			if err != nil {
				return err
			}
			part.Disk = diskReport
			diskReport.Partitions = append(diskReport.Partitions, part)
		}

		tsb += uint64(disk.Size)
		info.Disks = append(info.Disks, diskReport)
		info.Partitions = append(info.Partitions, diskReport.Partitions...)
	}
	info.TotalSizeBytes = tsb
	info.TotalPhysicalBytes = tsb

	return nil
}
