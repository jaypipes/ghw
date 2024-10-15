//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/pci"
)

type AcceleratorDevice struct {
	// the PCI address where the accelerator device can be found
	Address string `json:"address"`
	// pointer to a PCIDevice struct that describes the vendor and product
	// model, etc
	PCIDevice *pci.Device `json:"pci_device"`
}

func (dev *AcceleratorDevice) String() string {
	deviceStr := dev.Address
	if dev.PCIDevice != nil {
		deviceStr = dev.PCIDevice.String()
	}
	nodeStr := ""
	return fmt.Sprintf(
		"device %s@%s",
		nodeStr,
		deviceStr,
	)
}

type Info struct {
	ctx     *context.Context
	Devices []*AcceleratorDevice `json:"devices"`
}

// New returns a pointer to an Info struct that contains information about the
// accelerator devices on the host system
func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}

	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *Info) String() string {
	numDevsStr := "devices"
	if len(i.Devices) == 1 {
		numDevsStr = "device"
	}
	return fmt.Sprintf(
		"processing accelerators (%d %s)",
		len(i.Devices),
		numDevsStr,
	)
}

// simple private struct used to encapsulate processing accelerators information in a top-level
// "accelerator" YAML/JSON map/object key
type acceleratorPrinter struct {
	Info *Info `json:"accelerator"`
}

// YAMLString returns a string with the processing accelerators information formatted as YAML
// under a top-level "accelerator:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, acceleratorPrinter{i})
}

// JSONString returns a string with the processing accelerators information formatted as JSON
// under a top-level "accelerator:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, acceleratorPrinter{i}, indent)
}
