// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
)

func (ctx *context) netFillInfo(info *NetworkInfo) error {
	// Getting info from WMI
	var win32NetDescriptions []Win32_NetworkAdapter
	q1 := wmi.CreateQuery(&win32NetDescriptions, "")
	if err := wmi.Query(q1, &win32NetDescriptions); err != nil {
		return err
	}

	var win32NetConfigurationDescriptions []Win32_NetworkAdapterConfiguration
	q2 := wmi.CreateQuery(&win32NetConfigurationDescriptions, "")
	if err := wmi.Query(q2, &win32NetConfigurationDescriptions); err != nil {
		return err
	}

	var win32IP4RouteTableDescriptions []Win32_IP4RouteTable
	q3 := wmi.CreateQuery(&win32IP4RouteTableDescriptions, "")
	if err := wmi.Query(q3, &win32IP4RouteTableDescriptions); err != nil {
		return err
	}

	info.NICs = ctx.nics(win32NetDescriptions, win32NetConfigurationDescriptions, win32IP4RouteTableDescriptions)
	return nil
}

type Win32_NetworkAdapter struct {
	Description     string
	DeviceID        string
	Index           uint32
	InterfaceIndex  uint32
	MACAddress      string
	Manufacturer    string
	Name            string
	NetConnectionID string
	ProductName     string
	ServiceName     string
}

type Win32_NetworkAdapterConfiguration struct {
	Caption          string
	Description      string
	DefaultIPGateway []string
	DHCPEnabled      bool
	Index            uint32
	InterfaceIndex   uint32
	IPAddress        []string
}

type Win32_IP4RouteTable struct {
	Caption        string
	Description    string
	Destination    string
	Information    string
	InterfaceIndex int32
	Mask           string
	Metric1        int32
	Metric2        int32
	Metric3        int32
	Metric4        int32
	Metric5        int32
	Name           string
	NextHop        string
	Protocol       uint32
	Status         string
	Type           uint32
}

func (ctx *context) nics(win32NetDescriptions []Win32_NetworkAdapter, win32NetConfigurationDescriptions []Win32_NetworkAdapterConfiguration, win32IP4RouteTableDescriptions []Win32_IP4RouteTable) []*NIC {
	// Converting into standard structures
	nics := make([]*NIC, 0)
	for _, nicDescription := range win32NetDescriptions {
		nic := &NIC{
			Name:         ctx.netDeviceName(nicDescription),
			Vendor:       nicDescription.Manufacturer,
			MacAddress:   nicDescription.MACAddress,
			IsVirtual:    false,
			Capabilities: []*NICCapability{}, // TODO: add capabilities
		}
		//
		for _, configDescription := range win32NetConfigurationDescriptions {
			// Looking for configurations
			if nicDescription.InterfaceIndex == configDescription.InterfaceIndex {
				ipv4, ipv6 := ctx.netConfigIP(configDescription.IPAddress)
				var configuration = &NICConfiguration{
					DHCPenabled: configDescription.DHCPEnabled,
					IPv4:        ipv4,
					IPv6:        ipv6,
				}
				// Looking for gateway
				for _, routeDescription := range win32IP4RouteTableDescriptions {
					if nicDescription.InterfaceIndex == uint32(routeDescription.InterfaceIndex) {
						if routeDescription.Destination == "0.0.0.0" && routeDescription.Mask == "0.0.0.0" {
							configuration.Gateway = routeDescription.NextHop
							break
						}
					}
				}
				// Appending configuration to NIC configurations
				nic.Configurations = append(nic.Configurations, configuration)
			}
		}
		// Appenging NIC to NICs
		nics = append(nics, nic)
	}

	return nics
}

func (ctx *context) netDeviceName(description Win32_NetworkAdapter) string {
	var name string
	if strings.TrimSpace(description.NetConnectionID) != "" {
		name = description.NetConnectionID + " - " + description.Description
	} else {
		name = description.Description
	}
	return name
}

func (ctx *context) netConfigIP(IPs []string) (string, string) {
	var IPv4 string
	var IPv6 string
	if len(IPs) > 0 {
		IPv4 = IPs[0]
	}
	if len(IPs) > 1 {
		IPv6 = IPs[1]
	}
	return IPv4, IPv6
}
