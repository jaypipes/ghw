//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "io"

type PciClassInfo struct {
	Id   string // hex-encoded PCI_ID for the device class
	Name string // common string name for the class
}

type PciVendorInfo struct {
	Id   string // hex-encoded PCI_ID for the vendor
	Name string // common string name of the vendor
}

// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type PciProductInfo struct {
	Id   string // hex-encoded PCI_ID for the product/model
	Name string // common string name of the vendor
}

type PciDeviceInfo struct {
	Vendor           PciVendorInfo
	SubsystemVendor  PciVendorInfo // optional subvendor information
	Product          PciProductInfo
	SubsystemProduct PciProductInfo // optional sub-device information
	Class            PciClassInfo
	Subclass         PciClassInfo // optional sub-class for the device
}

// interface for a thing that can read a pci.ids database and return
// information about vendors, devices and classes
type pciDb interface {
	// Loads the database by reading from the supplied reader and populating
	// some internal state about PCI classes, devices and vendors
	loadFrom(*io.Reader)
	// Returns a pointer to a PciDevice struct given a vendor ID string
	GetVendorInfo(string) *PciVendorInfo
	// Returns a pointer to a PciProductInfo struct given a vendor and device
	// (product) ID string
	GetProductInfo(string, string) *PciProductInfo
	// Returns a pointer to a PciDeviceInfo struct given a string representing
	// the address of the device on the bus
	GetDeviceInfo(string) *PciDeviceInfo
}
