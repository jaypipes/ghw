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
	// pointer to a PCIDevice struct that describes the vendor and product
	// model, etc
	DeviceInfo *PCIDevice
	// Topology nodes that the graphics card is affined to. Will be nil if the
	// architecture is not NUMA.
	Node *TopologyNode
}

func (card *GraphicsCard) String() string {
	deviceStr := card.Address
	if card.DeviceInfo != nil {
		deviceStr = card.DeviceInfo.String()
	}
	nodeStr := ""
	if card.Node != nil {
		nodeStr = fmt.Sprintf(" [affined to NUMA node %d]", card.Node.Id)
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

func GPU(opts ...*WithOption) (*GPUInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &GPUInfo{}
	if err := ctx.gpuFillInfo(info); err != nil {
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
