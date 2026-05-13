//go:build linux
// +build linux

//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
)

// linuxSystemKinds is the set of Linux sysfs system directories that a
// PCI device may expose under its /sys/bus/pci/devices/<addr> entry to
// surface higher-level kernel-side device names.
var linuxSystemKinds = []string{
	"drm",
	"infiniband",
	"net",
	"nvme",
}

// DeviceNamesByLinuxSystem returns a map, keyed by well-known Linux
// system string ("drm", "infiniband", "net", "nvme"), of lists of
// Linux device names associated with the supplied PCI Device.
//
// For example, a ConnectX NIC may produce:
//
//	{
//	  "infiniband": ["mlx5_0"],
//	  "net":        ["enp4s0"],
//	}
//
// The map only contains keys for which at least one device name was
// found, so an empty map indicates no recognized system links for
// this device. Returns nil if the supplied device has no PCI address
// resolvable to a sysfs path.
func DeviceNamesByLinuxSystem(ctx context.Context, dev *Device) map[string][]string {
	if dev == nil || dev.Address == "" {
		return nil
	}
	pciAddr := pciaddr.FromString(dev.Address)
	if pciAddr == nil {
		return nil
	}
	paths := linuxpath.New(ctx)
	devPath := filepath.Join(paths.SysBusPciDevices, pciAddr.String())
	out := map[string][]string{}
	for _, kind := range linuxSystemKinds {
		entries, err := os.ReadDir(filepath.Join(devPath, kind))
		if err != nil {
			continue
		}
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)
		out[kind] = names
	}
	return out
}
