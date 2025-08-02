// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package chassis

import (
	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/util"
)

func (i *Info) load(opts *option.Options) error {
	i.AssetTag = linuxdmi.Item(opts, "chassis_asset_tag")
	i.SerialNumber = linuxdmi.Item(opts, "chassis_serial")
	i.Type = linuxdmi.Item(opts, "chassis_type")
	typeDesc, found := chassisTypeDescriptions[i.Type]
	if !found {
		typeDesc = util.UNKNOWN
	}
	i.TypeDescription = typeDesc
	i.Vendor = linuxdmi.Item(opts, "chassis_vendor")
	i.Version = linuxdmi.Item(opts, "chassis_version")

	return nil
}
