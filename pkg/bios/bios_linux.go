// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios

import (
	"github.com/jaypipes/ghw/pkg/linuxdmi"
	"github.com/jaypipes/ghw/pkg/option"
)

func (i *Info) load(opts *option.Options) error {
	i.Vendor = linuxdmi.Item(opts, "bios_vendor")
	i.Version = linuxdmi.Item(opts, "bios_version")
	i.Date = linuxdmi.Item(opts, "bios_date")

	return nil
}
