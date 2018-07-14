//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

type GraphicsCard struct {
	// the PCI address where the graphics card can be found
	Address string
	// The "index" of the card on the bus (generally not useful information,
	// but might as well include it)
	Index int
	// pointer to a PCIDeviceInfo struct that describes the vendor and product
	// model, etc
	DeviceInfo *PCIDeviceInfo
}

type GPUInfo struct {
	GraphicsCards []*GraphicsCard
}

func GPU() (*GPUInfo, error) {
	info := &GPUInfo{}
	err := gpuFillInfo(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (i *GPUInfo) String() string {
	numCardsStr := "cards"
	if len(i.GraphicsCards) == 1 {
		numCardsStr = "card"
	}
	return fmt.Sprintf(
		"gpu (%d graphics %s)",
		len(i.GraphicsCards),
		numCardsStr,
	)
}
