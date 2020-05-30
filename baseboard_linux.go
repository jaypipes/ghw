// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
)

func baseboardFillInfo(ctx *context.Context, info *BaseboardInfo) error {
	info.AssetTag = dmiItem(ctx, "board_asset_tag")
	info.SerialNumber = dmiItem(ctx, "board_serial")
	info.Vendor = dmiItem(ctx, "board_vendor")
	info.Version = dmiItem(ctx, "board_version")

	return nil
}
