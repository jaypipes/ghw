// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package net

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/safchain/ethtool"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func (i *Info) load() error {
	i.NICs = nics(i.ctx)
	return nil
}

func nics(ctx *context.Context) []*NIC {
	nics := make([]*NIC, 0)

	paths := linuxpath.New(ctx)
	files, err := os.ReadDir(paths.SysClassNet)
	if err != nil {
		return nics
	}

	for _, file := range files {
		filename := file.Name()
		// Ignore loopback...
		if filename == "lo" {
			continue
		}

		netPath := filepath.Join(paths.SysClassNet, filename)
		dest, _ := os.Readlink(netPath)
		isVirtual := false
		if strings.Contains(dest, "devices/virtual/net") {
			isVirtual = true
		}

		nic := &NIC{
			Name:      filename,
			IsVirtual: isVirtual,
		}

		mac := netDeviceMacAddress(paths, filename)
		nic.MacAddress = mac // keep this for backward compatibility
		nic.MACAddress = mac
		// Get speed and duplex from /sys/class/net/$DEVICE/ directory
		nic.Speed = readFile(filepath.Join(paths.SysClassNet, filename, "speed"))
		nic.Duplex = readFile(filepath.Join(paths.SysClassNet, filename, "duplex"))
		nic.Capabilities = netDeviceCapabilities(ctx, filename)
		nic.PCIAddress = netDevicePCIAddress(paths.SysClassNet, filename)

		nics = append(nics, nic)
	}
	return nics
}

func netDeviceMacAddress(paths *linuxpath.Paths, dev string) string {
	// Instead of use udevadm, we can get the device's MAC address by examing
	// the /sys/class/net/$DEVICE/address file in sysfs. However, for devices
	// that have addr_assign_type != 0, return None since the MAC address is
	// random.
	aatPath := filepath.Join(paths.SysClassNet, dev, "addr_assign_type")
	contents, err := os.ReadFile(aatPath)
	if err != nil {
		return ""
	}
	if strings.TrimSpace(string(contents)) != "0" {
		return ""
	}
	addrPath := filepath.Join(paths.SysClassNet, dev, "address")
	contents, err = os.ReadFile(addrPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}

func netDeviceCapabilities(ctx *context.Context, dev string) []*NICCapability {
	ethHandle, err := ethtool.NewEthtool()
	if err != nil {
		ctx.Warn("failed to create ethtool instance: %v", err)
		return []*NICCapability{}
	}
	defer ethHandle.Close()
	feats, err := netDeviceCapabilitiesFromEthHandle(ethHandle, dev)
	if err != nil {
		ctx.Warn(err.Error())
		return []*NICCapability{}
	}
	return feats
}

type ethtoolCollector interface {
	FeaturesWithState(intf string) (map[string]ethtool.FeatureState, error)
}

// make it mockable for test purposes
func netDeviceCapabilitiesFromEthHandle(collector ethtoolCollector, dev string) ([]*NICCapability, error) {
	feats, err := collector.FeaturesWithState(dev)
	if err != nil {
		return nil, err
	}

	caps := []*NICCapability{}
	for key, state := range feats {
		caps = append(caps, &NICCapability{
			Name:      key,
			IsEnabled: state.Active,
			CanEnable: state.Available,
		})
	}
	return caps, err
}

func netDevicePCIAddress(netDevDir, netDevName string) *string {
	// what we do here is not that hard in the end: we need to navigate the sysfs
	// up to the directory belonging to the device backing the network interface.
	// we can make few relatively safe assumptions, but the safest way is follow
	// the right links. And so we go.
	// First of all, knowing the network device name we need to resolve the backing
	// device path to its full sysfs path.
	// say we start with netDevDir="/sys/class/net" and netDevName="enp0s31f6"
	netPath := filepath.Join(netDevDir, netDevName)
	dest, err := os.Readlink(netPath)
	if err != nil {
		// bail out with empty value
		return nil
	}
	// now we have something like dest="../../devices/pci0000:00/0000:00:1f.6/net/enp0s31f6"
	// remember the path is relative to netDevDir="/sys/class/net"

	netDev := filepath.Clean(filepath.Join(netDevDir, dest))
	// so we clean "/sys/class/net/../../devices/pci0000:00/0000:00:1f.6/net/enp0s31f6"
	// leading to "/sys/devices/pci0000:00/0000:00:1f.6/net/enp0s31f6"
	// still not there. We need to access the data of the pci device. So we jump into the path
	// linked to the "device" pseudofile
	dest, err = os.Readlink(filepath.Join(netDev, "device"))
	if err != nil {
		// bail out with empty value
		return nil
	}
	// we expect something like="../../../0000:00:1f.6"

	devPath := filepath.Clean(filepath.Join(netDev, dest))
	// so we clean "/sys/devices/pci0000:00/0000:00:1f.6/net/enp0s31f6/../../../0000:00:1f.6"
	// leading to "/sys/devices/pci0000:00/0000:00:1f.6/"
	// finally here!

	// to which bus is this device connected to?
	dest, err = os.Readlink(filepath.Join(devPath, "subsystem"))
	if err != nil {
		// bail out with empty value
		return nil
	}
	// ok, this is hacky, but since we need the last *two* path components and we know we
	// are running on linux...
	if !strings.HasSuffix(dest, "/bus/pci") {
		// unsupported and unexpected bus!
		return nil
	}

	pciAddr := filepath.Base(devPath)
	return &pciAddr
}

func readFile(path string) string {
	contents, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}
