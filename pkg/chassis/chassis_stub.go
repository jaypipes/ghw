//go:build !linux && !windows
// +build !linux,!windows

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package chassis

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
)

func (i *Info) load(ctx context.Context) error {
	return errors.New("chassisFillInfo not implemented on " + runtime.GOOS)
}
