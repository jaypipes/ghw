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

func netFillInfo(info *NetworkInfo) error {
	return errors.New("netFillInfo not implemented on " + runtime.GOOS)
}

func NICs() []*NIC {
	return nil
}
