// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/jaypipes/pcidb"
)

func (ctx *context) pciFillInfo(info *PCIInfo) error {
	//info.Classes = db.Classes
	//info.Vendors = db.Vendors
	//info.Products = db.Products
	info.ctx = ctx
	return nil
}

type Win32_PnPEntity struct {
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	Manufacturer      string
	Name              string
	PNPClass          string
	PNPDeviceID       string
}

// GetDevice returns a pointer to a PCIDevice struct that describes the PCI
// device at the requested address. If no such device could be found, returns
// nil
func (info *PCIInfo) GetDevice(address string) *PCIDevice {
	// Backslashing address as requested by JSON and VMI query: https://docs.microsoft.com/en-us/windows/win32/wmisdk/where-clause
	var queryAddress = strings.Replace(address, "\\", `\\`, -1)
	// Preparing default structure
	var device = &PCIDevice{
		Address: queryAddress,
		Vendor: &pcidb.Vendor{
			ID:       UNKNOWN,
			Name:     UNKNOWN,
			Products: []*pcidb.Product{},
		},
		Subsystem: &pcidb.Product{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subsystems: []*pcidb.Product{},
		},
		Product: &pcidb.Product{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subsystems: []*pcidb.Product{},
		},
		Class: &pcidb.Class{
			ID:         UNKNOWN,
			Name:       UNKNOWN,
			Subclasses: []*pcidb.Subclass{},
		},
		Subclass: &pcidb.Subclass{
			ID:                    UNKNOWN,
			Name:                  UNKNOWN,
			ProgrammingInterfaces: []*pcidb.ProgrammingInterface{},
		},
		ProgrammingInterface: &pcidb.ProgrammingInterface{
			ID:   UNKNOWN,
			Name: UNKNOWN,
		},
	}
	// Getting disk drives from WMI
	var win32PnPDescriptions []Win32_PnPEntity
	q1 := wmi.CreateQuery(&win32PnPDescriptions, "WHERE PNPDeviceID='"+queryAddress+"'")
	if err := wmi.Query(q1, &win32PnPDescriptions); err != nil {
		return device
	}
	// Converting into standard structures
	device.Vendor.ID = win32PnPDescriptions[0].Manufacturer
	device.Vendor.Name = win32PnPDescriptions[0].Manufacturer
	device.Product.ID = win32PnPDescriptions[0].Name
	device.Product.Name = win32PnPDescriptions[0].Description
	return device
}
