// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func biosFillInfo(ctx *context.Context, info *BIOSInfo) error {
	info.Vendor = linuxdmi.Item(ctx, "bios_vendor")
	info.Version = linuxdmi.Item(ctx, "bios_version")
	info.Date = linuxdmi.Item(ctx, "bios_date")

	return nil
}
