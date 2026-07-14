// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios

import (
	"context"

	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/linuxdt"
	"github.com/jaypipes/ghw/pkg/util"
)

func (i *Info) load(ctx context.Context) error {
	if !linuxdmi.Available(ctx) && linuxdt.Available(ctx) {
		return i.loadDeviceTree(ctx)
	}

	i.Vendor = linuxdmi.Item(ctx, "bios_vendor")
	i.Version = linuxdmi.Item(ctx, "bios_version")
	i.Date = linuxdmi.Item(ctx, "bios_date")

	return nil
}

// loadDeviceTree populates BIOS/firmware information from the DeviceTree on
// systems without DMI/SMBIOS. U-Boot exposes its version under the "chosen"
// node; when present we report it as a "U-Boot" firmware. The DeviceTree has no
// firmware date, so that stays unknown.
func (i *Info) loadDeviceTree(ctx context.Context) error {
	i.Vendor = util.UNKNOWN
	i.Version = util.UNKNOWN
	i.Date = util.UNKNOWN
	if version := linuxdt.UBootVersion(ctx); version != util.UNKNOWN {
		i.Vendor = "U-Boot"
		i.Version = version
	}

	return nil
}
