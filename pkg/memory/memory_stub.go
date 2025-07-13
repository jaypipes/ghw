//go:build !linux && !windows
// +build !linux,!windows

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"runtime"

	"github.com/jaypipes/ghw/pkg/option"

	"github.com/pkg/errors"
)

func (i *Info) load(opts *option.Options) error {
	return errors.New("mem.Info.load not implemented on " + runtime.GOOS)
}
