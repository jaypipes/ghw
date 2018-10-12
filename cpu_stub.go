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

func cpuFillInfo(info *CPUInfo) error {
	return errors.New("cpuFillInfo not implemented on " + runtime.GOOS)
}

func Processors() []*Processor {
	return nil
}
