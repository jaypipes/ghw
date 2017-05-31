# `gwk` - Golang HardWare discovery/inspection library

`ghw` is a small Golang library providing hardware inspection and discovery.

## Design Principles

### No root privileges needed for discovery

`ghw` goes the extra mile to be useful without root priveleges. We query for
host hardware information as directly as possible without relying on shellouts
to programs like `dmidecode` that require root privileges to execute.

### Well-documented code and plenty of example code

The code itself should be well-documented, of course, with lots of usage
examples.

### Interfaces should be consistent across modules

Each module in the library should be structured in a consistent fashion, and
the structs returned by various library functions should have consistent
attribute and method names.

## Usage

You can use the functions in `ghw` to determine various hardware-related
information about the host computer:

* Memory
* CPU
* Block storage

### Memory

Information about the host computer's memory can be retrieved using the
`ghw.Memory()` function which returns a pointer to a `ghw.MemoryInfo` struct:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main(args []string) {
    memory, err := ghw.Memory()
    if err != nil {
        fmt.Printf("Error getting memory info: %v", err)
    }

    fmt.Println(mem.String())
}
```

The `ghw.MemoryInfo` struct contains three fields:

* `ghw.MemoryInfo.TotalPhysicalBytes` contains the amount of physical memory on
  the host
* `ghw.MemoryInfo.TotalUsableBytes` contains the amount of memory the
  system can actually use. Usable memory accounts for things like the kernel's
  resident memory size and some reserved system bits
* `ghw.SupportedPageSizes` is an array of integers representing the size, in
  bytes, of memory pages the system supports

### CPU

The `ghw.CPU()` function returns a `ghw.CPUInfo` struct that contains
information about the CPUs on the host system:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main(args []string) {
    cpu, err := ghw.CPU()
    if err != nil {
        fmt.Printf("Error getting CPU info: %v", err)
    }

    fmt.Println(cpu.String())

    for _, proc := range cpu.Processors {
        fmt.Println(proc.String())
        for _, core := range p.Cores {
            fmt.Println(core.String())
        }
    }
}
```

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
* `ghw.Processor.Cores is an array of `ghw.ProcessorCore` structs that are
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
* `ghw.LogicalProcessors` is an array of logical processor IDs assigned to any
  processing unit for the core

### Block storage

Information about the host computer's local block storage is returned from the
`ghw.Block()` function. This function returns a pointer to a `ghw.BlockInfo`
struct:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main(args []string) {
    block, err := ghw.Block()
    if err != nil {
        fmt.Printf("Error getting block storage info: %v", err)
    }

    fmt.Println(block.String())

    for _, disk := range block.Disks {
        fmt.Println(disk.String())
        for _, part := range disk.Partitions {
            fmt.Println(part.String())
        }
    }
}
```

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

## Developers

Contributions to `ghw` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:
