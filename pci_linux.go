// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"os"
	"strings"
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
	info.Classes = make(map[string]*PCIClassInfo, 20)
	info.Vendors = make(map[string]*PCIVendorInfo, 200)
	info.Products = make(map[string]*PCIProductInfo, 1000)
	subclasses := make([]*PCIClassInfo, 0)
	var curClass *PCIClassInfo
	var curSubclass *PCIClassInfo
	vendorProducts := make([]*PCIProductInfo, 0)
	var curVendor *PCIVendorInfo
	var curProduct *PCIProductInfo
	var curSubsystem *PCIProductInfo
	productSubsystems := make([]*PCIProductInfo, 0)
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
				subclasses = make([]*PCIClassInfo, 0)
			}
			inClassBlock = true
			classId := string(lineBytes[2:4])
			className := string(lineBytes[6:])
			curClass = &PCIClassInfo{
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
				vendorProducts = make([]*PCIProductInfo, 0)
			}
			inClassBlock = false
			vendorId := string(lineBytes[0:4])
			vendorName := string(lineBytes[6:])
			curVendor = &PCIVendorInfo{
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
				subclassId := string(lineBytes[1:3])
				subclassName := string(lineBytes[5:])
				curSubclass = &PCIClassInfo{
					Id:   subclassId,
					Name: subclassName,
				}
				subclasses = append(subclasses, curSubclass)
			} else {
				if curProduct != nil {
					// finalize existing vendor because we found a new vendor block
					curProduct.Subsystems = productSubsystems
					productSubsystems = make([]*PCIProductInfo, 0)
				}
				productId := string(lineBytes[1:5])
				productName := string(lineBytes[7:])
				productKey := curVendor.Id + productId
				curProduct = &PCIProductInfo{
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

			} else {
				vendorId := string(lineBytes[2:6])
				subsystemId := string(lineBytes[7:11])
				subsystemName := string(lineBytes[13:])
				curSubsystem = &PCIProductInfo{
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
