// +build !linux,!windows
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"

	"github.com/jaypipes/ghw/pkg/context"
)

func cpuFillInfo(ctx *context.Context, info *CPUInfo) error {
	return errors.New("cpuFillInfo not implemented on " + runtime.GOOS)
}
