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

### Memory

Information about the host computer's memory can be retrieved using the
`ghw.memory()` function which returns a pointer to a `ghw.MemoryInfo` struct:

```go
package main

import (
    "fmt"

    "github.com/jaypipes/ghw"
)

func main(args []string) {
    mem := ghw.NewMemoryInfo()

    fmt.Println(mem.String())
}
```

## Developers

Contributions to `ghw` are welcomed! Fork the repo on GitHub and submit a pull
request with your proposed changes. Or, feel free to log an issue for a feature
request or bug report.

### Running tests

You can run unit tests easily using the `make test` command, like so:
