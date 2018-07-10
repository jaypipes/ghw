//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"regexp"
	"strings"
)

var (
	RE_PCI_ADDRESS *regexp.Regexp = regexp.MustCompile("^(([0-9a-f]{0,4}):)?([0-9a-f]{2}):([0-9a-f]{2})\\.([0-9a-f]{1})$")
)

type PCIProgrammingInterfaceInfo struct {
	Id   string // hex-encoded PCI_ID of the programming interface
	Name string // common string name for the programming interface
}

type PCISubclassInfo struct {
	Id                    string                         // hex-encoded PCI_ID for the device subclass
	Name                  string                         // common string name for the subclass
	ProgrammingInterfaces []*PCIProgrammingInterfaceInfo // any programming interfaces this subclass might have
}

type PCIClassInfo struct {
	Id         string             // hex-encoded PCI_ID for the device class
	Name       string             // common string name for the class
	Subclasses []*PCISubclassInfo // any subclasses belonging to this class
}

// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type PCIProductInfo struct {
	VendorId   string            // vendor ID for the product
	Id         string            // hex-encoded PCI_ID for the product/model
	Name       string            // common string name of the vendor
	Subsystems []*PCIProductInfo // "subdevices" or "subsystems" for the product
}

type PCIVendorInfo struct {
	Id       string            // hex-encoded PCI_ID for the vendor
	Name     string            // common string name of the vendor
	Products []*PCIProductInfo // all top-level devices for the vendor
}

type PCIDeviceInfo struct {
	Vendor           PCIVendorInfo
	SubsystemVendor  PCIVendorInfo // optional subvendor information
	Product          PCIProductInfo
	SubsystemProduct PCIProductInfo // optional sub-device information
	Class            PCIClassInfo
	Subclass         PCIClassInfo // optional sub-class for the device
}

type PCIInfo struct {
	// hash of class ID -> class information
	Classes map[string]*PCIClassInfo
	// hash of vendor ID -> vendor information
	Vendors map[string]*PCIVendorInfo
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*PCIProductInfo
}

type PCIAddress struct {
	Domain   string
	Bus      string
	Slot     string
	Function string
}

// Given a string address, returns a complete PCIAddress struct, filled in with
// domain, bus, slot and function components. The address string may either
// be in $BUS:$SLOT.$FUNCTION (BSF) format or it can be a full PCI address
// that includes the 4-digit $DOMAIN information as well:
// $DOMAIN:$BUS:$SLOT.$FUNCTION.
//
// Returns "" if the address string wasn't a valid PCI address.
func PCIAddressFromString(address string) *PCIAddress {
	addrLowered := strings.ToLower(address)
	matches := RE_PCI_ADDRESS.FindStringSubmatch(addrLowered)
	if len(matches) == 6 {
		dom := "0000"
		if matches[1] != "" {
			dom = matches[2]
		}
		return &PCIAddress{
			Domain:   dom,
			Bus:      matches[3],
			Slot:     matches[4],
			Function: matches[5],
		}
	}
	return nil
}

func PCI() (*PCIInfo, error) {
	info := &PCIInfo{}
	err := pciFillInfo(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
