// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

import (
	"fmt"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/pci"
	"github.com/samber/lo"
)

// PCI IDs list available at https://admin.pci-ids.ucw.cz/read/PD
const (
	pciClassProcessingAccelerator    = "12"
	pciSubclassProcessingAccelerator = "00"
	pciClassController               = "03"
	pciSubclass3DController          = "02"
	pciSubclassDisplayController     = "80"
)

var (
	acceleratorPCIClasses = map[string][]string{
		pciClassProcessingAccelerator: []string{
			pciSubclassProcessingAccelerator,
		},
		pciClassController: []string{
			pciSubclass3DController,
			pciSubclassDisplayController,
		},
	}
)

func (i *Info) load() error {
	accelDevices := make([]*AcceleratorDevice, 0)

	// get PCI devices
	pciInfo, err := pci.New(context.WithContext(i.ctx))
	if err != nil {
		i.ctx.Warn("error loading PCI information: %s", err)
		return nil
	}

	// Prepare hardware filter based in the PCI Class + Subclass
	isAccelerator := func(dev *pci.Device) bool {
		class := dev.Class.ID
		subclass := dev.Subclass.ID
		if subclasses, ok := acceleratorPCIClasses[class]; ok {
			if lo.Contains(subclasses, subclass) {
				return true
			}
		}
		return false
	}

	// This loop iterates over the list of PCI devices and filters them based on discovery criteria
	for _, device := range pciInfo.Devices {
		if !isAccelerator(device) {
			continue
		}
		if len(i.DiscoveryFilters) > 0 {
			for _, filter := range i.DiscoveryFilters {
				if validate(filter, device) {
					accelDev := &AcceleratorDevice{
						Address:   device.Address,
						PCIDevice: device,
					}
					accelDevices = append(accelDevices, accelDev)
					break
				}
			}
		} else {
			accelDev := &AcceleratorDevice{
				Address:   device.Address,
				PCIDevice: device,
			}
			accelDevices = append(accelDevices, accelDev)
		}
	}

	i.Devices = accelDevices
	return nil
}

// validate checks if a given PCI device matches the provided filter string.
//
// The filter string is expected to be in the format "VendorID:ProductID:Class+Subclass".
// Each part of the filter (VendorID, ProductID, Class+Subclass) is optional and can be
// left empty, in which case the corresponding attribute is ignored during validation.
//
// Parameters:
//   - filter: A string in the form "VendorID:ProductID:Class+Subclass", where
//     any part of the string may be empty to represent a wildcard match.
//   - device: A pointer to a `pci.Device` structure.
//
// Returns:
//   - true:  If the device matches the filter criteria (wildcards are supported).
//   - false: If the device does not match the filter criteria.
//
// Matching criteria:
//   - VendorID must match `device.Vendor.ID` if provided.
//   - ProductID must match `device.Product.ID` if provided.
//   - Class and Subclass must match the concatenated result of `device.Class.ID` and `device.Subclass.ID` if provided.
//
// Example:
//
//	filter := "8086:1234:1200"
//	device := pci.Device{Vendor: Vendor{ID: "8086"}, Product: Product{ID: "1234"}, Class: Class{ID: "12"}, Subclass: Subclass{ID: "00"}}
//	isValid := validate(filter, &device)  // returns true
//
//	filter := "8086::1200"  // Wildcard for ProductID
//	isValid := validate(filter, &device)  // returns true
//
//	filter := "::1200"  // Wildcard for ProductID and VendorID
//	isValid := validate(filter, &device)  // returns true
func validate(filter string, device *pci.Device) bool {
	ids := strings.Split(filter, ":")

	if (ids[0] == "" || ids[0] == device.Vendor.ID) &&
		(len(ids) < 2 || (ids[1] == "" || ids[1] == device.Product.ID)) &&
		(len(ids) < 3 || (ids[2] == "" || ids[2] == fmt.Sprintf("%s%s", device.Class.ID, device.Subclass.ID))) {
		return true
	}

	return false
}
