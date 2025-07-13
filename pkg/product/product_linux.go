// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

import (
	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/option"
)

func (i *Info) load(opts *option.Options) error {
	i.Family = linuxdmi.Item(opts, "product_family")
	i.Name = linuxdmi.Item(opts, "product_name")
	i.Vendor = linuxdmi.Item(opts, "sys_vendor")
	i.SerialNumber = linuxdmi.Item(opts, "product_serial")
	i.UUID = linuxdmi.Item(opts, "product_uuid")
	i.SKU = linuxdmi.Item(opts, "product_sku")
	i.Version = linuxdmi.Item(opts, "product_version")

	return nil
}
