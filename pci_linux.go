// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	PATH_SYSFS_PCI_DEVICES = "/sys/bus/pci/devices"
)

var (
	pciIdsFilePaths = []string{
		"/usr/share/hwdata/pci.ids",
		"/usr/share/misc/pci.ids",
	}
)

func pciFillInfo(info *PCIInfo) error {
	for _, fp := range pciIdsFilePaths {
		if _, err := os.Stat(fp); err != nil {
			continue
		}
		f, err := os.Open(fp)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		err = parsePCIIdsFile(info, scanner)
		if err == nil {
			break
		}
	}
	return nil
}

func parsePCIIdsFile(info *PCIInfo, scanner *bufio.Scanner) error {
	inClassBlock := false
	info.Classes = make(map[string]*PCIClass, 20)
	info.Vendors = make(map[string]*PCIVendor, 200)
	info.Products = make(map[string]*PCIProduct, 1000)
	subclasses := make([]*PCISubclass, 0)
	progIfaces := make([]*PCIProgrammingInterface, 0)
	var curClass *PCIClass
	var curSubclass *PCISubclass
	var curProgIface *PCIProgrammingInterface
	vendorProducts := make([]*PCIProduct, 0)
	var curVendor *PCIVendor
	var curProduct *PCIProduct
	var curSubsystem *PCIProduct
	productSubsystems := make([]*PCIProduct, 0)
	for scanner.Scan() {
		line := scanner.Text()
		// skip comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lineBytes := []rune(line)

		// Lines starting with an uppercase "C" indicate a PCI top-level class
		// information block. These lines look like this:
		//
		// C 02  Network controller
		if lineBytes[0] == 'C' {
			if curClass != nil {
				// finalize existing class because we found a new class block
				curClass.Subclasses = subclasses
				subclasses = make([]*PCISubclass, 0)
			}
			inClassBlock = true
			classId := string(lineBytes[2:4])
			className := string(lineBytes[6:])
			curClass = &PCIClass{
				Id:         classId,
				Name:       className,
				Subclasses: subclasses,
			}
			info.Classes[curClass.Id] = curClass
			continue
		}

		// Lines not beginning with an uppercase "C" or a TAB character
		// indicate a top-level vendor information block. These lines look like
		// this:
		//
		// 0a89  BREA Technologies Inc
		if lineBytes[0] != '\t' {
			if curVendor != nil {
				// finalize existing vendor because we found a new vendor block
				curVendor.Products = vendorProducts
				vendorProducts = make([]*PCIProduct, 0)
			}
			inClassBlock = false
			vendorId := string(lineBytes[0:4])
			vendorName := string(lineBytes[6:])
			curVendor = &PCIVendor{
				Id:       vendorId,
				Name:     vendorName,
				Products: vendorProducts,
			}
			info.Vendors[curVendor.Id] = curVendor
			continue
		}

		// Lines beginning with only a single TAB character are *either* a
		// subclass OR are a device information block. If we're in a class
		// block (i.e. the last parsed block header was for a PCI class), then
		// we parse a subclass block. Otherwise, we parse a device information
		// block.
		//
		// A subclass information block looks like this:
		//
		// \t00  Non-VGA unclassified device
		//
		// A device information block looks like this:
		//
		// \t0002  PCI to MCA Bridge
		if len(lineBytes) > 1 && lineBytes[1] != '\t' {
			if inClassBlock {
				if curSubclass != nil {
					// finalize existing subclass because we found a new subclass block
					curSubclass.ProgrammingInterfaces = progIfaces
					progIfaces = make([]*PCIProgrammingInterface, 0)
				}
				subclassId := string(lineBytes[1:3])
				subclassName := string(lineBytes[5:])
				curSubclass = &PCISubclass{
					Id:   subclassId,
					Name: subclassName,
					ProgrammingInterfaces: progIfaces,
				}
				subclasses = append(subclasses, curSubclass)
			} else {
				if curProduct != nil {
					// finalize existing product because we found a new product block
					curProduct.Subsystems = productSubsystems
					productSubsystems = make([]*PCIProduct, 0)
				}
				productId := string(lineBytes[1:5])
				productName := string(lineBytes[7:])
				productKey := curVendor.Id + productId
				curProduct = &PCIProduct{
					VendorId: curVendor.Id,
					Id:       productId,
					Name:     productName,
				}
				vendorProducts = append(vendorProducts, curProduct)
				info.Products[productKey] = curProduct
			}
		} else {
			// Lines beginning with two TAB characters are *either* a subsystem
			// (subdevice) OR are a programming interface for a PCI device
			// subclass. If we're in a class block (i.e. the last parsed block
			// header was for a PCI class), then we parse a programming
			// interface block, otherwise we parse a subsystem block.
			//
			// A programming interface block looks like this:
			//
			// \t\t00  UHCI
			//
			// A subsystem block looks like this:
			//
			// \t\t0e11 4091  Smart Array 6i
			if inClassBlock {
				progIfaceId := string(lineBytes[2:4])
				progIfaceName := string(lineBytes[6:])
				curProgIface = &PCIProgrammingInterface{
					Id:   progIfaceId,
					Name: progIfaceName,
				}
				progIfaces = append(progIfaces, curProgIface)
			} else {
				vendorId := string(lineBytes[2:6])
				subsystemId := string(lineBytes[7:11])
				subsystemName := string(lineBytes[13:])
				curSubsystem = &PCIProduct{
					VendorId: vendorId,
					Id:       subsystemId,
					Name:     subsystemName,
				}
				productSubsystems = append(productSubsystems, curSubsystem)
			}
		}
	}
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
		vendor = &PCIVendor{
			Id:       vendorId,
			Name:     "UNKNOWN",
			Products: []*PCIProduct{},
		}
	}

	// Find the product
	product := info.Products[vendorId+productId]
	if product == nil {
		product = &PCIProduct{
			Id:         productId,
			Name:       "UNKNOWN",
			Subsystems: []*PCIProduct{},
		}
	}

	// Find the subsystem information
	subvendor := info.Vendors[subvendorId]
	var subsystem *PCIProduct
	if subvendor != nil && product != nil {
		for _, p := range product.Subsystems {
			if p.Id == subproductId {
				subsystem = p
			}
		}
	}
	if subsystem == nil {
		subsystem = &PCIProduct{
			VendorId: subvendorId,
			Id:       subproductId,
			Name:     "UNKNOWN",
		}
	}

	// Find the class and subclass
	class := info.Classes[classId]
	var subclass *PCISubclass
	if class != nil {
		for _, sc := range class.Subclasses {
			if sc.Id == subclassId {
				subclass = sc
			}
		}
	} else {
		class = &PCIClass{
			Id:         classId,
			Name:       "UNKNOWN",
			Subclasses: []*PCISubclass{},
		}
	}

	// Find the programming interface
	var progIface *PCIProgrammingInterface
	if subclass != nil {
		for _, pi := range subclass.ProgrammingInterfaces {
			if pi.Id == progIfaceId {
				progIface = pi
			}
		}
	} else {
		subclass = &PCISubclass{
			Id:   subclassId,
			Name: "UNKNOWN",
			ProgrammingInterfaces: []*PCIProgrammingInterface{},
		}
	}

	if progIface == nil {
		progIface = &PCIProgrammingInterface{
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
