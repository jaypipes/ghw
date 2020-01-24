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

	info.NICs = ctx.nics(win32NetDescriptions)
	return nil
}

type Win32_NetworkAdapter struct {
	Description     string
	DeviceID        string
	MACAddress      string
	Manufacturer    string
	Name            string
	NetConnectionID string
	ProductName     string
	ServiceName     string
}

func (ctx *context) nics(win32NetDescriptions []Win32_NetworkAdapter) []*NIC {
	// Converting into standard structures
	nics := make([]*NIC, 0)
	for _, description := range win32NetDescriptions {
		nic := &NIC{
			Name:         ctx.netDeviceName(description),
			Vendor:       description.Manufacturer,
			MacAddress:   description.MACAddress,
			IsVirtual:    false,
			Capabilities: []*NICCapability{}, // TODO: add capabilities
		}
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
