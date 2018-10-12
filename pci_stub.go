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

func pciFillInfo(info *PCIInfo) error {
	return errors.New("pciFillInfo not implemented on " + runtime.GOOS)
}

func (info *PCIInfo) GetDevice(address string) *PCIDevice {
	return nil
}

func (info *PCIInfo) ListDevices() []*PCIDevice {
	return nil
}
