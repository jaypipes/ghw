//go:build !linux && !windows && !darwin
// +build !linux,!windows,!darwin

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cpu

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("cpu.Info.load not implemented on " + runtime.GOOS)
}
