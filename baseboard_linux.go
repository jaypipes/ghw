// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func baseboardFillInfo(ctx *context.Context, info *BaseboardInfo) error {
	info.AssetTag = linuxdmi.Item(ctx, "board_asset_tag")
	info.SerialNumber = linuxdmi.Item(ctx, "board_serial")
	info.Vendor = linuxdmi.Item(ctx, "board_vendor")
	info.Version = linuxdmi.Item(ctx, "board_version")

	return nil
}
