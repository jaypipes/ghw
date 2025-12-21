// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard

import (
	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/option"
)

func (i *Info) load(opts *option.Options) error {
	i.AssetTag = linuxdmi.Item(opts, "board_asset_tag")
	i.SerialNumber = linuxdmi.Item(opts, "board_serial")
	i.Vendor = linuxdmi.Item(opts, "board_vendor")
	i.Version = linuxdmi.Item(opts, "board_version")
	i.Product = linuxdmi.Item(opts, "board_name")

	return nil
}
