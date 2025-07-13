// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package linuxdmi

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/util"
)

func Item(opts *option.Options, value string) string {
	paths := linuxpath.New(opts)
	path := filepath.Join(paths.SysClassDMI, "id", value)

	b, err := os.ReadFile(path)
	if err != nil {
		opts.Warn("Unable to read %s: %s\n", value, err)
		return util.UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
