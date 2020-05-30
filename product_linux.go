// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
)

func productFillInfo(ctx *context.Context, info *ProductInfo) error {

	info.Family = dmiItem(ctx, "product_family")
	info.Name = dmiItem(ctx, "product_name")
	info.Vendor = dmiItem(ctx, "sys_vendor")
	info.SerialNumber = dmiItem(ctx, "product_serial")
	info.UUID = dmiItem(ctx, "product_uuid")
	info.SKU = dmiItem(ctx, "product_sku")
	info.Version = dmiItem(ctx, "product_version")

	return nil
}
