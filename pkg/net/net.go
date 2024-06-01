//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package net

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

// NICCapability is a feature/capability of a Network Interface Controller
// (NIC)
type NICCapability struct {
	// Name is the string name for the capability, e.g.
	// "tcp-segmentation-offload"
	Name string `json:"name"`
	// IsEnabled is true if the capability is currently enabled on the NIC,
	// false otherwise.
	IsEnabled bool `json:"is_enabled"`
	// CanEnable is true if the capability can be enabled on the NIC, false
	// otherwise.
	CanEnable bool `json:"can_enable"`
}

// NIC contains information about a single Network Interface Controller (NIC).
type NIC struct {
	// Name is the string identifier the system gave this NIC.
	Name string `json:"name"`
	// MACAddress is the Media Access Control (MAC) address of this NIC.
	MACAddress string `json:"mac_address"`
	// DEPRECATED: Please use MACAddress instead.
	MacAddress string `json:"-"`
	// IsVirtual is true if the NIC is entirely virtual/emulated, false
	// otherwise.
	IsVirtual bool `json:"is_virtual"`
	// Capabilities is a slice of pointers to `NICCapability` structs
	// describing a feature/capability of this NIC.
	Capabilities []*NICCapability `json:"capabilities"`
	// PCIAddress is a pointer to the PCI address for this NIC, or nil if there
	// is no PCI address for this NIC.
	PCIAddress *string `json:"pci_address,omitempty"`
	// Speed is a string describing the link speed of this NIC, e.g. "1000Mb/s"
	Speed string `json:"speed"`
	// Duplex is a string indicating the current duplex setting of this NIC,
	// e.g. "Full"
	Duplex string `json:"duplex"`
	// SupportedLinkModes is a slice of strings containing the supported link
	// modes of this NIC, e.g. "10baseT/Half", "1000baseT/Full", etc.
	SupportedLinkModes []string `json:"supported_link_modes,omitempty"`
	// SupportedPorts is a slice of strings containing the supported physical
	// ports on this NIC, e.g. "Twisted Pair"
	SupportedPorts []string `json:"supported_ports,omitempty"`
	// SupportedFECModes is a slice of strings containing the supported Forward
	// Error Correction (FEC) modes for this NIC.
	SupportedFECModes []string `json:"supported_fec_modes,omitempty"`
	// AdvertiseLinkModes is a slice of strings containing the advertised
	// (during auto-negotiation) link modes of this NIC, e.g. "10baseT/Half",
	// "1000baseT/Full", etc.
	AdvertisedLinkModes []string `json:"advertised_link_modes,omitempty"`
	// AvertisedFECModes is a slice of strings containing the advertised
	// (during auto-negotiation) Forward Error Correction (FEC) modes for this
	// NIC.
	AdvertisedFECModes []string `json:"advertised_fec_modes,omitempty"`
	// TODO(fromani): add other hw addresses (USB) when we support them
}

// String returns a short string with information about the NIC capability.
func (nc *NICCapability) String() string {
	return fmt.Sprintf(
		"{Name:%s IsEnabled:%t CanEnable:%t}",
		nc.Name,
		nc.IsEnabled,
		nc.CanEnable,
	)
}

// String returns a short string with information about the NIC.
func (n *NIC) String() string {
	isVirtualStr := ""
	if n.IsVirtual {
		isVirtualStr = " (virtual)"
	}
	return fmt.Sprintf(
		"%s%s",
		n.Name,
		isVirtualStr,
	)
}

// Info describes all network interface controllers (NICs) in the host system.
type Info struct {
	ctx *context.Context
	// NICs is a slice of pointers to `NIC` structs describing the network
	// interface controllers (NICs) on the host system.
	NICs []*NIC `json:"nics"`
}

// New returns a pointer to an Info struct that contains information about the
// network interface controllers (NICs) on the host system
func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}
	return info, nil
}

// String returns a short string with information about the networking on the
// host system.
func (i *Info) String() string {
	return fmt.Sprintf(
		"net (%d NICs)",
		len(i.NICs),
	)
}

// simple private struct used to encapsulate net information in a
// top-level "net" YAML/JSON map/object key
type netPrinter struct {
	Info *Info `json:"network"`
}

// YAMLString returns a string with the net information formatted as YAML
// under a top-level "net:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, netPrinter{i})
}

// JSONString returns a string with the net information formatted as JSON
// under a top-level "net:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, netPrinter{i}, indent)
}
