//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/accelerator"
	"github.com/jaypipes/ghw/pkg/baseboard"
	"github.com/jaypipes/ghw/pkg/bios"
	"github.com/jaypipes/ghw/pkg/block"
	"github.com/jaypipes/ghw/pkg/chassis"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/gpu"
	"github.com/jaypipes/ghw/pkg/memory"
	"github.com/jaypipes/ghw/pkg/net"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/pci"
	pciaddress "github.com/jaypipes/ghw/pkg/pci/address"
	"github.com/jaypipes/ghw/pkg/product"
	"github.com/jaypipes/ghw/pkg/topology"
)

type WithOption = option.Option

var (
	WithChroot      = option.WithChroot
	WithSnapshot    = option.WithSnapshot
	WithAlerter     = option.WithAlerter
	WithNullAlerter = option.WithNullAlerter
	// match the existing environ variable to minimize surprises
	WithDisableWarnings = option.WithNullAlerter
	WithDisableTools    = option.WithDisableTools
	WithPathOverrides   = option.WithPathOverrides
)

type SnapshotOptions = option.SnapshotOptions

type PathOverrides = option.PathOverrides

type CPUInfo = cpu.Info

var (
	CPU = cpu.New
)

type MemoryArea = memory.Area
type MemoryInfo = memory.Info
type MemoryCache = memory.Cache
type MemoryCacheType = memory.CacheType
type MemoryModule = memory.Module

const (
	MemoryCacheTypeUnified = memory.CacheTypeUnified
	// DEPRECATED: Please use MemoryCacheTypeUnified
	MEMORY_CACHE_TYPE_UNIFIED  = memory.CACHE_TYPE_UNIFIED
	MemoryCacheTypeInstruction = memory.CacheTypeInstruction
	// DEPRECATED: Please use MemoryCacheTypeInstruction
	MEMORY_CACHE_TYPE_INSTRUCTION = memory.CACHE_TYPE_INSTRUCTION
	MemoryCacheTypeData           = memory.CacheTypeData
	// DEPRECATED: Please use MemoryCacheTypeData
	MEMORY_CACHE_TYPE_DATA = memory.CACHE_TYPE_DATA
)

var (
	Memory = memory.New
)

type BlockInfo = block.Info
type Disk = block.Disk
type Partition = block.Partition

var (
	Block = block.New
)

type DriveType = block.DriveType

const (
	DriveTypeUnknown = block.DriveTypeUnknown
	// DEPRECATED: Please use DriveTypeUnknown
	DRIVE_TYPE_UNKNOWN = block.DRIVE_TYPE_UNKNOWN
	DriveTypeHDD       = block.DriveTypeHDD
	// DEPRECATED: Please use DriveTypeHDD
	DRIVE_TYPE_HDD = block.DRIVE_TYPE_HDD
	DriveTypeFDD   = block.DriveTypeFDD
	// DEPRECATED: Please use DriveTypeFDD
	DRIVE_TYPE_FDD = block.DRIVE_TYPE_FDD
	DriveTypeODD   = block.DriveTypeODD
	// DEPRECATED: Please use DriveTypeODD
	DRIVE_TYPE_ODD = block.DRIVE_TYPE_ODD
	DriveTypeSSD   = block.DriveTypeSSD
	// DEPRECATED: Please use DriveTypeSSD
	DRIVE_TYPE_SSD = block.DRIVE_TYPE_SSD
)

type StorageController = block.StorageController

const (
	StorageControllerUnknown = block.StorageControllerUnknown
	// DEPRECATED: Please use StorageControllerUnknown
	STORAGE_CONTROLLER_UNKNOWN = block.STORAGE_CONTROLLER_UNKNOWN
	StorageControllerIDE       = block.StorageControllerIDE
	// DEPRECATED: Please use StorageControllerIDE
	STORAGE_CONTROLLER_IDE = block.STORAGE_CONTROLLER_IDE
	StorageControllerSCSI  = block.StorageControllerSCSI
	// DEPRECATED: Please use StorageControllerSCSI
	STORAGE_CONTROLLER_SCSI = block.STORAGE_CONTROLLER_SCSI
	StorageControllerNVMe   = block.StorageControllerNVMe
	// DEPRECATED: Please use StorageControllerNVMe
	STORAGE_CONTROLLER_NVME = block.STORAGE_CONTROLLER_NVME
	StorageControllerVirtIO = block.StorageControllerVirtIO
	// DEPRECATED: Please use StorageControllerVirtIO
	STORAGE_CONTROLLER_VIRTIO = block.STORAGE_CONTROLLER_VIRTIO
	StorageControllerMMC      = block.StorageControllerMMC
	// DEPRECATED: Please use StorageControllerMMC
	STORAGE_CONTROLLER_MMC = block.STORAGE_CONTROLLER_MMC
)

type NetworkInfo = net.Info
type NIC = net.NIC
type NICCapability = net.NICCapability

var (
	Network = net.New
)

type BIOSInfo = bios.Info

var (
	BIOS = bios.New
)

type ChassisInfo = chassis.Info

var (
	Chassis = chassis.New
)

type BaseboardInfo = baseboard.Info

var (
	Baseboard = baseboard.New
)

type TopologyInfo = topology.Info
type TopologyNode = topology.Node

var (
	Topology = topology.New
)

type Architecture = topology.Architecture

const (
	ArchitectureSMP = topology.ArchitectureSMP
	// DEPRECATED: Please use ArchitectureSMP
	ARCHITECTURE_SMP = topology.ArchitectureSMP
	ArchitectureNUMA = topology.ArchitectureNUMA
	// DEPRECATED: Please use ArchitectureNUMA
	ARCHITECTURE_NUMA = topology.ArchitectureNUMA
)

type PCIInfo = pci.Info
type PCIAddress = pciaddress.Address
type PCIDevice = pci.Device

var (
	PCI                  = pci.New
	PCIAddressFromString = pciaddress.FromString
)

type ProductInfo = product.Info

var (
	Product = product.New
)

type GPUInfo = gpu.Info
type GraphicsCard = gpu.GraphicsCard

var (
	GPU = gpu.New
)

type AcceleratorInfo = accelerator.Info
type AcceleratorDevice = accelerator.AcceleratorDevice

var (
	Accelerator = accelerator.New
)
