# `ghw` - Go HardWare discovery/inspection library

[![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/ghw.svg)](https://pkg.go.dev/github.com/jaypipes/ghw)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaypipes/ghw)](https://goreportcard.com/report/github.com/jaypipes/ghw)
[![Build Status](https://github.com/jaypipes/ghw/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jaypipes/ghw/actions)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

![ghw mascot](images/ghw-gopher.png)

`ghw` is a Go library providing hardware inspection and discovery for Linux and
Windows. There currently exists partial support for MacOSX.

## Design Principles

* No root privileges needed for discovery

  `ghw` goes the extra mile to be useful without root priveleges. We query for
  host hardware information as directly as possible without relying on shellouts
  to programs like `dmidecode` that require root privileges to execute.

  Elevated privileges are indeed required to query for some information, but
  `ghw` will never error out if blocked from reading that information. Instead,
  `ghw` will print a warning message about the information that could not be
  retrieved. You may disable these warning messages with the
  `GHW_DISABLE_WARNINGS` environment variable.

* Well-documented code and plenty of example code

  The code itself should be well-documented with lots of usage examples.

* Interfaces should be consistent across modules

  Each module in the library should be structured in a consistent fashion, and
  the structs returned by various library functions should have consistent
  attribute and method names.

## Inspecting != Monitoring

`ghw` is a tool for gathering information about your hardware's **capacity**
and **capabilities**.

It is important to point out that `ghw` does **NOT** report information that is
temporary or variable. It is **NOT** a system monitor nor is it an appropriate
tool for gathering data points for metrics that change over time.  If you are
looking for a system that tracks **usage** of CPU, memory, network I/O or disk
I/O, there are plenty of great open source tools that do this! Check out the
[Prometheus project](https://prometheus.io/) for a great example.

## Usage

`ghw` has functions that return an `Info` object about a particular hardware
domain (e.g. CPU, Memory, Block storage, etc).

Use the following functions in `ghw` to inspect information about the host
hardware:

* [`ghw.CPU()`](#cpu)
* [`ghw.Memory()`](#memory)
* [`ghw.Block()`](#block-storage) (block storage)
* [`ghw.Topology()`](#topology) (processor architecture, NUMA topology and
  memory cache hierarchy)
* [`ghw.Network()`](#network)
* [`ghw.PCI()`](#pci)
* [`ghw.GPU()`](#gpu) (graphical processing unit)
* [`ghw.Accelerator()`](#accelerator) (processing accelerators, AI)
* [`ghw.Chassis()`](#chassis)
* [`ghw.BIOS()`](#bios)
* [`ghw.Baseboard()`](#baseboard)
* [`ghw.Product()`](#product)

### CPU

The `ghw.CPU()` function returns a `ghw.CPUInfo` struct that contains
information about the CPUs on the host system.

`ghw.CPUInfo` contains the following fields:

* `ghw.CPUInfo.TotalCores` has the total number of physical cores the host
  system contains
* `ghw.CPUInfo.TotalHardwareThreads` has the total number of hardware threads
  the host system contains
* `ghw.CPUInfo.Processors` is an array of `ghw.Processor` structs, one for each
  physical processor package contained in the host

Each `ghw.Processor` struct contains a number of fields:

* `ghw.Processor.ID` is the physical processor `uint32` ID according to the
  system
* `ghw.Processor.TotalCores` is the number of physical cores in the processor
  package
* `ghw.Processor.TotalHardwareThreads` is the number of hardware threads in the
  processor package
* `ghw.Processor.Vendor` is a string containing the vendor name
* `ghw.Processor.Model` is a string containing the vendor's model name
* `ghw.Processor.Capabilities` (Linux only) is an array of strings indicating
  the features the processor has enabled
* `ghw.Processor.Cores` (Linux only) is an array of `ghw.ProcessorCore` structs
  that are packed onto this physical processor

A `ghw.ProcessorCore` has the following fields:

* `ghw.ProcessorCore.ID` is the `uint32` identifier that the host gave this
  core. Note that this does *not* necessarily equate to a zero-based index of
  the core within a physical package. For example, the core IDs for an Intel Core
  i7 are 0, 1, 2, 8, 9, and 10
* `ghw.ProcessorCore.TotalHardwareThreads` is the number of hardware threads
  associated with the core
* `ghw.ProcessorCore.LogicalProcessors` is an array of ints representing the
  logical processor IDs assigned to any processing unit for the core. These are
  sometimes called the "thread siblings". Logical processor IDs are the
  *zero-based* index of the processor on the host and are *not* related to the
  core ID.

```go
package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/jaypipes/ghw"
)

func main() {
	cpu, err := ghw.CPU()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v", err)
	}

	fmt.Printf("%v\n", cpu)

	for _, proc := range cpu.Processors {
		fmt.Printf(" %v\n", proc)
		for _, core := range proc.Cores {
			fmt.Printf("  %v\n", core)
		}
		if len(proc.Capabilities) > 0 {
			// pretty-print the (large) block of capability strings into rows
			// of 6 capability strings
			rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
			for row := 1; row < rows; row = row + 1 {
				rowStart := (row * 6) - 1
				rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
				rowElems := proc.Capabilities[rowStart:rowEnd]
				capStr := strings.Join(rowElems, " ")
				if row == 1 {
					fmt.Printf("  capabilities: [%s\n", capStr)
				} else if rowEnd < len(proc.Capabilities) {
					fmt.Printf("                 %s\n", capStr)
				} else {
					fmt.Printf("                 %s]\n", capStr)
				}
			}
		}
	}
}
```

Example output from my personal workstation:

```
cpu (1 physical package, 6 cores, 12 hardware threads)
 physical package #0 (6 cores, 12 hardware threads)
  processor core #0 (2 threads), logical processors [0 6]
  processor core #1 (2 threads), logical processors [1 7]
  processor core #2 (2 threads), logical processors [2 8]
  processor core #3 (2 threads), logical processors [3 9]
  processor core #4 (2 threads), logical processors [4 10]
  processor core #5 (2 threads), logical processors [5 11]
  capabilities: [msr pae mce cx8 apic sep
                 mtrr pge mca cmov pat pse36
                 clflush dts acpi mmx fxsr sse
                 sse2 ss ht tm pbe syscall
                 nx pdpe1gb rdtscp lm constant_tsc arch_perfmon
                 pebs bts rep_good nopl xtopology nonstop_tsc
                 cpuid aperfmperf pni pclmulqdq dtes64 monitor
                 ds_cpl vmx est tm2 ssse3 cx16
                 xtpr pdcm pcid sse4_1 sse4_2 popcnt
                 aes lahf_lm pti retpoline tpr_shadow vnmi
                 flexpriority ept vpid dtherm ida arat]
```

### Memory

The `ghw.Memory()` function returns a `ghw.MemoryInfo` struct that contains
information about the RAM on the host system.

`ghw.MemoryInfo` contains the following fields:

* `ghw.MemoryInfo.TotalPhysicalBytes` contains the amount of physical memory on
  the host
* `ghw.MemoryInfo.TotalUsableBytes` contains the amount of memory the
  system can actually use. Usable memory accounts for things like the kernel's
  resident memory size and some reserved system bits. Please note this value is
  **NOT** the amount of memory currently in use by processes in the system. See
  [the discussion][#physical-versus-usage-memory] about the difference.
* `ghw.MemoryInfo.SupportedPageSizes` is an array of integers representing the
  size, in bytes, of memory pages the system supports
* `ghw.MemoryInfo.Modules` is an array of pointers to `ghw.MemoryModule`
  structs, one for each physical [DIMM](https://en.wikipedia.org/wiki/DIMM).
  Currently, this information is only included on Windows, with Linux support
  [planned](https://github.com/jaypipes/ghw/pull/171#issuecomment-597082409).

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	fmt.Println(memory.String())
}
```

Example output from my personal workstation:

```
memory (24GB physical, 24GB usable)
```

#### Physical versus Usable Memory

There has been [some](https://github.com/jaypipes/ghw/pull/171)
[confusion](https://github.com/jaypipes/ghw/issues/183) regarding the
difference between the total physical bytes versus total usable bytes of
memory.

Some of this confusion has been due to a misunderstanding of the term "usable".
As mentioned [above](#inspection!=monitoring), `ghw` does inspection of the
system's capacity.

A host computer has two capacities when it comes to RAM. The first capacity is
the amount of RAM that is contained in all memory banks (DIMMs) that are
attached to the motherboard. `ghw.MemoryInfo.TotalPhysicalBytes` refers to this
first capacity.

There is a (usually small) amount of RAM that is consumed by the bootloader
before the operating system is started (booted). Once the bootloader has booted
the operating system, the amount of RAM that may be used by the operating
system and its applications is fixed. `ghw.MemoryInfo.TotalUsableBytes` refers
to this second capacity.

You can determine the amount of RAM that the bootloader used (that is not made
available to the operating system) by subtracting
`ghw.MemoryInfo.TotalUsableBytes` from `ghw.MemoryInfo.TotalPhysicalBytes`:

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

        phys := memory.TotalPhysicalBytes
        usable := memory.TotalUsableBytes

	fmt.Printf("The bootloader consumes %d bytes of RAM\n", phys - usable)
}
```

Example output from my personal workstation booted into a Windows10 operating
system with a Linux GRUB bootloader:

```
The bootloader consumes 3832720 bytes of RAM
```

### Block storage

The `ghw.Block()` function returns a `ghw.BlockInfo` struct that contains
information about the block storage on the host system.

`ghw.BlockInfo` contains the following fields:

* `ghw.BlockInfo.TotalSizeBytes` contains the amount of physical block storage
  on the host.
* `ghw.BlockInfo.Disks` is an array of pointers to `ghw.Disk` structs, one for
  each disk found by the system

Each `ghw.Disk` struct contains the following fields:

* `ghw.Disk.Name` contains a string with the short name of the disk, e.g. "sda"
* `ghw.Disk.SizeBytes` contains the amount of storage the disk provides
* `ghw.Disk.PhysicalBlockSizeBytes` contains the size of the physical blocks
  used on the disk, in bytes. This is typically the minimum amount of data that
  will be written in a single write operation for the disk.
* `ghw.Disk.IsRemovable` contains a boolean indicating if the disk drive is
  removable
* `ghw.Disk.DriveType` is the type of drive. It is of type `ghw.DriveType`
  which has a `ghw.DriveType.String()` method that can be called to return a
  string representation of the bus. This string will be `HDD`, `FDD`, `ODD`,
  or `SSD`, which correspond to a hard disk drive (rotational), floppy drive,
  optical (CD/DVD) drive and solid-state drive.
* `ghw.Disk.StorageController` is the type of storage controller. It is of type
  `ghw.StorageController` which has a `ghw.StorageController.String()` method
  that can be called to return a string representation of the bus. This string
  will be `SCSI`, `IDE`, `virtio`, `MMC`, or `NVMe`
* `ghw.Disk.BusPath` (Linux, Darwin only) is the filepath to the bus used by
  the disk.
* `ghw.Disk.NUMANodeID` (Linux only) is the numeric index of the NUMA node this
  disk is local to, or -1 if the host system is not a NUMA system or is not
  Linux.
* `ghw.Disk.Vendor` contains a string with the name of the hardware vendor for
  the disk
* `ghw.Disk.Model` contains a string with the vendor-assigned disk model name
* `ghw.Disk.SerialNumber` contains a string with the disk's serial number
* `ghw.Disk.WWN` contains a string with the disk's
  [World Wide Name](https://en.wikipedia.org/wiki/World_Wide_Name)
* `ghw.Disk.Partitions` contains an array of pointers to `ghw.Partition`
  structs, one for each partition on the disk

Each `ghw.Partition` struct contains these fields:

* `ghw.Partition.Name` contains a string with the short name of the partition,
  e.g. `sda1`
* `ghw.Partition.Label` contains the label for the partition itself. On Linux
  systems, this is derived from the `ID_PART_ENTRY_NAME` [udev][udev] entry for
  the partition.
* `ghw.Partition.FilesystemLabel` contains the label for the filesystem housed
  on the partition. On Linux systems, this is derived from the `ID_FS_NAME`
  [udev][udev] entry for the partition.
* `ghw.Partition.SizeBytes` contains the amount of storage the partition
  provides
* `ghw.Partition.MountPoint` contains a string with the partition's mount
  point, or `""` if no mount point was discovered
* `ghw.Partition.Type` contains a string indicated the filesystem type for the
  partition, or `""` if the system could not determine the type
* `ghw.Partition.IsReadOnly` is a bool indicating the partition is read-only
* `ghw.Partition.Disk` is a pointer to the `ghw.Disk` object associated with
  the partition.
* `ghw.Partition.UUID` is a string containing the partition UUID on Linux, the
  partition UUID on MacOS and nothing on Windows. On Linux systems, this is
  derived from the `ID_PART_ENTRY_UUID` [udev][udev] entry for the partition.

[udev]: https://en.wikipedia.org/wiki/Udev

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		fmt.Printf(" %v\n", disk)
		for _, part := range disk.Partitions {
			fmt.Printf("  %v\n", part)
		}
	}
}
```

Example output from my personal workstation:

```
block storage (1 disk, 2TB physical storage)
 sda HDD (2TB) SCSI [@pci-0000:04:00.0-scsi-0:1:0:0 (node #0)] vendor=LSI model=Logical_Volume serial=600508e000000000f8253aac9a1abd0c WWN=0x600508e000000000f8253aac9a1abd0c
  /dev/sda1 (100MB)
  /dev/sda2 (187GB)
  /dev/sda3 (449MB)
  /dev/sda4 (1KB)
  /dev/sda5 (15GB)
  /dev/sda6 (2TB) [ext4] mounted@/
```

> **NOTE**: `ghw` looks in the udev runtime database for some information. If
> you are using `ghw` in a container, remember to bind mount `/dev/disk` and
> `/run` into your container, otherwise `ghw` won't be able to query the udev
> DB or sysfs paths for information.

### Topology

> **NOTE**: Topology support is currently Linux-only. Windows support is
> [planned](https://github.com/jaypipes/ghw/issues/166).

The `ghw.Topology()` function returns a `ghw.TopologyInfo` struct that contains
information about the host computer's architecture (NUMA vs. SMP), the host's
NUMA node layout and processor-specific memory caches.

The `ghw.TopologyInfo` struct contains two fields:

* `ghw.TopologyInfo.Architecture` contains an enum with the value `ghw.NUMA` or
  `ghw.SMP` depending on what the topology of the system is
* `ghw.TopologyInfo.Nodes` is an array of pointers to `ghw.TopologyNode`
  structs, one for each topology node (typically physical processor package)
  found by the system

Each `ghw.TopologyNode` struct contains the following fields:

* `ghw.TopologyNode.ID` is the system's `uint32` identifier for the node
* `ghw.TopologyNode.Memory` is a `ghw.MemoryArea` struct describing the memory
  attached to this node.
* `ghw.TopologyNode.Cores` is an array of pointers to `ghw.ProcessorCore` structs that
  are contained in this node
* `ghw.TopologyNode.Caches` is an array of pointers to `ghw.MemoryCache` structs that
  represent the low-level caches associated with processors and cores on the
  system
* `ghw.TopologyNode.Distance` is an array of distances between NUMA nodes as reported
  by the system.

`ghw.MemoryArea` describes a collection of *physical* RAM on the host.

In the simplest and most common case, all system memory fits in a single memory
area. In more complex host systems, like [NUMA systems][numa], many memory
areas may be present in the host system (e.g. one for each NUMA cell).

[numa]: https://en.wikipedia.org/wiki/Non-uniform_memory_access

The `ghw.MemoryArea` struct contains the following fields:

* `ghw.MemoryArea.TotalPhysicalBytes` contains the amount of physical memory
  associated with this memory area.
* `ghw.MemoryArea.TotalUsableBytes` contains the amount of memory of this
  memory area the system can actually use. Usable memory accounts for things
  like the kernel's resident memory size and some reserved system bits. Please
  note this value is **NOT** the amount of memory currently in use by processes
  in the system. See [the discussion][#physical-versus-usage-memory] about
  the difference.

See above in the [CPU](#cpu) section for information about the
`ghw.ProcessorCore` struct and how to use and query it.

Each `ghw.MemoryCache` struct contains the following fields:

* `ghw.MemoryCache.Type` is an enum that contains one of `ghw.DATA`,
  `ghw.INSTRUCTION` or `ghw.UNIFIED` depending on whether the cache stores CPU
  instructions, program data, or both
* `ghw.MemoryCache.Level` is a positive integer indicating how close the cache
  is to the processor. The lower the number, the closer the cache is to the
  processor and the faster the processor can access its contents
* `ghw.MemoryCache.SizeBytes` is an integer containing the number of bytes the
  cache can contain
* `ghw.MemoryCache.LogicalProcessors` is an array of integers representing the
  logical processors that use the cache

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	topology, err := ghw.Topology()
	if err != nil {
		fmt.Printf("Error getting topology info: %v", err)
	}

	fmt.Printf("%v\n", topology)

	for _, node := range topology.Nodes {
		fmt.Printf(" %v\n", node)
		for _, cache := range node.Caches {
			fmt.Printf("  %v\n", cache)
		}
	}
}
```

Example output from my personal workstation:

```
topology SMP (1 nodes)
 node #0 (6 cores)
  L1i cache (32 KB) shared with logical processors: 3,9
  L1i cache (32 KB) shared with logical processors: 2,8
  L1i cache (32 KB) shared with logical processors: 11,5
  L1i cache (32 KB) shared with logical processors: 10,4
  L1i cache (32 KB) shared with logical processors: 0,6
  L1i cache (32 KB) shared with logical processors: 1,7
  L1d cache (32 KB) shared with logical processors: 11,5
  L1d cache (32 KB) shared with logical processors: 10,4
  L1d cache (32 KB) shared with logical processors: 3,9
  L1d cache (32 KB) shared with logical processors: 1,7
  L1d cache (32 KB) shared with logical processors: 0,6
  L1d cache (32 KB) shared with logical processors: 2,8
  L2 cache (256 KB) shared with logical processors: 2,8
  L2 cache (256 KB) shared with logical processors: 3,9
  L2 cache (256 KB) shared with logical processors: 0,6
  L2 cache (256 KB) shared with logical processors: 10,4
  L2 cache (256 KB) shared with logical processors: 1,7
  L2 cache (256 KB) shared with logical processors: 11,5
  L3 cache (12288 KB) shared with logical processors: 0,1,10,11,2,3,4,5,6,7,8,9
```

### Network

The `ghw.Network()` function returns a `ghw.NetworkInfo` struct that contains
information about the host computer's networking hardware.

The `ghw.NetworkInfo` struct contains one field:

* `ghw.NetworkInfo.NICs` is an array of pointers to `ghw.NIC` structs, one
  for each network interface controller found for the systen

Each `ghw.NIC` struct contains the following fields:

* `ghw.NIC.Name` is the system's identifier for the NIC
* `ghw.NIC.MACAddress` is the Media Access Control (MAC) address for the NIC,
  if any
* `ghw.NIC.IsVirtual` is a boolean indicating if the NIC is a virtualized
  device
* `ghw.NIC.Capabilities` (Linux only) is an array of pointers to
  `ghw.NICCapability` structs that can describe the things the NIC supports.
  These capabilities match the returned values from the `ethtool -k <DEVICE>`
  call on Linux as well as the AutoNegotiation and PauseFrameUse capabilities
  from `ethtool`.
* `ghw.NIC.PCIAddress` (Linux only) is the PCI device address of the device
  backing the NIC.  this is not-nil only if the backing device is indeed a PCI
  device; more backing devices (e.g. USB) will be added in future versions.
* `ghw.NIC.Speed` (Linux only) is a string showing the current link speed.  On
  Linux, this field will be present even if `ethtool` is not available.
* `ghw.NIC.Duplex` (Linux only) is a string showing the current link duplex. On
  Linux, this field will be present even if `ethtool` is not available.
* `ghw.NIC.SupportedLinkModes` (Linux only) is a string slice containing a list
  of supported link modes, e.g. "10baseT/Half", "1000baseT/Full".
* `ghw.NIC.SupportedPorts` (Linux only) is a string slice containing the list
  of supported port types, e.g. "MII", "TP", "FIBRE", "Twisted Pair".
* `ghw.NIC.SupportedFECModes` (Linux only) is a string slice containing a list
  of supported Forward Error Correction (FEC) Modes.
* `ghw.NIC.AdvertisedLinkModes` (Linux only) is a string slice containing the
  link modes being advertised during auto negotiation.
* `ghw.NIC.AdvertisedFECModes` (Linux only) is a string slice containing the
  Forward Error Correction (FEC) modes advertised during auto negotiation.

The `ghw.NICCapability` struct contains the following fields:

* `ghw.NICCapability.Name` is the string name of the capability (e.g.
  "tcp-segmentation-offload")
* `ghw.NICCapability.IsEnabled` is a boolean indicating whether the capability
  is currently enabled/active on the NIC
* `ghw.NICCapability.CanEnable` is a boolean indicating whether the capability
  may be enabled

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main() {
    net, err := ghw.Network()
    if err != nil {
        fmt.Printf("Error getting network info: %v", err)
    }

    fmt.Printf("%v\n", net)

    for _, nic := range net.NICs {
        fmt.Printf(" %v\n", nic)

        enabledCaps := make([]int, 0)
        for x, cap := range nic.Capabilities {
            if cap.IsEnabled {
                enabledCaps = append(enabledCaps, x)
            }
        }
        if len(enabledCaps) > 0 {
            fmt.Printf("  enabled capabilities:\n")
            for _, x := range enabledCaps {
                fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
            }
        }
    }
}
```

Example output from my personal laptop:

```
net (3 NICs)
 docker0
  enabled capabilities:
   - tx-checksumming
   - tx-checksum-ip-generic
   - scatter-gather
   - tx-scatter-gather
   - tx-scatter-gather-fraglist
   - tcp-segmentation-offload
   - tx-tcp-segmentation
   - tx-tcp-ecn-segmentation
   - tx-tcp-mangleid-segmentation
   - tx-tcp6-segmentation
   - udp-fragmentation-offload
   - generic-segmentation-offload
   - generic-receive-offload
   - tx-vlan-offload
   - highdma
   - tx-lockless
   - netns-local
   - tx-gso-robust
   - tx-fcoe-segmentation
   - tx-gre-segmentation
   - tx-gre-csum-segmentation
   - tx-ipxip4-segmentation
   - tx-ipxip6-segmentation
   - tx-udp_tnl-segmentation
   - tx-udp_tnl-csum-segmentation
   - tx-gso-partial
   - tx-sctp-segmentation
   - tx-esp-segmentation
   - tx-vlan-stag-hw-insert
 enp58s0f1
  enabled capabilities:
   - rx-checksumming
   - generic-receive-offload
   - rx-vlan-offload
   - tx-vlan-offload
   - highdma
   - auto-negotiation
 wlp59s0
  enabled capabilities:
   - scatter-gather
   - tx-scatter-gather
   - generic-segmentation-offload
   - generic-receive-offload
   - highdma
   - netns-local
```

### PCI

`ghw` contains a PCI database inspection and querying facility that allows
developers to not only gather information about devices on a local PCI bus but
also query for information about hardware device classes, vendor and product
information.

> **NOTE**: Parsing of the PCI-IDS file database is provided by the separate
> [github.com/jaypipes/pcidb library](http://github.com/jaypipes/pcidb). You
> can read that library's README for more information about the various structs
> that are exposed on the `ghw.PCIInfo` struct.

The `ghw.PCI()` function returns a `ghw.PCIInfo` struct that contains
information about the host computer's PCI devices.

The `ghw.PCIInfo` struct contains one field:

* `ghw.PCIInfo.Devices` is a slice of pointers to `ghw.PCIDevice` structs that
  describe the PCI devices on the host system

> **NOTE**: PCI products are often referred to by their "device ID". We use the
> term "product ID" in `ghw` because it more accurately reflects what the
> identifier is for: a specific product line produced by the vendor.

The `ghw.PCIDevice` struct has the following fields:

* `ghw.PCIDevice.Vendor` is a pointer to a `pcidb.Vendor` struct that
  describes the device's primary vendor. This will always be non-nil.
* `ghw.PCIDevice.Product` is a pointer to a `pcidb.Product` struct that
  describes the device's primary product. This will always be non-nil.
* `ghw.PCIDevice.Subsystem` is a pointer to a `pcidb.Product` struct that
  describes the device's secondary/sub-product. This will always be non-nil.
* `ghw.PCIDevice.Class` is a pointer to a `pcidb.Class` struct that
  describes the device's class. This will always be non-nil.
* `ghw.PCIDevice.Subclass` is a pointer to a `pcidb.Subclass` struct
  that describes the device's subclass. This will always be non-nil.
* `ghw.PCIDevice.ProgrammingInterface` is a pointer to a
  `pcidb.ProgrammingInterface` struct that describes the device subclass'
  programming interface. This will always be non-nil.
* `ghw.PCIDevice.Driver` is a string representing the device driver the
  system is using to handle this device. Can be empty string if this
  information is not available. If the information is not available, this does
  not mean the device is not functioning, but rather that `ghw` was not able to
  retrieve driver information.

The `ghw.PCIAddress` (which is an alias for the `ghw.pci.address.Address`
struct) contains the PCI address fields. It has a `ghw.PCIAddress.String()`
method that returns the canonical Domain:Bus:Device.Function ([D]BDF)
representation of this Address.

The `ghw.PCIAddress` struct has the following fields:

* `ghw.PCIAddress.Domain` is a string representing the PCI domain component of
  the address.
* `ghw.PCIAddress.Bus` is a string representing the PCI bus component of
  the address.
* `ghw.PCIAddress.Device` is a string representing the PCI device component of
  the address.
* `ghw.PCIAddress.Function` is a string representing the PCI function component of
  the address.

> **NOTE**: Older versions (pre-`v0.9.0`) erroneously referred to the `Device`
> field as the `Slot` field. As noted by [@pearsonk](https://github.com/pearsonk)
> in [#220](https://github.com/jaypipes/ghw/issues/220), this was a misnomer.

The following code snippet shows how to list the PCI devices on the host system
and output a simple list of PCI address and vendor/product information:

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}
	fmt.Printf("host PCI devices:\n")
	fmt.Println("====================================================")

	for _, device := range pci.Devices {
		vendor := device.Vendor
		vendorName := vendor.Name
		if len(vendor.Name) > 20 {
			vendorName = string([]byte(vendorName)[0:17]) + "..."
		}
		product := device.Product
		productName := product.Name
		if len(product.Name) > 40 {
			productName = string([]byte(productName)[0:37]) + "..."
		}
		fmt.Printf("%-12s\t%-20s\t%-40s\n", device.Address, vendorName, productName)
	}
}
```

on my local workstation the output of the above looks like the following:

```
host PCI devices:
====================================================
0000:00:00.0	Intel Corporation   	5520/5500/X58 I/O Hub to ESI Port
0000:00:01.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:02.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:03.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:07.0	Intel Corporation   	5520/5500/X58 I/O Hub PCI Express Roo...
0000:00:10.0	Intel Corporation   	7500/5520/5500/X58 Physical and Link ...
0000:00:10.1	Intel Corporation   	7500/5520/5500/X58 Routing and Protoc...
0000:00:14.0	Intel Corporation   	7500/5520/5500/X58 I/O Hub System Man...
0000:00:14.1	Intel Corporation   	7500/5520/5500/X58 I/O Hub GPIO and S...
0000:00:14.2	Intel Corporation   	7500/5520/5500/X58 I/O Hub Control St...
0000:00:14.3	Intel Corporation   	7500/5520/5500/X58 I/O Hub Throttle R...
0000:00:19.0	Intel Corporation   	82567LF-2 Gigabit Network Connection
0000:00:1a.0	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.1	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.2	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1a.7	Intel Corporation   	82801JI (ICH10 Family) USB2 EHCI Cont...
0000:00:1b.0	Intel Corporation   	82801JI (ICH10 Family) HD Audio Contr...
0000:00:1c.0	Intel Corporation   	82801JI (ICH10 Family) PCI Express Ro...
0000:00:1c.1	Intel Corporation   	82801JI (ICH10 Family) PCI Express Po...
0000:00:1c.4	Intel Corporation   	82801JI (ICH10 Family) PCI Express Ro...
0000:00:1d.0	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.1	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.2	Intel Corporation   	82801JI (ICH10 Family) USB UHCI Contr...
0000:00:1d.7	Intel Corporation   	82801JI (ICH10 Family) USB2 EHCI Cont...
0000:00:1e.0	Intel Corporation   	82801 PCI Bridge
0000:00:1f.0	Intel Corporation   	82801JIR (ICH10R) LPC Interface Contr...
0000:00:1f.2	Intel Corporation   	82801JI (ICH10 Family) SATA AHCI Cont...
0000:00:1f.3	Intel Corporation   	82801JI (ICH10 Family) SMBus Controller
0000:01:00.0	NEC Corporation     	uPD720200 USB 3.0 Host Controller
0000:02:00.0	Marvell Technolog...	88SE9123 PCIe SATA 6.0 Gb/s controller
0000:02:00.1	Marvell Technolog...	88SE912x IDE Controller
0000:03:00.0	NVIDIA Corporation  	GP107 [GeForce GTX 1050 Ti]
0000:03:00.1	NVIDIA Corporation  	UNKNOWN
0000:04:00.0	LSI Logic / Symbi...	SAS2004 PCI-Express Fusion-MPT SAS-2 ...
0000:06:00.0	Qualcomm Atheros    	AR5418 Wireless Network Adapter [AR50...
0000:08:03.0	LSI Corporation     	FW322/323 [TrueFire] 1394a Controller
0000:3f:00.0	Intel Corporation   	UNKNOWN
0000:3f:00.1	Intel Corporation   	Xeon 5600 Series QuickPath Architectu...
0000:3f:02.0	Intel Corporation   	Xeon 5600 Series QPI Link 0
0000:3f:02.1	Intel Corporation   	Xeon 5600 Series QPI Physical 0
0000:3f:02.2	Intel Corporation   	Xeon 5600 Series Mirror Port Link 0
0000:3f:02.3	Intel Corporation   	Xeon 5600 Series Mirror Port Link 1
0000:3f:03.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:03.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:03.4	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:04.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:05.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.0	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.1	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.2	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
0000:3f:06.3	Intel Corporation   	Xeon 5600 Series Integrated Memory Co...
```

#### Finding a PCI device by PCI address

In addition to the above information, the `ghw.PCIInfo` struct has the
following method:

* `ghw.PCIInfo.GetDevice(address string)`

The following code snippet shows how to call the `ghw.PCIInfo.GetDevice()`
method and use its returned `ghw.PCIDevice` struct pointer:

```go
package main

import (
	"fmt"
	"os"

	"github.com/jaypipes/ghw"
)

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	addr := "0000:00:00.0"
	if len(os.Args) == 2 {
		addr = os.Args[1]
	}
	fmt.Printf("PCI device information for %s\n", addr)
	fmt.Println("====================================================")
	deviceInfo := pci.GetDevice(addr)
	if deviceInfo == nil {
		fmt.Printf("could not retrieve PCI device information for %s\n", addr)
		return
	}

	vendor := deviceInfo.Vendor
	fmt.Printf("Vendor: %s [%s]\n", vendor.Name, vendor.ID)
	product := deviceInfo.Product
	fmt.Printf("Product: %s [%s]\n", product.Name, product.ID)
	subsystem := deviceInfo.Subsystem
	subvendor := pci.Vendors[subsystem.VendorID]
	subvendorName := "UNKNOWN"
	if subvendor != nil {
		subvendorName = subvendor.Name
	}
	fmt.Printf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.ID, subvendorName)
	class := deviceInfo.Class
	fmt.Printf("Class: %s [%s]\n", class.Name, class.ID)
	subclass := deviceInfo.Subclass
	fmt.Printf("Subclass: %s [%s]\n", subclass.Name, subclass.ID)
	progIface := deviceInfo.ProgrammingInterface
	fmt.Printf("Programming Interface: %s [%s]\n", progIface.Name, progIface.ID)
}
```

Here's a sample output from my local workstation:

```
PCI device information for 0000:03:00.0
====================================================
Vendor: NVIDIA Corporation [10de]
Product: GP107 [GeForce GTX 1050 Ti] [1c82]
Subsystem: UNKNOWN [8613] (Subvendor: ASUSTeK Computer Inc.)
Class: Display controller [03]
Subclass: VGA compatible controller [00]
Programming Interface: VGA controller [00]
```

### GPU

The `ghw.GPU()` function returns a `ghw.GPUInfo` struct that contains
information about the host computer's graphics hardware.

The `ghw.GPUInfo` struct contains one field:

* `ghw.GPUInfo.GraphicCards` is an array of pointers to `ghw.GraphicsCard`
  structs, one for each graphics card found for the system

Each `ghw.GraphicsCard` struct contains the following fields:

* `ghw.GraphicsCard.Index` is the system's numeric zero-based index for the
  card on the bus
* `ghw.GraphicsCard.Address` is the PCI address for the graphics card
* `ghw.GraphicsCard.DeviceInfo` is a pointer to a `ghw.PCIDevice` struct
  describing the graphics card. This may be `nil` if no PCI device information
  could be determined for the card.
* `ghw.GraphicsCard.Node` is an pointer to a `ghw.TopologyNode` struct that the
  GPU/graphics card is affined to. On non-NUMA systems, this will always be
  `nil`.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	gpu, err := ghw.GPU()
	if err != nil {
		fmt.Printf("Error getting GPU info: %v", err)
	}

	fmt.Printf("%v\n", gpu)

	for _, card := range gpu.GraphicsCards {
		fmt.Printf(" %v\n", card)
	}
}
```

Example output from my personal workstation:

```
gpu (1 graphics card)
 card #0 @0000:03:00.0 -> class: 'Display controller' vendor: 'NVIDIA Corporation' product: 'GP107 [GeForce GTX 1050 Ti]'
```

**NOTE**: You can [read more](#pci) about the fields of the `ghw.PCIDevice`
struct if you'd like to dig deeper into PCI subsystem and programming interface
information

**NOTE**: You can [read more](#topology) about the fields of the
`ghw.TopologyNode` struct if you'd like to dig deeper into the NUMA/topology
subsystem

### Accelerator

The `ghw.Accelerator()` function returns a `ghw.AcceleratorInfo` struct that contains
information about the host computer's processing accelerator hardware. In this category
we can find used hardware for AI. The hardware detected in this category will be
processing accelerators (PCI class `1200`), 3D controllers (`0302`) and Display
controllers (`0380`).

The `ghw.AcceleratorInfo` struct contains one field:

* `ghw.AcceleratorInfo.Devices` is an array of pointers to `ghw.AcceleratorDevice`
  structs, one for each processing accelerator card found for the system.

Each `ghw.AcceleratorDevice` struct contains the following fields:

* `ghw.AcceleratorDevice.Address` is the PCI address for the processing accelerator card.
* `ghw.AcceleratorDevice.PCIDevice` is a pointer to a `ghw.PCIDevice` struct.
  describing the processing accelerator card. This may be `nil` if no PCI device
  information could be determined for the card.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	accel, err := ghw.Accelerator()
	if err != nil {
		fmt.Printf("Error getting processing accelerator info: %v", err)
	}

	fmt.Printf("%v\n", accel)

	for _, card := range accel.Devices {
		fmt.Printf(" %v\n", device)
	}
}
```

Example output from a testing machine:

```
processing accelerators (1 device)
 device @0000:00:04.0 -> driver: 'fake_pci_driver' class: 'Processing accelerators' vendor: 'Red Hat, Inc.' product: 'QEMU PCI Test Device'
```

**NOTE**: You can [read more](#pci) about the fields of the `ghw.PCIDevice`
struct if you'd like to dig deeper into PCI subsystem and programming interface
information

### Chassis

The `ghw.Chassis()` function returns a `ghw.ChassisInfo` struct that contains
information about the host computer's hardware chassis.

The `ghw.ChassisInfo` struct contains multiple fields:

* `ghw.ChassisInfo.AssetTag` is a string with the chassis asset tag
* `ghw.ChassisInfo.SerialNumber` is a string with the chassis serial number
* `ghw.ChassisInfo.Type` is a string with the chassis type *code*
* `ghw.ChassisInfo.TypeDescription` is a string with a description of the
  chassis type
* `ghw.ChassisInfo.Vendor` is a string with the chassis vendor
* `ghw.ChassisInfo.Version` is a string with the chassis version

> **NOTE**: These fields are often missing for non-server hardware. Don't be
> surprised to see empty string or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	chassis, err := ghw.Chassis()
	if err != nil {
		fmt.Printf("Error getting chassis info: %v", err)
	}

	fmt.Printf("%v\n", chassis)
}
```

Example output from my personal workstation:

```
chassis type=Desktop vendor=System76 version=thelio-r1
```

> **NOTE**: Some of the values such as serial numbers are shown as unknown
> because the Linux kernel by default disallows access to those fields if
> you're not running as root. They will be populated if it runs as root or
> otherwise you may see warnings like the following:

```
WARNING: Unable to read chassis_serial: open /sys/class/dmi/id/chassis_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

### BIOS

The `ghw.BIOS()` function returns a `ghw.BIOSInfo` struct that contains
information about the host computer's basis input/output system (BIOS).

The `ghw.BIOSInfo` struct contains multiple fields:

* `ghw.BIOSInfo.Vendor` is a string with the BIOS vendor
* `ghw.BIOSInfo.Version` is a string with the BIOS version
* `ghw.BIOSInfo.Date` is a string with the date the BIOS was flashed/created

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	bios, err := ghw.BIOS()
	if err != nil {
		fmt.Printf("Error getting BIOS info: %v", err)
	}

	fmt.Printf("%v\n", bios)
}
```

Example output from my personal workstation:

```
bios vendor=System76 version=F2 Z5 date=11/14/2018
```

### Baseboard

The `ghw.Baseboard()` function returns a `ghw.BaseboardInfo` struct that
contains information about the host computer's hardware baseboard.

The `ghw.BaseboardInfo` struct contains multiple fields:

* `ghw.BaseboardInfo.AssetTag` is a string with the baseboard asset tag
* `ghw.BaseboardInfo.SerialNumber` is a string with the baseboard serial number
* `ghw.BaseboardInfo.Vendor` is a string with the baseboard vendor
* `ghw.BaseboardInfo.Product` is a string with the baseboard name on Linux and
  Product on Windows
* `ghw.BaseboardInfo.Version` is a string with the baseboard version

> **NOTE**: These fields are often missing for non-server hardware. Don't be
> surprised to see empty string or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	baseboard, err := ghw.Baseboard()
	if err != nil {
		fmt.Printf("Error getting baseboard info: %v", err)
	}

	fmt.Printf("%v\n", baseboard)
}
```

Example output from my personal workstation:

```
baseboard vendor=System76 version=thelio-r1
```

> **NOTE**: Some of the values such as serial numbers are shown as unknown
> because the Linux kernel by default disallows access to those fields if
> you're not running as root. They will be populated if it runs as root or
> otherwise you may see warnings like the following:

```
WARNING: Unable to read board_serial: open /sys/class/dmi/id/board_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

### Product

The `ghw.Product()` function returns a `ghw.ProductInfo` struct that
contains information about the host computer's hardware product line.

The `ghw.ProductInfo` struct contains multiple fields:

* `ghw.ProductInfo.Family` is a string describing the product family
* `ghw.ProductInfo.Name` is a string with the product name
* `ghw.ProductInfo.SerialNumber` is a string with the product serial number
* `ghw.ProductInfo.UUID` is a string with the product UUID
* `ghw.ProductInfo.SKU` is a string with the product stock unit identifier
  (SKU)
* `ghw.ProductInfo.Vendor` is a string with the product vendor
* `ghw.ProductInfo.Version` is a string with the product version

> **NOTE**: These fields are often missing for non-server hardware. Don't be
> surprised to see empty string, "Default string" or "None" values.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	product, err := ghw.Product()
	if err != nil {
		fmt.Printf("Error getting product info: %v", err)
	}

	fmt.Printf("%v\n", product)
}
```

Example output from my personal workstation:

```
product family=Default string name=Thelio vendor=System76 sku=Default string version=thelio-r1
```

> **NOTE**: Some of the values such as serial numbers are shown as unknown
> because the Linux kernel by default disallows access to those fields if
> you're not running as root.  They will be populated if it runs as root or
> otherwise you may see warnings like the following:

```
WARNING: Unable to read product_serial: open /sys/class/dmi/id/product_serial: permission denied
```

You can ignore them or use the [Disabling warning messages](#disabling-warning-messages)
feature to quiet things down.

## Advanced Usage

### Disabling warning messages

When `ghw` isn't able to retrieve some information, it may print certain
warning messages to `stderr`. To disable these warnings, simply set the
`GHW_DISABLE_WARNINGS` environs variable:

```
$ ghwc memory
WARNING:
Could not determine total physical bytes of memory. This may
be due to the host being a virtual machine or container with no
/var/log/syslog file, or the current user may not have necessary
privileges to read the syslog. We are falling back to setting the
total physical amount of memory to the total usable amount of memory
memory (24GB physical, 24GB usable)
```

```
$ GHW_DISABLE_WARNINGS=1 ghwc memory
memory (24GB physical, 24GB usable)
```

You can disable warning programmatically using the `WithDisableWarnings` option:

```go

import (
	"github.com/jaypipes/ghw"
)

mem, err := ghw.Memory(ghw.WithDisableWarnings())
```

`WithDisableWarnings` is a alias for the `WithNullAlerter` option, which in turn
leverages the more general `Alerter` feature of ghw.

You may supply a `Alerter` to ghw to redirect all the warnings there, like
logger objects (see for example golang's stdlib `log.Logger`).
`Alerter` is in fact the minimal logging interface `ghw needs.
To learn more, please check the `option.Alerter` interface and the `ghw.WithAlerter()`
function.

### Overriding the root mountpoint `ghw` uses

When `ghw` looks for information about the host system, it considers `/` as its
root mountpoint. So, for example, when looking up CPU information on a Linux
system, `ghw.CPU()` will use the path `/proc/cpuinfo`.

If you are calling `ghw` from a system that has an alternate root mountpoint,
you can either set the `GHW_CHROOT` environment variable to that alternate
path, or call one of the functions like `ghw.CPU()` or `ghw.Memory()` with the
`ghw.WithChroot()` modifier.

For example, if you are executing from within an application container that has
bind-mounted the root host filesystem to the mount point `/host`, you would set
`GHW_CHROOT` to `/host` so that `ghw` can find `/proc/cpuinfo` at
`/host/proc/cpuinfo`.

Alternately, you can use the `ghw.WithChroot()` function like so:

```go
cpu, err := ghw.CPU(ghw.WithChroot("/host"))
```

### Serialization to JSON or YAML

All of the `ghw` `XXXInfo` structs -- e.g. `ghw.CPUInfo` -- have two methods
for producing a serialized JSON or YAML string representation of the contained
information:

* `JSONString()` returns a string containing the information serialized into
  JSON. It accepts a single boolean parameter indicating whether to use
  indentation when outputting the string
* `YAMLString()` returns a string containing the information serialized into
  YAML

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	mem, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	fmt.Printf("%s", mem.YAMLString())
}
```

the above example code prints the following out on my local workstation:

```
memory:
  supported_page_sizes:
  - 1073741824
  - 2097152
  total_physical_bytes: 25263415296
  total_usable_bytes: 25263415296
```

### Overriding a specific mountpoint (Linux only)

When running inside containers, it can be cumbersome to only override the root
mountpoint. Inside containers, when granting access to the host file systems,
it is common to bind-mount them to a non-standard location, like `/sys` on
`/host-sys` or `/proc` to `/host-proc`.  It is rarer to mount them to a common
subtree (e.g. `/sys` to `/host/sys` and `/proc` to `/host/proc`...)

To better cover this use case, `ghw.WithPathOverrides()` can be used to supply
a mapping of directories to mountpoints, like this example shows:

```go
cpu, err := ghw.CPU(ghw.WithPathOverrides(ghw.PathOverrides{
	"/proc": "/host-proc",
	"/sys": "/host-sys",
}))
```

**NOTE**: This feature works in addition and is composable with the
`ghw.WithChroot()` function and `GHW_CHROOT` environment variable.

### Reading hardware information from a `ghw` snapshot (Linux only)

The `ghw-snapshot` tool can create a snapshot of a host's hardware information.

Please read [`SNAPSHOT.md`](SNAPSHOT.md) to learn about creating snapshots with
the `ghw-snapshot` tool.

You can make `ghw` read hardware information from a snapshot created with
`ghw-snapshot` using environment variables or programmatically.

Use the `GHW_SNAPSHOT_PATH` environment variable to specify the filepath to a
snapshot that `ghw` will read to determine hardware information. All the needed
chroot changes will be automatically performed. By default, the snapshot is
unpacked into a temporary directory managed by `ghw`. This temporary directory
is automatically deleted when `ghw` is finished reading the snapshot.

Three other environment variables are relevant if and only if `GHW_SNAPSHOT_PATH`
is not empty:

* `GHW_SNAPSHOT_ROOT` let users specify the directory on which the snapshot
  should be unpacked. This moves the ownership of that directory from `ghw` to
  users. For this reason, `ghw` will *not* automatically clean up the content
  unpacked into `GHW_SNAPSHOT_ROOT`.
* `GHW_SNAPSHOT_EXCLUSIVE` tells `ghw` that the directory is meant only to
  contain the given snapshot, thus `ghw` will *not* attempt to unpack it unless
  the directory is empty.  You can use both `GHW_SNAPSHOT_ROOT` and
  `GHW_SNAPSHOT_EXCLUSIVE` to make sure `ghw` unpacks the snapshot only once
  regardless of how many `ghw` packages (e.g. cpu, memory) access it. Set the
  value of this environment variable to any non-empty string.
* `GHW_SNAPSHOT_PRESERVE` tells `ghw` not to clean up the unpacked snapshot.
  Set the value of this environment variable to any non-empty string.

```go
cpu, err := ghw.CPU(ghw.WithSnapshot(ghw.SnapshotOptions{
	Path: "/path/to/linux-amd64-d4771ed3300339bc75f856be09fc6430.tar.gz",
}))


myRoot := "/my/safe/directory"
cpu, err := ghw.CPU(ghw.WithSnapshot(ghw.SnapshotOptions{
	Path: "/path/to/linux-amd64-d4771ed3300339bc75f856be09fc6430.tar.gz",
	Root: &myRoot,
}))

myOtherRoot := "/my/other/safe/directory"
cpu, err := ghw.CPU(ghw.WithSnapshot(ghw.SnapshotOptions{
	Path:      "/path/to/linux-amd64-d4771ed3300339bc75f856be09fc6430.tar.gz",
	Root:      &myOtherRoot,
	Exclusive: true,
}))
```

### Creating snapshots

You can create `ghw` snapshots using the `ghw-snapshot` tool or
programmatically using the `pkg/snapshot` package.

Below is an example of creating a `ghw` snapshot using the `pkg/snapshot`
package.

```go

import (
	"fmt"
	"os"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

// ...

scratchDir, err := os.MkdirTemp("", "ghw-snapshot-*")
if err != nil {
	fmt.Printf("Error creating clone directory: %v", err)
}
defer os.RemoveAll(scratchDir)

// this step clones all the files and directories ghw cares about
if err := snapshot.CloneTreeInto(scratchDir); err != nil {
	fmt.Printf("error cloning into %q: %v", scratchDir, err)
}

// optionally, you may add extra content into your snapshot.
// ghw will ignore the extra content.
// Glob patterns like `filepath.Glob` are supported.
fileSpecs := []string{
	"/proc/cmdline",
}

// options allows the client code to optionally deference symlinks, or copy
// them into the cloned tree as symlinks
var opts *snapshot.CopyFileOptions
if err := snapshot.CopyFilesInto(fileSpecs, scratchDir, opts); err != nil {
	fmt.Printf("error cloning extra files into %q: %v", scratchDir, err)
}

// automates the creation of the gzipped tarball out of the given tree.
if err := snapshot.PackFrom("my-snapshot.tgz", scratchDir); err != nil {
	fmt.Printf("error packing %q into %q: %v", scratchDir, *output, err)
}
```

## Calling external programs

By default `ghw` may call external programs, for example `ethtool`, to learn
about hardware capabilities.  In some rare circumstances it may be useful to
opt out from this behaviour and rely only on the data provided by
pseudo-filesystems, like sysfs.

The most common use case is when we want to read a snapshot from `ghw`. In
these cases the information provided by tools will be inconsistent with the
data from the snapshot - since they will be run on a different host than the
host the snapshot was created for.

To prevent `ghw` from calling external tools, set the `GHW_DISABLE_TOOLS`
environment variable to any value, or, programmatically, use the
`ghw.WithDisableTools()` function.  The default behaviour of ghw is to call
external tools when available.

> **WARNING**: on all platforms, disabling external tools make ghw return less
> data.  Unless noted otherwise, there is _no fallback flow_ if external tools
> are disabled. On MacOSX/Darwin, disabling external tools disables block
> support entirely

## Developers

[Contributions](CONTRIBUTING.md) to `ghw` are welcomed! Fork the repo on GitHub
and submit a pull request with your proposed changes. Or, feel free to log an
issue for a feature request or bug report.
