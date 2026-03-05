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
