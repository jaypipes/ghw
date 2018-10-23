//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jaypipes/pcidb"
)

var (
	regexPCIAddress *regexp.Regexp = regexp.MustCompile(
		`^(([0-9a-f]{0,4}):)?([0-9a-f]{2}):([0-9a-f]{2})\.([0-9a-f]{1})$`,
	)
)

type PCIDevice struct {
	Address              string // The PCI address of the device
	Vendor               *pcidb.PCIVendor
	Product              *pcidb.PCIProduct
	Subsystem            *pcidb.PCIProduct // optional subvendor/sub-device information
	Class                *pcidb.PCIClass
	Subclass             *pcidb.PCISubclass             // optional sub-class for the device
	ProgrammingInterface *pcidb.PCIProgrammingInterface // optional programming interface
}

func (di *PCIDevice) String() string {
	vendorName := UNKNOWN
	if di.Vendor != nil {
		vendorName = di.Vendor.Name
	}
	productName := UNKNOWN
	if di.Product != nil {
		productName = di.Product.Name
	}
	className := UNKNOWN
	if di.Class != nil {
		className = di.Class.Name
	}
	return fmt.Sprintf(
		"%s -> class: '%s' vendor: '%s' product: '%s'",
		di.Address,
		className,
		vendorName,
		productName,
	)
}

type PCIInfo struct {
	// hash of class ID -> class information
	Classes map[string]*pcidb.PCIClass
	// hash of vendor ID -> vendor information
	Vendors map[string]*pcidb.PCIVendor
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*pcidb.PCIProduct
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
	matches := regexPCIAddress.FindStringSubmatch(addrLowered)
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
