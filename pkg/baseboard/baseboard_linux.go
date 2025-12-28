// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard

import (
	"context"

	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func (i *Info) load(ctx context.Context) error {
	i.AssetTag = linuxdmi.Item(ctx, "board_asset_tag")
	i.SerialNumber = linuxdmi.Item(ctx, "board_serial")
	i.Vendor = linuxdmi.Item(ctx, "board_vendor")
	i.Version = linuxdmi.Item(ctx, "board_version")
	i.Product = linuxdmi.Item(ctx, "board_name")

	return nil
}
