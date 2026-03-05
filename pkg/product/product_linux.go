// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

import (
	"context"

	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func (i *Info) load(ctx context.Context) error {
	i.Family = linuxdmi.Item(ctx, "product_family")
	i.Name = linuxdmi.Item(ctx, "product_name")
	i.Vendor = linuxdmi.Item(ctx, "sys_vendor")
	i.SerialNumber = linuxdmi.Item(ctx, "product_serial")
	i.UUID = linuxdmi.Item(ctx, "product_uuid")
	i.SKU = linuxdmi.Item(ctx, "product_sku")
	i.Version = linuxdmi.Item(ctx, "product_version")

	return nil
}
