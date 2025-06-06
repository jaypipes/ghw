// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

// GetSRIOVDevices returns only the PCI devices that are
// Single Root I/O Virtualization (SR-IOV) capable -- either
// physical of virtual functions.
func (i *Info) GetSRIOVDevices() []*Device {
	res := []*Device{}
	for _, dev := range i.Devices {
		if dev.Function == nil {
			continue
		}
		res = append(res, dev)
	}
	return res
}

func (info *Info) fillSRIOVDevices() error {
	for _, dev := range info.Devices {
		isPF, err := info.fillPhysicalFunctionForDevice(dev)
		if !isPF {
			// not a physical function, nothing to do
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (info *Info) fillPhysicalFunctionForDevice(dev *Device) (bool, error) {
	paths := linuxpath.New(info.ctx)
	devPath := filepath.Join(paths.SysBusPciDevices, dev.Address)

	buf, err := os.ReadFile(filepath.Join(devPath, "sriov_totalvfs"))
	if err != nil {
		// is not a physfn. Since we will fill virtfn from physfn, we can give up now
		// note we intentionally swallow the error.
		return false, nil
	}

	maxVFs, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return true, fmt.Errorf("cannot reading sriov_totalvfn: %w", err)
	}

	pf := &Function{
		MaxVirtual: maxVFs,
	}
	err = info.fillVirtualFunctionsForPhysicalFunction(pf, devPath)
	if err != nil {
		return true, fmt.Errorf("cannot inspect VFs: %w", err)
	}
	dev.Function = pf
	return true, nil
}

func (info *Info) fillVirtualFunctionsForPhysicalFunction(parentFn *Function, parentPath string) error {
	numVfs := util.SafeIntFromFile(info.ctx, filepath.Join(parentPath, "sriov_numvfs"))
	if numVfs == -1 {
		return fmt.Errorf("invalid number of virtual functions: %v", numVfs)
	}

	var vfs []*Function
	for vfnIdx := 0; vfnIdx < numVfs; vfnIdx++ {
		virtFn := fmt.Sprintf("virtfn%d", vfnIdx)
		vfnDest, err := os.Readlink(filepath.Join(parentPath, virtFn))
		if err != nil {
			return fmt.Errorf("error reading backing device for virtfn %q: %w", virtFn, err)
		}

		vfnAddr := filepath.Base(vfnDest)
		vfnDev := info.GetDevice(vfnAddr)
		if vfnDev == nil {
			return fmt.Errorf("error finding the PCI device for virtfn %s", vfnAddr)
		}

		// functions must be ordered by their index
		vf := &Function{
			Parent: parentFn,
		}

		vfs = append(vfs, vf)
		vfnDev.Function = vf
	}

	parentFn.Virtual = vfs
	return nil
}
