//go:build linux
// +build linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package main

import (
	"github.com/jaypipes/ghw/cmd/ghw-snapshot/command"
)

func main() {
	command.Execute()
}
