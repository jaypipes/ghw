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

func productFillInfo(ctx *context.Context, info *ProductInfo) error {
	return errors.New("productFillInfo not implemented on " + runtime.GOOS)
}
