// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios

import (
	"context"

	"github.com/jaypipes/ghw/pkg/linuxdmi"
)

func (i *Info) load(ctx context.Context) error {
	i.Vendor = linuxdmi.Item(ctx, "bios_vendor")
	i.Version = linuxdmi.Item(ctx, "bios_version")
	i.Date = linuxdmi.Item(ctx, "bios_date")

	return nil
}
