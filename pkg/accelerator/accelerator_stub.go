//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

import (
	"runtime"

	"github.com/pkg/errors"

	"github.com/jaypipes/ghw/pkg/option"
)

func (i *Info) load(opt ...option.Option) error {
	return errors.New("accelerator.Info.load not implemented on " + runtime.GOOS)
}
