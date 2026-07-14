//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package tpm

import (
	"context"
	"errors"
	"runtime"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("tpm load not implemented on " + runtime.GOOS)
}
