//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"strconv"
)

type GraphicsCard struct {
	// the PCI address where the graphics card can be found
	Address string
	// The "index" of the card on the bus (generally not useful information,
	// but might as well include it)
	Index int
	// pointer to a PCIDevice struct that describes the vendor and product
	// model, etc
	DeviceInfo *PCIDevice
	// Array of topology nodes that the graphics card is affined to. Will be empty
	// if the architecture is not NUMA.
	Nodes []*TopologyNode
}

func (card *GraphicsCard) String() string {
	deviceStr := card.Address
	if card.DeviceInfo != nil {
		deviceStr = card.DeviceInfo.String()
	}
	nodeStr := ""
	if len(card.Nodes) > 0 {
		x := 0
		nodeStr = " NUMA nodes ["
		for _, node := range card.Nodes {
			if x > 0 {
				nodeStr += ","
			}
			nodeStr += strconv.Itoa(int(node.Id))
			x++
		}
		nodeStr += "] "
	}
	return fmt.Sprintf(
		"card #%d %s@%s",
		card.Index,
		nodeStr,
		deviceStr,
	)
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
