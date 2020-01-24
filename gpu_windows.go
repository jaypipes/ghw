// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

type Win32_VideoController struct {
	Caption                 string
	CreationClassName       string
	Description             string
	DeviceID                string
	Name                    string
	PNPDeviceID             string
	SystemCreationClassName string
	SystemName              string
	VideoArchitecture       uint16
	VideoMemoryType         uint16
	VideoModeDescription    string
	VideoProcessor          string
}

func (ctx *context) gpuFillInfo(info *GPUInfo) error {
	// Getting disk drives from WMI
	var win32VideoControllerDescriptions []Win32_VideoController
	q1 := wmi.CreateQuery(&win32VideoControllerDescriptions, "")
	if err := wmi.Query(q1, &win32VideoControllerDescriptions); err != nil {
		return err
	}
	// Converting into standard structures
	cards := make([]*GraphicsCard, 0)
	pci, err := PCI()
	if err != nil {
		return err
	}
	for _, description := range win32VideoControllerDescriptions {
		card := &GraphicsCard{
			Address:    description.DeviceID, // https://stackoverflow.com/questions/32073667/how-do-i-discover-the-pcie-bus-topology-and-slot-numbers-on-the-board
			Index:      0,
			DeviceInfo: pci.GetDevice(description.PNPDeviceID),
		}
		cards = append(cards, card)
	}
	//car.DeviceInfo
	//card.Node
	info.GraphicsCards = cards
	return nil
}
