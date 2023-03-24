//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package net

import (
	"fmt"
	"encoding/json"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

type NICCapability struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
	CanEnable bool   `json:"can_enable"`
}

type NICLinkInfo struct {
	Speed                     string   `json:"speed"`
	Duplex                    string   `json:"duplex"`
	AutoNegotiation           *bool    `json:"auto-negotiation,omitempty"`
	Port                      string   `json:"port,omitempty"`
	PHYAD                     string   `json:"phyad,omitempty"`
	Transceiver               string   `json:"transceiver,omitempty"`
	MDIX                      []string `json:"mdi-x,omitempty"`
	SupportsWakeOn            string   `json:"supports_wake-on,omitempty"`
	WakeOn                    string   `json:"wake-on,omitempty"`
	LinkDetected              *bool    `json:"link_detected"`
	SupportedPorts            []string `json:"supported_ports,omitempty"`
	SupportedLinkModes        []string `json:"supported_link_modes,omitempty"`
	SupportedPauseFrameUse    *bool     `json:"supported_pause_frame_use,omitempty"`
	SupportsAutoNegotiation   *bool    `json:"supports_auto-negotiation,omitempty"`
	SupportedFECModes         []string `json:"supported_fec_modes,omitempty"`
	AdvertisedLinkModes       []string `json:"advertised_link_modes,omitempty"`
	AdvertisedPauseFrameUse   *bool    `json:"advertised_pause_frame_use,omitempty"`
	AdvertisedAutoNegotiation *bool    `json:"advertised_auto-negotiation,omitempty"`
	AdvertisedFECModes        []string `json:"advertised_fec_modes,omitempty"`
	NETIFMsgLevel             []string `json:"netif_msg_level,omitempty"`
}

func (h *NICLinkInfo) String() string {
    s, _ := json.Marshal(h)
    return string(s)
}

type NIC struct {
	Name         string           `json:"name"`
	MacAddress   string           `json:"mac_address"`
	IsVirtual    bool             `json:"is_virtual"`
	Capabilities []*NICCapability `json:"capabilities"`
	PCIAddress   *string          `json:"pci_address,omitempty"`
	LinkInfo     *NICLinkInfo     `json:"link_info"`
	// TODO(fromani): add other hw addresses (USB) when we support them
}

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

type Info struct {
	ctx  *context.Context
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
