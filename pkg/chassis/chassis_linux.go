// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package chassis

import (
	"context"

	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/util"
)

func (i *Info) load(ctx context.Context) error {
	i.AssetTag = linuxdmi.Item(ctx, "chassis_asset_tag")
	i.SerialNumber = linuxdmi.Item(ctx, "chassis_serial")
	i.Type = linuxdmi.Item(ctx, "chassis_type")
	typeDesc, found := chassisTypeDescriptions[i.Type]
	if !found {
		typeDesc = util.UNKNOWN
	}
	i.TypeDescription = typeDesc
	i.Vendor = linuxdmi.Item(ctx, "chassis_vendor")
	i.Version = linuxdmi.Item(ctx, "chassis_version")

	return nil
}
