//go:build !linux && !cgo
// +build !linux,!cgo

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load() error {
	return errors.New("usb load not implemented without cgo on " + runtime.GOOS)
}
