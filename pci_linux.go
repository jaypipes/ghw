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
	subclasses := make([]*PCIClassInfo, 0)
	var curClass *PCIClassInfo
	for scanner.Scan() {
		line := scanner.Text()
		// skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Lines starting with an uppercase "C" indicate a PCI top-level class
		// information block. These lines look like this:
		//
		// C 02  Network controller
		if strings.HasPrefix(line, "C") {
			if inClassBlock {
				// finalize existing class because we found a new class block
				curClass.Subclasses = subclasses
				subclasses = make([]*PCIClassInfo, 0)
			}
			inClassBlock = true
			lineBytes := []rune(line)
			classId := string(lineBytes[2:4])
			className := string(lineBytes[6:])
			fmt.Printf("Found Id: '%v' and Name: '%v'\n", classId, className)
			curClass = &PCIClassInfo{
				Id:         classId,
				Name:       className,
				Subclasses: subclasses,
			}
			info.Classes[curClass.Id] = curClass
		}
	}
	return nil
}
