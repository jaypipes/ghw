# `ghw` - Golang HardWare discovery/inspection library [![Build Status](https://travis-ci.org/jaypipes/ghw.svg?branch=master)](https://travis-ci.org/jaypipes/ghw)

`ghw` is a small Golang library providing hardware inspection and discovery.

## Design Principles

* No root privileges needed for discovery

  `ghw` goes the extra mile to be useful without root priveleges. We query for
  host hardware information as directly as possible without relying on shellouts
  to programs like `dmidecode` that require root privileges to execute.

* Well-documented code and plenty of example code

  The code itself should be well-documented, of course, with lots of usage
  examples.

* Interfaces should be consistent across modules

  Each module in the library should be structured in a consistent fashion, and
  the structs returned by various library functions should have consistent
  attribute and method names.

## Usage

You can use the functions in `ghw` to determine various hardware-related
information about the host computer:

* [Memory](#memory)
* [CPU](#cpu)
* [Block storage](#block-storage)
* [Topology](#topology)
* [Network](#network)
* [PCI](#pci)
* [GPU](#gpu)

### Memory

Information about the host computer's memory can be retrieved using the
`ghw.Memory()` function which returns a pointer to a `ghw.MemoryInfo` struct.

The `ghw.MemoryInfo` struct contains three fields:

* `ghw.MemoryInfo.TotalPhysicalBytes` contains the amount of physical memory on
  the host
* `ghw.MemoryInfo.TotalUsableBytes` contains the amount of memory the
  system can actually use. Usable memory accounts for things like the kernel's
  resident memory size and some reserved system bits
* `ghw.MemoryInfo.SupportedPageSizes` is an array of integers representing the
  size, in bytes, of memory pages the system supports

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

	fmt.Println(mem.String())
}
```

Example output from my personal workstation:

```
memory (24GB physical, 24GB usable)
```

### CPU

The `ghw.CPU()` function returns a `ghw.CPUInfo` struct that contains
information about the CPUs on the host system.

`ghw.CPUInfo` contains the following fields:

* `ghw.CPUInfo.TotalCores` has the total number of physical cores the host
  system contains
* `ghw.CPUInfo.TotalCores` has the total number of hardware threads the
  host system contains
* `ghw.CPUInfo.Processors` is an array of `ghw.Processor` structs, one for each
  physical processor package contained in the host

Each `ghw.Processor` struct contains a number of fields:

* `ghw.Processor.Id` is the physical processor `uint32` ID according to the
  system
* `ghw.Processor.NumCores` is the number of physical cores in the processor
  package
* `ghw.Processor.NumThreads` is the number of hardware threads in the processor
  package
* `ghw.Processor.Vendor` is a string containing the vendor name
* `ghw.Processor.Model` is a string containing the vendor's model name
* `ghw.Processor.Capabilities` is an array of strings indicating the features
  the processor has enabled
* `ghw.Processor.Cores` is an array of `ghw.ProcessorCore` structs that are
  packed onto this physical processor

A `ghw.ProcessorCore` has the following fields:

* `ghw.ProcessorCore.Id` is the `uint32` identifier that the host gave this
  core. Note that this does *not* necessarily equate to a zero-based index of
  the core within a physical package. For example, the core IDs for an Intel Core
  i7 are 0, 1, 2, 8, 9, and 10
* `ghw.ProcessorCore.Index` is the zero-based index of the core on the physical
  processor package
* `ghw.ProcessorCore.NumThreads` is the number of hardware threads associated
  with the core
* `ghw.ProcessorCore.LogicalProcessors` is an array of logical processor IDs
  assigned to any processing unit for the core

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

### Block storage

Information about the host computer's local block storage is returned from the
`ghw.Block()` function. This function returns a pointer to a `ghw.BlockInfo`
struct.

The `ghw.BlockInfo` struct contains two fields:

* `ghw.BlockInfo.TotalPhysicalBytes` contains the amount of physical block
  storage on the host
* `ghw.BlockInfo.Disks` is an array of pointers to `ghw.Disk` structs, one for
  each disk drive found by the system

Each `ghw.Disk` struct contains the following fields:

* `ghw.Disk.Name` contains a string with the short name of the disk, e.g. "sda"
* `ghw.Disk.SizeBytes` contains the amount of storage the disk provides
* `ghw.Disk.SectorSizeBytes` contains the size of the sector used on the disk,
  in bytes
* `ghw.Disk.BusType` will be either "scsi" or "ide"
* `ghw.Disk.Vendor` contains a string with the name of the hardware vendor for
  the disk drive
* `ghw.Disk.SerialNumber` contains a string with the disk's serial number
* `ghw.Disk.Partitions` contains an array of pointers to `ghw.Partition`
  structs, one for each partition on the disk

Each `ghw.Partition` struct contains these fields:

* `ghw.Partition.Name` contains a string with the short name of the partition,
  e.g. "sda1"
* `ghw.Partition.SizeBytes` contains the amount of storage the partition
  provides
* `ghw.Partition.MountPoint` contains a string with the partition's mount
  point, or "" if no mount point was discovered
* `ghw.Partition.Type` contains a string indicated the filesystem type for the
  partition, or "" if the system could not determine the type
* `ghw.Partition.IsReadOnly` is a bool indicating the partition is read-only
* `ghw.Partition.Disk` is a pointer to the `ghw.Disk` object associated with
  the partition. This will be `nil` if the `ghw.Partition` struct was returned
  by the `ghw.DiskPartitions()` library function.

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
 /dev/sda (2TB) [SCSI]  LSI - SN #3600508e000000000f8253aac9a1abd0c
  /dev/sda1 (100MB)
  /dev/sda2 (187GB)
  /dev/sda3 (449MB)
  /dev/sda4 (1KB)
  /dev/sda5 (15GB)
  /dev/sda6 (2TB) [ext4] mounted@/
```

### Topology

Information about the host computer's architecture (NUMA vs. SMP), the host's
node layout and processor caches can be retrieved from the `ghw.Topology()`
function. This function returns a pointer to a `ghw.TopologyInfo` struct.

The `ghw.TopologyInfo` struct contains two fields:

* `ghw.TopologyInfo.Architecture` contains an enum with the value `ghw.NUMA` or
  `ghw.SMP` depending on what the topology of the system is
* `ghw.TopologyInfo.Nodes` is an array of pointers to `ghw.TopologyNode`
  structs, one for each topology node (typically physical processor package)
  found by the system

Each `ghw.TopologyNode` struct contains the following fields:

* `ghw.TopologyNode.Id` is the system's `uint32` identifier for the node
* `ghw.TopologyNode.Cores` is an array of pointers to `ghw.ProcessorCore` structs that
  are contained in this node
* `ghw.TopologyNode.Caches` is an array of pointers to `ghw.MemoryCache` structs that
  represent the low-level caches associated with processors and cores on the
  system

See above in the [CPU](#cpu) section for information about the
`ghw.ProcessorCore` struct and how to use and query it.

Each `ghw.MemoryCache` struct contains the following fields:

* `ghw.MemoryCache.Type` is an enum that contains one of `ghw.DATA`,
  `ghw.INSTRUCTION` or `ghw.UNIFIED` depending on whether the cache stores CPU
  instructions, program data, or both
* `ghw.MemoryCache.Level` is a positive integer indicating how close the cache
  is to the processor
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

Information about the host computer's networking hardware is returned from the
`ghw.Network()` function. This function returns a pointer to a
`ghw.NetworkInfo` struct.

The `ghw.NetworkInfo` struct contains one field:

* `ghw.NetworkInfo.NICs` is an array of pointers to `ghw.NIC` structs, one
  for each network interface controller found for the systen

Each `ghw.NIC` struct contains the following fields:

* `ghw.NIC.Name` is the system's identifier for the NIC
* `ghw.NIC.MacAddress` is the MAC address for the NIC, if any
* `ghw.NIC.IsVirtual` is a boolean indicating if the NIC is a virtualized
  device

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
	}
}
```

Example output from my personal workstation:

```
net (2 NICs)
 enp0s25
 wls1
```

### PCI

`ghw` contains a PCI database inspection and querying facility that allows
developers to not only gather information about devices on a local PCI bus but
also query for information about hardware device classes, vendor and product
information.

The `ghw.PCI()` function returns a `ghw.PCIInfo` struct. The `ghw.PCIInfo`
struct contains a number of fields that may be queried for PCI information:

* `ghw.PCIInfo.Classes` is a map, keyed by the PCI class ID (a hex-encoded
  string) of pointers to `ghw.PCIClass` structs, one for each class of PCI
  device known to `ghw`
* `ghw.PCIInfo.Vendors` is a map, keyed by the PCI vendor ID (a hex-encoded
  string) of pointers to `ghw.PCIVendor` structs, one for each PCI vendor
  known to `ghw`
* `ghw.PCIInfo.Products` is a map, keyed by the PCI product ID* (a hex-encoded
  string) of pointers to `ghw.PCIProduct` structs, one for each PCI product
  known to `ghw`

**NOTE**: PCI products are often referred to by their "device ID". We use
the term "product ID" in `ghw` because it more accurately reflects what the
identifier is for: a specific product line produced by the vendor.

#### PCI device classes

Let's take a look at the PCI device class information and how to query the PCI
database for class, subclass, and programming interface information.

Each `ghw.PCIClass` struct contains the following fields:

* `ghw.PCIClass.Id` is the hex-encoded string identifier for the device
  class
* `ghw.PCIClass.Name` is the common name/description of the class
* `ghw.PCIClass.Subclasses` is an array of pointers to
  `ghw.PCISubclass` structs, one for each subclass in the device class

Each `ghw.PCISubclass` struct contains the following fields:

* `ghw.PCISubclass.Id` is the hex-encoded string identifier for the device
  subclass
* `ghw.PCISubclass.Name` is the common name/description of the subclass
* `ghw.PCISubclass.ProgrammingInterfaces` is an array of pointers to
  `ghw.PCIProgrammingInterface` structs, one for each programming interface
   for the device subclass

Each `ghw.PCIProgrammingInterface` struct contains the following fields:

* `ghw.PCIProgrammingInterface.Id` is the hex-encoded string identifier for
  the programming interface
* `ghw.PCIProgrammingInterface.Name` is the common name/description for the
  programming interface

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

	for _, devClass := range pci.Classes {
		fmt.Printf(" Device class: %v ('%v')\n", devClass.Name, devClass.Id)
        for _, devSubclass := range devClass.Subclasses {
            fmt.Printf("    Device subclass: %v ('%v')\n", devSubclass.Name, devSubclass.Id)
            for _, progIface := range devSubclass.ProgrammingInterfaces {
                fmt.Printf("        Programming interface: %v ('%v')\n", progIface.Name, progIface.Id)
            }
        }
	}
}
```

Example output from my personal workstation, snipped for brevity:

```
...
 Device class: Serial bus controller ('0c')
    Device subclass: FireWire (IEEE 1394) ('00')
        Programming interface: Generic ('00')
        Programming interface: OHCI ('10')
    Device subclass: ACCESS Bus ('01')
    Device subclass: SSA ('02')
    Device subclass: USB controller ('03')
        Programming interface: UHCI ('00')
        Programming interface: OHCI ('10')
        Programming interface: EHCI ('20')
        Programming interface: XHCI ('30')
        Programming interface: Unspecified ('80')
        Programming interface: USB Device ('fe')
    Device subclass: Fibre Channel ('04')
    Device subclass: SMBus ('05')
    Device subclass: InfiniBand ('06')
    Device subclass: IPMI SMIC interface ('07')
    Device subclass: SERCOS interface ('08')
    Device subclass: CANBUS ('09')
...
```

#### PCI vendors and products

Let's take a look at the PCI vendor information and how to query the PCI
database for vendor information and the products a vendor supplies.

Each `ghw.PCIVendor` struct contains the following fields:

* `ghw.PCIVendor.Id` is the hex-encoded string identifier for the vendor
* `ghw.PCIVendor.Name` is the common name/description of the vendor
* `ghw.PCIVendor.Products` is an array of pointers to `ghw.PCIProduct`
  structs, one for each product supplied by the vendor

Each `ghw.PCIProduct` struct contains the following fields:

* `ghw.PCIProduct.VendorId` is the hex-encoded string identifier for the
  product's vendor
* `ghw.PCIProduct.Id` is the hex-encoded string identifier for the product
* `ghw.PCIProduct.Name` is the common name/description of the subclass
* `ghw.PCIProduct.Subsystems` is an array of pointers to
  `ghw.PCIProduct` structs, one for each "subsystem" (sometimes called
  "sub-device" in PCI literature) for the product

**NOTE**: A subsystem product may have a different vendor than its "parent" PCI
product. This is sometimes referred to as the "sub-vendor".

Here's some example code that demonstrates listing the PCI vendors with the
most known products:

```go
package main

import (
	"fmt"
	"sort"

	"github.com/jaypipes/ghw"
)

type ByCountProducts []*ghw.PCIVendor

func (v ByCountProducts) Len() int {
	return len(v)
}

func (v ByCountProducts) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByCountProducts) Less(i, j int) bool {
	return len(v[i].Products) > len(v[j].Products)
}

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	vendors := make([]*ghw.PCIVendor, len(pci.Vendors))
	x := 0
	for _, vendor := range pci.Vendors {
		vendors[x] = vendor
		x++
	}

	sort.Sort(ByCountProducts(vendors))

	fmt.Println("Top 5 vendors by product")
	fmt.Println("====================================================")
	for _, vendor := range vendors[0:5] {
		fmt.Printf("%v ('%v') has %d products\n", vendor.Name, vendor.Id, len(vendor.Products))
	}
}
```

which yields (on my local workstation as of July 7th, 2018):

```
Top 5 vendors by product
====================================================
Intel Corporation ('8086') has 3389 products
NVIDIA Corporation ('10de') has 1358 products
Advanced Micro Devices, Inc. [AMD/ATI] ('1002') has 886 products
National Instruments ('1093') has 601 products
Chelsio Communications Inc ('1425') has 525 products
```

The following is an example of querying the PCI product and subsystem
information to find the products which have the most number of subsystems that
have a different vendor than the top-level product. In other words, the two
products which have been re-sold or re-manufactured with the most number of
different companies.

```go
package main

import (
	"fmt"
	"sort"

	"github.com/jaypipes/ghw"
)

type ByCountSeparateSubvendors []*ghw.PCIProduct

func (v ByCountSeparateSubvendors) Len() int {
	return len(v)
}

func (v ByCountSeparateSubvendors) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByCountSeparateSubvendors) Less(i, j int) bool {
	iVendor := v[i].VendorId
	iSetSubvendors := make(map[string]bool, 0)
	iNumDiffSubvendors := 0
	jVendor := v[j].VendorId
	jSetSubvendors := make(map[string]bool, 0)
	jNumDiffSubvendors := 0

	for _, sub := range v[i].Subsystems {
		if sub.VendorId != iVendor {
			iSetSubvendors[sub.VendorId] = true
		}
	}
	iNumDiffSubvendors = len(iSetSubvendors)

	for _, sub := range v[j].Subsystems {
		if sub.VendorId != jVendor {
			jSetSubvendors[sub.VendorId] = true
		}
	}
	jNumDiffSubvendors = len(jSetSubvendors)

	return iNumDiffSubvendors > jNumDiffSubvendors
}

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	products := make([]*ghw.PCIProduct, len(pci.Products))
	x := 0
	for _, product := range pci.Products {
		products[x] = product
		x++
	}

	sort.Sort(ByCountSeparateSubvendors(products))

	fmt.Println("Top 2 products by # different subvendors")
	fmt.Println("====================================================")
	for _, product := range products[0:2] {
		vendorId := product.VendorId
		vendor := pci.Vendors[vendorId]
		setSubvendors := make(map[string]bool, 0)

		for _, sub := range product.Subsystems {
			if sub.VendorId != vendorId {
				setSubvendors[sub.VendorId] = true
			}
		}
		fmt.Printf("%v ('%v') from %v\n", product.Name, product.Id, vendor.Name)
		fmt.Printf(" -> %d subsystems under the following different vendors:\n", len(setSubvendors))
		for subvendorId, _ := range setSubvendors {
			subvendor, exists := pci.Vendors[subvendorId]
			subvendorName := "Unknown subvendor"
			if exists {
				subvendorName = subvendor.Name
			}
			fmt.Printf("      - %v ('%v')\n", subvendorName, subvendorId)
		}
	}
}
```

which yields (on my local workstation as of July 7th, 2018):

```
Top 2 products by # different subvendors
====================================================
RTL-8100/8101L/8139 PCI Fast Ethernet Adapter ('8139') from Realtek Semiconductor Co., Ltd.
 -> 34 subsystems under the following different vendors:
      - OVISLINK Corp. ('149c')
      - EPoX Computer Co., Ltd. ('1695')
      - Red Hat, Inc ('1af4')
      - Mitac ('1071')
      - Netgear ('1385')
      - Micro-Star International Co., Ltd. [MSI] ('1462')
      - Hangzhou Silan Microelectronics Co., Ltd. ('1904')
      - Compex ('11f6')
      - Edimax Computer Co. ('1432')
      - KYE Systems Corporation ('1489')
      - ZyXEL Communications Corporation ('187e')
      - Acer Incorporated [ALI] ('1025')
      - Matsushita Electric Industrial Co., Ltd. ('10f7')
      - Ruby Tech Corp. ('146c')
      - Belkin ('1799')
      - Allied Telesis ('1259')
      - Unex Technology Corp. ('1429')
      - CIS Technology Inc ('1436')
      - D-Link System Inc ('1186')
      - Ambicom Inc ('1395')
      - AOPEN Inc. ('a0a0')
      - TTTech Computertechnik AG (Wrong ID) ('0357')
      - Gigabyte Technology Co., Ltd ('1458')
      - Packard Bell B.V. ('1631')
      - Billionton Systems Inc ('14cb')
      - Kingston Technologies ('2646')
      - Accton Technology Corporation ('1113')
      - Samsung Electronics Co Ltd ('144d')
      - Biostar Microtech Int'l Corp ('1565')
      - U.S. Robotics ('16ec')
      - KTI ('8e2e')
      - Hewlett-Packard Company ('103c')
      - ASUSTeK Computer Inc. ('1043')
      - Surecom Technology ('10bd')
Bt878 Video Capture ('036e') from Brooktree Corporation
 -> 30 subsystems under the following different vendors:
      - iTuner ('aa00')
      - Nebula Electronics Ltd. ('0071')
      - DViCO Corporation ('18ac')
      - iTuner ('aa05')
      - iTuner ('aa0d')
      - LeadTek Research Inc. ('107d')
      - Avermedia Technologies Inc ('1461')
      - Chaintech Computer Co. Ltd ('270f')
      - iTuner ('aa07')
      - iTuner ('aa0a')
      - Microtune, Inc. ('1851')
      - iTuner ('aa01')
      - iTuner ('aa04')
      - iTuner ('aa06')
      - iTuner ('aa0f')
      - iTuner ('aa02')
      - iTuner ('aa0b')
      - Pinnacle Systems, Inc. (Wrong ID) ('bd11')
      - Rockwell International ('127a')
      - Askey Computer Corp. ('144f')
      - Twinhan Technology Co. Ltd ('1822')
      - Anritsu Corp. ('1852')
      - iTuner ('aa08')
      - Hauppauge computer works Inc. ('0070')
      - Pinnacle Systems Inc. ('11bd')
      - Conexant Systems, Inc. ('14f1')
      - iTuner ('aa09')
      - iTuner ('aa03')
      - iTuner ('aa0c')
      - iTuner ('aa0e')
```

#### Listing and accessing host PCI device information

In addition to the above information, the `ghw.PCIInfo` struct has the
following methods:

* `ghw.PCIInfo.ListDevices() []*PCIDevice`
* `ghw.PCIInfo.GetDevice(address string) *PCIDevice`

This methods return either an array of or a single pointer to a `ghw.PCIDevice`
struct, which has the following fields:


* `ghw.PCIDevice.Vendor` is a pointer to a `ghw.PCIVendor` struct that
  describes the device's primary vendor. This will always be non-nil.
* `ghw.PCIDevice.Product` is a pointer to a `ghw.PCIProduct` struct that
  describes the device's primary product. This will always be non-nil.
* `ghw.PCIDevice.Subsystem` is a pointer to a `ghw.PCIProduct` struct that
  describes the device's secondary/sub-product. This will always be non-nil.
* `ghw.PCIDevice.Class` is a pointer to a `ghw.PCIClass` struct that
  describes the device's class. This will always be non-nil.
* `ghw.PCIDevice.Subclass` is a pointer to a `ghw.PCISubclass` struct
  that describes the device's subclass. This will always be non-nil.
* `ghw.PCIDevice.ProgrammingInterface` is a pointer to a
  `ghw.PCIProgrammingInterface` struct that describes the device subclass'
  programming interface. This will always be non-nil.

The following code snippet shows how to call the `ghw.PCIInfo.ListDevices()`
method and output a simple list of PCI address and vendor/product information:

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
	devices := pci.ListDevices()
	if len(devices) == 0 {
		fmt.Printf("error: could not retrieve PCI devices\n")
		return
	}

	for _, device := range devices {
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
	fmt.Printf("Vendor: %s [%s]\n", vendor.Name, vendor.Id)
	product := deviceInfo.Product
	fmt.Printf("Product: %s [%s]\n", product.Name, product.Id)
	subsystem := deviceInfo.Subsystem
	subvendor := pci.Vendors[subsystem.VendorId]
	subvendorName := "UNKNOWN"
	if subvendor != nil {
		subvendorName = subvendor.Name
	}
	fmt.Printf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.Id, subvendorName)
	class := deviceInfo.Class
	fmt.Printf("Class: %s [%s]\n", class.Name, class.Id)
	subclass := deviceInfo.Subclass
	fmt.Printf("Subclass: %s [%s]\n", subclass.Name, subclass.Id)
	progIface := deviceInfo.ProgrammingInterface
	fmt.Printf("Programming Interface: %s [%s]\n", progIface.Name, progIface.Id)
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

Information about the host computer's graphics hardware is returned from the
`ghw.GPU()` function. This function returns a pointer to a `ghw.GPUInfo`
struct.

The `ghw.GPUInfo` struct contains one field:

* `ghw.GPUInfo.GraphicCards` is an array of pointers to `ghw.GraphicsCard`
  structs, one for each graphics card found for the systen

Each `ghw.GraphicsCard` struct contains the following fields:

* `ghw.GraphicsCard.Index` is the system's numeric zero-based index for the
  card on the bus
* `ghw.GraphicsCard.Address` is the PCI address for the graphics card
* `ghw.GraphicsCard.DeviceInfo` is a pointer to a `ghw.PCIDevice` struct
  describing the graphics card. This may be `nil` if no PCI device information
  could be determined for the card.
* `ghw.GraphicsCard.Nodes` is an array of pointers to `ghw.TopologyNode`
  structs, one for each NUMA node that the GPU/graphics card is affined to. On
  non-NUMA systems, this will always be an empty array.

```go
package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	gpu, err := ghw.GPU())
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

## Developers

Contributions to `ghw` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:


```
[jaypipes@uberbox ghw]$ make test
go test github.com/jaypipes/ghw github.com/jaypipes/ghw/ghwc
ok      github.com/jaypipes/ghw 0.084s
?       github.com/jaypipes/ghw/ghwc    [no test files]
```
