// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
)

const wqlDiskDrive = "SELECT Caption, CreationClassName, Description, DeviceID, Index, InterfaceType, Manufacturer, MediaType, Model, Name, Partitions, SerialNumber, Size, TotalCylinders, TotalHeads, TotalSectors, TotalTracks, TracksPerCylinder FROM Win32_DiskDrive"

type win32DiskDrive struct {
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

const wqlDiskPartition = "SELECT Access, BlockSize, Caption, CreationClassName, Description, DeviceID, DiskIndex, Index, Name, Size, SystemName, Type FROM Win32_DiskPartition"

type win32DiskPartition struct {
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

const wqlLogicalDiskToPartition = "SELECT Antecedent, Dependent FROM Win32_LogicalDiskToPartition"

type win32LogicalDiskToPartition struct {
	Antecedent string
	Dependent  string
}

const wqlLogicalDisk = "SELECT Caption, CreationClassName, Description, DeviceID, FileSystem, FreeSpace, Name, Size, SystemName FROM Win32_LogicalDisk"

type win32LogicalDisk struct {
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
	win32DiskDriveDescriptions, err := getDiskDrives()
	if err != nil {
		return err
	}

	win32DiskPartitionDescriptions, err := getDiskPartitions()
	if err != nil {
		return err
	}

	win32LogicalDiskToPartitionDescriptions, err := getLogicalDisksToPartitions()
	if err != nil {
		return err
	}

	win32LogicalDiskDescriptions, err := getLogicalDisks()
	if err != nil {
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
			Vendor:                 UNKNOWN, // TODO: add information
			Model:                  diskdrive.Caption,
			SerialNumber:           diskdrive.SerialNumber,
			WWN:                    UNKNOWN, // TODO: add information
			Partitions:             make([]*Partition, 0),
		}
		for _, diskpartition := range win32DiskPartitionDescriptions {
			// Finding disk partition linked to current disk drive
			if diskdrive.Index == diskpartition.DiskIndex {
				disk.PhysicalBlockSizeBytes = diskpartition.BlockSize
				// Finding logical partition linked to current disk partition
				for _, logicaldisk := range win32LogicalDiskDescriptions {
					for _, logicaldisktodiskpartition := range win32LogicalDiskToPartitionDescriptions {
						var desiredAntecedent = "\\\\" + diskpartition.SystemName + "\\root\\cimv2:" + diskpartition.CreationClassName + ".DeviceID=\"" + diskpartition.DeviceID + "\""
						var desiredDependent = "\\\\" + logicaldisk.SystemName + "\\root\\cimv2:" + logicaldisk.CreationClassName + ".DeviceID=\"" + logicaldisk.DeviceID + "\""
						if logicaldisktodiskpartition.Antecedent == desiredAntecedent && logicaldisktodiskpartition.Dependent == desiredDependent {
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

func getDiskDrives() ([]win32DiskDrive, error) {
	// Getting disks drives data from WMI
	var win3232DiskDriveDescriptions []win32DiskDrive
	if err := wmi.Query(wqlDiskDrive, &win3232DiskDriveDescriptions); err != nil {
		return nil, err
	}
	return win3232DiskDriveDescriptions, nil
}

func getDiskPartitions() ([]win32DiskPartition, error) {
	// Getting disk partitions from WMI
	var win32DiskPartitionDescriptions []win32DiskPartition
	if err := wmi.Query(wqlDiskPartition, &win32DiskPartitionDescriptions); err != nil {
		return nil, err
	}
	return win32DiskPartitionDescriptions, nil
}

func getLogicalDisksToPartitions() ([]win32LogicalDiskToPartition, error) {
	// Getting links between logical disks and partitions from WMI
	var win32LogicalDiskToPartitionDescriptions []win32LogicalDiskToPartition
	if err := wmi.Query(wqlLogicalDiskToPartition, &win32LogicalDiskToPartitionDescriptions); err != nil {
		return nil, err
	}
	return win32LogicalDiskToPartitionDescriptions, nil
}

func getLogicalDisks() ([]win32LogicalDisk, error) {
	// Getting logical disks from WMI
	var win32LogicalDiskDescriptions []win32LogicalDisk
	if err := wmi.Query(wqlLogicalDisk, &win32LogicalDiskDescriptions); err != nil {
		return nil, err
	}
	return win32LogicalDiskDescriptions, nil
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
	// See Access property from: https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-diskpartition
	return access == 0x1
}
