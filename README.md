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
* Block devices

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
    memory := ghw.Memory()

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
    block := ghw.Block()

    fmt.Println(block.String())
}
```

The `ghw.MemoryInfo` struct contains two fields:

* `ghw.BlockInfo.TotalPhysicalBytes` contains the amount of physical block
  storage on the host
* `ghw.BlockInfo.Disks` is an array of pointers to `ghw.Disk` structs, one for
  each disk drive found by the system

Each `ghw.Disk` struct contains the following fields:

* `ghw.Disk.Name` contains a string with the short name of the disk, e.g. "sda"
* `ghw.Disk.SizeBytes` contains the amount of storage the disk provides
* `ghw.Disk.SectorSize` contains the size of the sector used on the disk,
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

## Developers

Contributions to `ghw` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:
