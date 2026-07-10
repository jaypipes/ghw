// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package linuxdmi

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/internal/log"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

// Available returns true if DMI/SMBIOS data is exposed by the kernel, i.e. the
// host has a populated /sys/class/dmi/id directory. It is used to decide whether
// to fall back to other sources (such as the DeviceTree) for hardware identity.
func Available(ctx context.Context) bool {
	paths := linuxpath.New(ctx)
	fi, err := os.Stat(filepath.Join(paths.SysClassDMI, "id"))
	return err == nil && fi.IsDir()
}

func Item(ctx context.Context, value string) string {
	paths := linuxpath.New(ctx)
	path := filepath.Join(paths.SysClassDMI, "id", value)

	log.Debug(ctx, "reading from %q", path)
	b, err := os.ReadFile(path)
	if err != nil {
		log.Warn(ctx, "Unable to read %s: %s\n", value, err)
		return util.UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
