// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
)

func chassisFillInfo(ctx *context.Context, info *ChassisInfo) error {
	info.AssetTag = dmiItem(ctx, "chassis_asset_tag")
	info.SerialNumber = dmiItem(ctx, "chassis_serial")
	info.Type = dmiItem(ctx, "chassis_type")
	typeDesc, found := chassisTypeDescriptions[info.Type]
	if !found {
		typeDesc = UNKNOWN
	}
	info.TypeDescription = typeDesc
	info.Vendor = dmiItem(ctx, "chassis_vendor")
	info.Version = dmiItem(ctx, "chassis_version")

	return nil
}
