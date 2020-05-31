// +build !linux,!darwin,!windows
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

func blockFillInfo(ctx *context.Context, info *BlockInfo) error {
	return errors.New("blockFillInfo not implemented on " + runtime.GOOS)
}
