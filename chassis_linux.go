// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func chassisFillInfo(ctx *context.Context, info *ChassisInfo) error {
	info.AssetTag = linuxdmi.Item(ctx, "chassis_asset_tag")
	info.SerialNumber = linuxdmi.Item(ctx, "chassis_serial")
	info.Type = linuxdmi.Item(ctx, "chassis_type")
	typeDesc, found := chassisTypeDescriptions[info.Type]
	if !found {
		typeDesc = UNKNOWN
	}
	info.TypeDescription = typeDesc
	info.Vendor = linuxdmi.Item(ctx, "chassis_vendor")
	info.Version = linuxdmi.Item(ctx, "chassis_version")

	return nil
}
