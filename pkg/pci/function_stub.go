//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

func (i *Info) GetSRIOVDevices() []*Device {
	return nil
}
