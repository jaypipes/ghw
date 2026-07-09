// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

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

	i.Family = linuxdmi.Item(ctx, "product_family")
	i.Name = linuxdmi.Item(ctx, "product_name")
	i.Vendor = linuxdmi.Item(ctx, "sys_vendor")
	i.SerialNumber = linuxdmi.Item(ctx, "product_serial")
	i.UUID = linuxdmi.Item(ctx, "product_uuid")
	i.SKU = linuxdmi.Item(ctx, "product_sku")
	i.Version = linuxdmi.Item(ctx, "product_version")

	return nil
}

// loadDeviceTree populates product information from the DeviceTree on systems
// without DMI/SMBIOS. The DeviceTree only carries the model, vendor and serial
// number, so the remaining fields stay unknown.
func (i *Info) loadDeviceTree(ctx context.Context) error {
	i.Family = util.UNKNOWN
	i.Name = linuxdt.Model(ctx)
	i.Vendor = linuxdt.Vendor(ctx)
	i.SerialNumber = linuxdt.SerialNumber(ctx)
	i.UUID = util.UNKNOWN
	i.SKU = util.UNKNOWN
	i.Version = util.UNKNOWN

	return nil
}
