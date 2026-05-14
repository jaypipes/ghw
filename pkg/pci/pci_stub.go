//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"context"
	"errors"
	"runtime"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("pciFillInfo not implemented on " + runtime.GOOS)
}

// GetDevice returns a pointer to a Device struct that describes the PCI
// device at the requested address. If no such device could be found, returns
// nil
func (info *Info) GetDevice(address string) *Device {
	return nil
}

// ListDevices returns a list of pointers to Device structs present on the
// host system
func (info *Info) ListDevices() []*Device {
	return nil
}

// DeviceNamesByLinuxSystem returns nil on non-Linux platforms. The
// underlying data source is /sys/bus/pci, which is Linux-only.
func DeviceNamesByLinuxSystem(ctx context.Context, dev *Device) map[string][]string {
	return nil
}

// VPD is a placeholder for the parsed PCI Vital Product Data block
// declared in vpd_linux.go. The sysfs `vpd` file is Linux-only; on
// other platforms this type exists only so callers can compile
// against the same Device struct.
type VPD struct {
	Identifier string            `json:"identifier,omitempty"`
	ReadOnly   map[string]string `json:"read_only,omitempty"`
	ReadWrite  map[string]string `json:"read_write,omitempty"`
}

// ErrVPDUnavailable is returned by Device.VPD on non-Linux platforms
// and by Linux callers when the device has no resolvable sysfs path.
var ErrVPDUnavailable = errors.New("vpd: no sysfs directory associated with device")

// ErrVPDNotPresent is returned by Device.VPD when the sysfs vpd file
// does not exist.
var ErrVPDNotPresent = errors.New("vpd: not present for device")

// VPD returns ErrVPDUnavailable on non-Linux platforms.
func (d *Device) VPD(ctx context.Context) (*VPD, error) {
	return nil, ErrVPDUnavailable
}
