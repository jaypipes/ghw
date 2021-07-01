// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sriov

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/pci"
	pciaddress "github.com/jaypipes/ghw/pkg/pci/address"
	"github.com/jaypipes/ghw/pkg/util"
)

func (info *Info) load() error {
	// SRIOV device do not have a specific class (as in "entry in /sys/class"),
	// so we need to iterate over all the PCI devices.
	pciInfo, err := pci.NewWithContext(info.ctx)
	if err != nil {
		return err
	}

	for _, dev := range pciInfo.Devices {
		err := info.scanDevice(pciInfo, dev)
		if err != nil {
			return err
		}
	}

	return nil
}

func (info *Info) scanDevice(pciInfo *pci.Info, dev *pci.Device) error {
	paths := linuxpath.New(info.ctx)
	devPath := filepath.Join(paths.SysBusPciDevices, dev.Address)

	buf, err := ioutil.ReadFile(filepath.Join(devPath, "sriov_totalvfs"))
	if err != nil {
		// is not a physfn. Since we will fill virtfn from physfn, we can give up now
		return nil
	}

	maxVFs, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		info.ctx.Warn("error reading sriov_totalvfn for %q: %v", devPath, err)
		return nil
	}

	virtFNs := findVFsFromPF(info, pciInfo, dev.Address, devPath)
	physFN := PhysicalFunction{
		Device:   info.newDevice(dev, devPath),
		MaxVFNum: maxVFs,
		VFs:      virtFNs,
	}

	info.PhysicalFunctions = append(info.PhysicalFunctions, &physFN)
	for idx := 0; idx < len(virtFNs); idx++ {
		info.VirtualFunctions = append(info.VirtualFunctions, &virtFNs[idx])
	}

	return nil
}

func findVFsFromPF(info *Info, pciInfo *pci.Info, parentAddr, parentPath string) []VirtualFunction {
	numVfs := util.SafeIntFromFile(info.ctx, filepath.Join(parentPath, "sriov_numvfs"))
	if numVfs == -1 {
		return nil
	}

	var vfs []VirtualFunction
	for vfnIdx := 0; vfnIdx < numVfs; vfnIdx++ {
		virtFn := fmt.Sprintf("virtfn%d", vfnIdx)
		vfnDest, err := os.Readlink(filepath.Join(parentPath, virtFn))
		if err != nil {
			info.ctx.Warn("error reading backing device for virtfn %q physfn %q: %v", virtFn, parentPath, err)
			return nil
		}

		vfnPath := filepath.Clean(filepath.Join(parentPath, vfnDest))
		vfnAddr := filepath.Base(vfnDest)
		vfnDev := pciInfo.GetDevice(vfnAddr)
		if vfnDev == nil {
			info.ctx.Warn("error finding the PCI device for virtfn %s physfn %s", vfnAddr, parentAddr)
			return nil
		}

		vfs = append(vfs, VirtualFunction{
			Device:        info.newDevice(vfnDev, vfnPath),
			ID:            vfnIdx,
			ParentAddress: pciaddress.FromString(parentAddr),
		})
	}
	return vfs
}

func (info *Info) newDevice(dev *pci.Device, devPath string) Device {
	// see: https://doc.dpdk.org/guides/linux_gsg/linux_drivers.html
	return Device{
		Address:    pciaddress.FromString(dev.Address),
		Interfaces: findNetworks(info.ctx, devPath),
		PCI:        dev,
	}
}

func findNetworks(ctx *context.Context, devPath string) []string {
	netPath := filepath.Join(devPath, "net")

	netEntries, err := ioutil.ReadDir(netPath)
	if err != nil {
		ctx.Warn("cannot enumerate network names for %q: %v", devPath, err)
		return nil
	}

	var networks []string
	for _, netEntry := range netEntries {
		networks = append(networks, netEntry.Name())
	}

	return networks
}
