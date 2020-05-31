// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func dmiItem(ctx *context.Context, value string) string {
	paths := linuxpath.New(ctx)
	path := filepath.Join(paths.SysClassDMI, "id", value)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		warn("Unable to read %s: %s\n", value, err)
		return UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
