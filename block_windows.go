// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
)

type Win32_DiskDrive struct {
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	Index             uint32 // Used to link with partition
	InterfaceType     string
	Manufacturer      string
	MediaType         string
	Model             string
	Name              string
	Partitions        int32
	SerialNumber      string
	Size              uint64
	TotalCylinders    int64
	TotalHeads        int32
	TotalSectors      int64
	TotalTracks       int64
	TracksPerCylinder int32
}

/*type Win32_DiskDriveToDiskPartition struct {
}*/

type Win32_DiskPartition struct {
	Access            uint16
	BlockSize         uint64
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	DiskIndex         uint32 // Used to link with Disk Drive
	Index             uint32
	Name              string
	Size              int64
	SystemName        string
	Type              string
}

type Win32_LogicalDiskToPartition struct {
	Antecedent string
	Dependent  string
}

type Win32_LogicalDisk struct {
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	FileSystem        string
	FreeSpace         uint64
	Name              string
	Size              uint64
	SystemName        string
}

func (ctx *context) blockFillInfo(info *BlockInfo) error {
	// Getting disk drives from WMI
	var win32DiskDriveDescriptions []Win32_DiskDrive
	q1 := wmi.CreateQuery(&win32DiskDriveDescriptions, "")
	if err := wmi.Query(q1, &win32DiskDriveDescriptions); err != nil {
		return err
	}

	// Getting disk partitions from WMI
	var win32DiskPartitionDescriptions []Win32_DiskPartition
	q2 := wmi.CreateQuery(&win32DiskPartitionDescriptions, "")
	if err := wmi.Query(q2, &win32DiskPartitionDescriptions); err != nil {
		return err
	}

	// Getting links between logical disks and partitions from WMI
	var win32LogicalDiskToPartitionDescriptions []Win32_LogicalDiskToPartition
	q3 := wmi.CreateQuery(&win32LogicalDiskToPartitionDescriptions, "")
	if err := wmi.Query(q3, &win32LogicalDiskToPartitionDescriptions); err != nil {
		return err
	}

	// Getting logical disks from WMI
	var win32LogicalDiskDescriptions []Win32_LogicalDisk
	q4 := wmi.CreateQuery(&win32LogicalDiskDescriptions, "")
	if err := wmi.Query(q4, &win32LogicalDiskDescriptions); err != nil {
		return err
	}

	// Converting into standard structures
	disks := make([]*Disk, 0)
	for _, diskdrive := range win32DiskDriveDescriptions {
		disk := &Disk{
			Name:                   diskdrive.Name,
			SizeBytes:              diskdrive.Size,
			PhysicalBlockSizeBytes: 0,
			DriveType:              toDriveType(diskdrive.MediaType),
			StorageController:      toStorageController(diskdrive.InterfaceType),
			BusType:                toBusType(diskdrive.InterfaceType),
			BusPath:                UNKNOWN, // TODO: add information
			//NUMANodeID:             node,
			Vendor:       UNKNOWN, // TODO: add information
			Model:        diskdrive.Caption,
			SerialNumber: diskdrive.SerialNumber,
			WWN:          UNKNOWN, // TODO: add information
			Partitions:   make([]*Partition, 0),
		}
		for _, diskpartition := range win32DiskPartitionDescriptions {
			// Finding disk partition linked to current disk drive
			if diskdrive.Index == diskpartition.DiskIndex {
				disk.PhysicalBlockSizeBytes = diskpartition.BlockSize
				//fmt.Printf("Disk Partition %#v\n", diskdrive)
				// Finding logical partition linked to current disk partition
				for _, logicaldisk := range win32LogicalDiskDescriptions {
					for _, logicaldisktodiskpartition := range win32LogicalDiskToPartitionDescriptions {
						//fmt.Printf("\nRelation %#v\n", logicaldisktodiskpartition)
						//fmt.Printf("logicaldisk %#v\n", logicaldisk)
						//fmt.Printf("diskpartition %#v\n", diskpartition)
						var desiredAntecedent = "\\\\" + diskpartition.SystemName + "\\root\\cimv2:" + diskpartition.CreationClassName + ".DeviceID=\"" + diskpartition.DeviceID + "\""
						var desiredDependent = "\\\\" + logicaldisk.SystemName + "\\root\\cimv2:" + logicaldisk.CreationClassName + ".DeviceID=\"" + logicaldisk.DeviceID + "\""
						//fmt.Printf("-- Antecedent\n%#v\n%#v\n", desiredAntecedent, logicaldisktodiskpartition.Antecedent)
						//fmt.Printf("-- Dependent\n%#v\n%#v\n", desiredDependent, logicaldisktodiskpartition.Dependent)
						if logicaldisktodiskpartition.Antecedent == desiredAntecedent && logicaldisktodiskpartition.Dependent == desiredDependent {
							//fmt.Printf("--------------- Disk drive %#v\n", diskdrive)
							//fmt.Printf("--------------- Disk Partition %#v\n", diskpartition)
							//fmt.Printf("--------------- Logical Disk %#v\n", logicaldisk)
							// Appending Partition
							p := &Partition{
								Name:        logicaldisk.Caption,
								Label:       logicaldisk.Caption,
								SizeBytes:   logicaldisk.Size,
								UsableBytes: logicaldisk.FreeSpace,
								MountPoint:  logicaldisk.DeviceID,
								Type:        diskpartition.Type,
								IsReadOnly:  toReadOnly(diskpartition.Access), // TODO: add information
							}
							disk.Partitions = append(disk.Partitions, p)
							break
						}
					}
				}
			}
		}
		// Appending Disk
		disks = append(disks, disk)
	}

	info.Disks = disks
	var tpb uint64
	for _, d := range info.Disks {
		tpb += d.SizeBytes
	}
	info.TotalPhysicalBytes = tpb
	return nil
}

// TODO: improve
func toDriveType(mediaType string) DriveType {
	var driveType DriveType
	mediaType = strings.ToLower(mediaType)
	if strings.Contains(mediaType, "fixed") || strings.Contains(mediaType, "ssd") {
		driveType = DRIVE_TYPE_SSD
	} else if strings.ContainsAny(mediaType, "hdd") {
		driveType = DRIVE_TYPE_HDD
	} else {
		driveType = DRIVE_TYPE_UNKNOWN
	}
	return driveType
}

// TODO: improve
func toStorageController(interfaceType string) StorageController {
	var storageController StorageController
	switch interfaceType {
	case "SCSI":
		storageController = STORAGE_CONTROLLER_SCSI
	case "IDE":
		storageController = STORAGE_CONTROLLER_IDE
	default:
		storageController = STORAGE_CONTROLLER_UNKNOWN
	}
	return storageController
}

// TODO: improve
func toBusType(interfaceType string) BusType {
	var busType BusType
	switch interfaceType {
	case "SCSI":
		busType = BUS_TYPE_SCSI
	case "IDE":
		busType = BUS_TYPE_IDE
	default:
		busType = BUS_TYPE_UNKNOWN
	}
	return busType
}

// TODO: improve
func toReadOnly(access uint16) bool {
	var readOnly bool
	switch access {
	case 0x1: // See Access property from: https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-diskpartition
		readOnly = true
	default:
		readOnly = false
	}
	return readOnly
}
