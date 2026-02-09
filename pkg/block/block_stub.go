//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package block

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load(_ context.Context) error {
	return errors.New("blockFillInfo not implemented on " + runtime.GOOS)
}
