//go:build !linux && !windows
// +build !linux,!windows

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package net

import (
	"context"
	"errors"
	"runtime"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("netFillInfo not implemented on " + runtime.GOOS)
}
