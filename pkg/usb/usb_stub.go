//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"context"
	"runtime"

	"errors"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("usb load not implemented on " + runtime.GOOS)
}
