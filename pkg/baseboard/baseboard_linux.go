// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard

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

	i.AssetTag = linuxdmi.Item(ctx, "board_asset_tag")
	i.SerialNumber = linuxdmi.Item(ctx, "board_serial")
	i.Vendor = linuxdmi.Item(ctx, "board_vendor")
	i.Version = linuxdmi.Item(ctx, "board_version")
	i.Product = linuxdmi.Item(ctx, "board_name")

	return nil
}

// loadDeviceTree populates baseboard information from the DeviceTree on systems
// without DMI/SMBIOS. The DeviceTree carries the model, vendor and serial
// number, so the asset tag and version stay unknown.
func (i *Info) loadDeviceTree(ctx context.Context) error {
	i.AssetTag = util.UNKNOWN
	i.SerialNumber = linuxdt.SerialNumber(ctx)
	i.Vendor = linuxdt.Vendor(ctx)
	i.Version = util.UNKNOWN
	i.Product = linuxdt.Model(ctx)

	return nil
}
