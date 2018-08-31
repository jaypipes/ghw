// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/pcidb"
)

const (
	PATH_SYSFS_PCI_DEVICES = "/sys/bus/pci/devices"
)

func pciFillInfo(info *PCIInfo) error {
	db, err := pcidb.New()
	if err != nil {
		return err
	}
	info.Classes = db.Classes
	info.Vendors = db.Vendors
	info.Products = db.Products
	return nil
}

func getPCIDeviceModaliasPath(address string) string {
	pciAddr := PCIAddressFromString(address)
	if pciAddr == nil {
		return ""
	}
	return filepath.Join(
		PATH_SYSFS_PCI_DEVICES,
		pciAddr.Domain+":"+pciAddr.Bus+":"+pciAddr.Slot+"."+pciAddr.Function,
		"modalias",
	)
}

// Returns a pointer to a PCIDevice struct that describes the PCI device at
// the requested address. If no such device could be found, returns nil
func (info *PCIInfo) GetDevice(address string) *PCIDevice {
	fp := getPCIDeviceModaliasPath(address)
	if fp == "" {
		return nil
	}
	if _, err := os.Stat(fp); err != nil {
		return nil
	}
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil
	}
	// The modalias file is an encoded file that looks like this:
	//
	// $ cat /sys/devices/pci0000\:00/0000\:00\:03.0/0000\:03\:00.0/modalias
	// pci:v000010DEd00001C82sv00001043sd00008613bc03sc00i00
	//
	// It is interpreted like so:
	//
	// pci: -- ignore
	// v000010DE -- PCI vendor ID
	// d00001C82 -- PCI device ID (the product/model ID)
	// sv00001043 -- PCI subsystem vendor ID
	// sd00008613 -- PCI subsystem device ID (subdevice product/model ID)
	// bc03 -- PCI base class
	// sc00 -- PCI subclass
	// i00 -- programming interface
	vendorId := strings.ToLower(string(data[9:13]))
	productId := strings.ToLower(string(data[18:22]))
	subvendorId := strings.ToLower(string(data[28:32]))
	subproductId := strings.ToLower(string(data[38:42]))
	classId := string(data[44:46])
	subclassId := string(data[48:50])
	progIfaceId := string(data[51:53])

	// Find the vendor
	vendor := info.Vendors[vendorId]
	if vendor == nil {
		vendor = &pcidb.PCIVendor{
			Id:       vendorId,
			Name:     "UNKNOWN",
			Products: []*pcidb.PCIProduct{},
		}
	}

	// Find the product
	product := info.Products[vendorId+productId]
	if product == nil {
		product = &pcidb.PCIProduct{
			Id:         productId,
			Name:       "UNKNOWN",
			Subsystems: []*pcidb.PCIProduct{},
		}
	}

	// Find the subsystem information
	subvendor := info.Vendors[subvendorId]
	var subsystem *pcidb.PCIProduct
	if subvendor != nil && product != nil {
		for _, p := range product.Subsystems {
			if p.Id == subproductId {
				subsystem = p
			}
		}
	}
	if subsystem == nil {
		subsystem = &pcidb.PCIProduct{
			VendorId: subvendorId,
			Id:       subproductId,
			Name:     "UNKNOWN",
		}
	}

	// Find the class and subclass
	class := info.Classes[classId]
	var subclass *pcidb.PCISubclass
	if class != nil {
		for _, sc := range class.Subclasses {
			if sc.Id == subclassId {
				subclass = sc
			}
		}
	} else {
		class = &pcidb.PCIClass{
			Id:         classId,
			Name:       "UNKNOWN",
			Subclasses: []*pcidb.PCISubclass{},
		}
	}

	// Find the programming interface
	var progIface *pcidb.PCIProgrammingInterface
	if subclass != nil {
		for _, pi := range subclass.ProgrammingInterfaces {
			if pi.Id == progIfaceId {
				progIface = pi
			}
		}
	} else {
		subclass = &pcidb.PCISubclass{
			Id:   subclassId,
			Name: "UNKNOWN",
			ProgrammingInterfaces: []*pcidb.PCIProgrammingInterface{},
		}
	}

	if progIface == nil {
		progIface = &pcidb.PCIProgrammingInterface{
			Id:   progIfaceId,
			Name: "UNKNOWN",
		}
	}

	return &PCIDevice{
		Address:              address,
		Vendor:               vendor,
		Subsystem:            subsystem,
		Product:              product,
		Class:                class,
		Subclass:             subclass,
		ProgrammingInterface: progIface,
	}
}

// Returns a list of pointers to PCIDevice structs present on the host system
func (info *PCIInfo) ListDevices() []*PCIDevice {
	devs := make([]*PCIDevice, 0)
	// We scan the /sys/bus/pci/devices directory which contains a collection
	// of symlinks. The names of the symlinks are all the known PCI addresses
	// for the host. For each address, we grab a *PCIDevice matching the
	// address and append to the returned array.
	links, err := ioutil.ReadDir(PATH_SYSFS_PCI_DEVICES)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read /sys/bus/pci/devices")
		return nil
	}
	var dev *PCIDevice
	for _, link := range links {
		addr := link.Name()
		dev = info.GetDevice(addr)
		if dev == nil {
			fmt.Fprintf(os.Stderr, "error: failed to get device information for PCI address %s\n", addr)
		} else {
			devs = append(devs, dev)
		}
	}
	return devs
}
