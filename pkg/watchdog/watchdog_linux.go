//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package watchdog

import (
	"context"
	"os"

	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func (i *Info) load(ctx context.Context) error {
	paths := linuxpath.New(ctx)

	// A registered watchdog driver exposes a directory per device under
	// /sys/class/watchdog (kernel >= 4.5).
	entries, err := os.ReadDir(paths.SysClassWatchdog)
	if err == nil && len(entries) > 0 {
		i.Present = true
		return nil
	}

	// Fall back to the character device for older kernels without the
	// watchdog sysfs class.
	if _, err := os.Stat(paths.DevWatchdog); err == nil {
		i.Present = true
	}

	return nil
}
