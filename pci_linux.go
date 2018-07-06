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
	info.Classes = make(map[string]*PCIClassInfo, 20)
	info.Vendors = make(map[string]*PCIVendorInfo, 200)
	subclasses := make([]*PCIClassInfo, 0)
	var curClass *PCIClassInfo
	var curVendor *PCIVendorInfo
	for scanner.Scan() {
		line := scanner.Text()
		// skip comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Lines starting with an uppercase "C" indicate a PCI top-level class
		// information block. These lines look like this:
		//
		// C 02  Network controller
		if strings.HasPrefix(line, "C") {
			if curClass != nil {
				// finalize existing class because we found a new class block
				curClass.Subclasses = subclasses
				subclasses = make([]*PCIClassInfo, 0)
			}
			lineBytes := []rune(line)
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
		// indicate a top-level vendor information block. These lines look like this:
		//
		// 0a89  BREA Technologies Inc
		if !strings.HasPrefix(line, "\t") {
			if curVendor != nil {
				// finalize existing vendor because we found a new vendor block
			}
			lineBytes := []rune(line)
			vendorId := string(lineBytes[0:4])
			vendorName := string(lineBytes[6:])
			curVendor = &PCIVendorInfo{
				Id:   vendorId,
				Name: vendorName,
			}
			info.Vendors[curVendor.Id] = curVendor
			continue
		}
	}
	return nil
}
