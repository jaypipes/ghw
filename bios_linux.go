// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
)

func biosFillInfo(ctx *context.Context, info *BIOSInfo) error {
	info.Vendor = dmiItem(ctx, "bios_vendor")
	info.Version = dmiItem(ctx, "bios_version")
	info.Date = dmiItem(ctx, "bios_date")

	return nil
}
