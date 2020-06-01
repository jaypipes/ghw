// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func productFillInfo(ctx *context.Context, info *ProductInfo) error {

	info.Family = linuxdmi.Item(ctx, "product_family")
	info.Name = linuxdmi.Item(ctx, "product_name")
	info.Vendor = linuxdmi.Item(ctx, "sys_vendor")
	info.SerialNumber = linuxdmi.Item(ctx, "product_serial")
	info.UUID = linuxdmi.Item(ctx, "product_uuid")
	info.SKU = linuxdmi.Item(ctx, "product_sku")
	info.Version = linuxdmi.Item(ctx, "product_version")

	return nil
}
