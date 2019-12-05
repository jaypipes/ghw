// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func (ctx *context) dmiFillInfo(info *DMIInfo) error {
	return errors.New("dmiFillInfo not implemented on " + runtime.GOOS)
}
