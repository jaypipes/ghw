//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sriov

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/pci"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
)

type Device struct {
	Interfaces []string `json:"interfaces"`
	// the PCI address where the SRIOV instance can be found
	Address *pciaddr.Address `json:"address"`
	PCI     *pci.Device      `json:"pci"`
}

func (d Device) ToString(devType string) string {
	deviceStr := d.Address.String()
	nodeStr := ""
	if d.PCI != nil {
		deviceStr = d.PCI.String()
		if d.PCI.Node != nil {
			nodeStr = fmt.Sprintf(" [affined to NUMA node %d]", d.PCI.Node.ID)
		}
	}
	return fmt.Sprintf("%s function %s@%s", devType, nodeStr, deviceStr)
}

type PhysicalFunction struct {
	Device
	MaxVFNum int               `json:"max_vf_num,omitempty"`
	VFs      []VirtualFunction `json:"vfs,omitempty"`
}

type VirtualFunction struct {
	Device
	ID int `json:"id"`
	// Address of the (parent) Physical Function this Virtual Function pertains to.
	ParentAddress *pciaddr.Address `json:"parent_address,omitempty"`
}

func (pf *PhysicalFunction) String() string {
	return fmt.Sprintf("%s with %d/%d virtual functions",
		pf.Device.ToString("physical"),
		len(pf.VFs),
		pf.MaxVFNum,
	)
}

func (vf *VirtualFunction) String() string {
	return fmt.Sprintf("%s index %d from %s",
		vf.Device.ToString("virtual"),
		vf.ID,
		vf.ParentAddress,
	)
}

type Info struct {
	ctx *context.Context
	// All the Physical Functions found in the host system,
	PhysicalFunctions []*PhysicalFunction `json:"physical_functions,omitempty"`
	// All the Virtual Functions found in the host system,
	// This is the very same data found navigating the `PhysicalFunctions`;
	// These pointers point back to the corresponding structs in the `PhysicalFunctions`
	// slice.
	VirtualFunctions []*VirtualFunction `json:"virtual_functions,omitempty"`
}

// New returns a pointer to an Info struct that contains information about the
// SRIOV devices on the host system.
func New(opts ...*option.Option) (*Info, error) {
	return NewWithContext(context.New(opts...))
}

// New returns a pointer to an Info struct that contains information about the
// SRIOV devices on the host system, reusing a given context.
// Use this function when you want to consume this package from another,
// ensuring the two see a coherent set of resources.
func NewWithContext(ctx *context.Context) (*Info, error) {
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *Info) String() string {
	return fmt.Sprintf(
		"sriov (%d phsyical %d virtual devices)",
		len(i.PhysicalFunctions),
		len(i.VirtualFunctions),
	)
}

// simple private struct used to encapsulate SRIOV information in a top-level
// "sriov" YAML/JSON map/object key
type sriovPrinter struct {
	Info *Info `json:"sriov,omitempty"`
}

// YAMLString returns a string with the SRIOV information formatted as YAML
// under a top-level "sriov:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, sriovPrinter{i})
}

// JSONString returns a string with the SRIOV information formatted as JSON
// under a top-level "sriov:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, sriovPrinter{i}, indent)
}
