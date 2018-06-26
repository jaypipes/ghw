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

* `ghw.Processor.Id` is the physical processor ID according to the system
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

* `ghw.ProcessorCore.Id` is the identifier that the host gave this core. Note
  that this does *not* necessarily equate to a zero-based index of the core
  within a physical package. For example, the core IDs for an Intel Core i7
  are 0, 1, 2, 8, 9, and 10
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
* `ghw.TopologyInfo.Nodes` is an array of pointers to `ghw.Node` structs, one
  for each topology node (typically physical processor package) found by the
  system

Each `ghw.Node` struct contains the following fields:

* `ghw.Node.Id` is the system's identifier for the node
* `ghw.Node.Cores` is an array of pointers to `ghw.ProcessorCore` structs that
  are contained in this node
* `ghw.Node.Caches` is an array of pointers to `ghw.MemoryCache` structs that
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
