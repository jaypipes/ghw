// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	PathSysClassNet = "/sys/class/net"
)

func netFillInfo(info *NetworkInfo) error {
	info.NICs = NICs()
	return nil
}

func NICs() []*NIC {
	nics := make([]*NIC, 0)

	files, err := ioutil.ReadDir(PathSysClassNet)
	if err != nil {
		return nics
	}

	for _, file := range files {
		filename := file.Name()
		// Ignore loopback...
		if filename == "lo" {
			continue
		}

		netPath := filepath.Join(PathSysClassNet, filename)
		dest, _ := os.Readlink(netPath)
		isVirtual := false
		if strings.Contains(dest, "virtio") {
			isVirtual = true
		}

		nic := &NIC{
			Name:      filename,
			IsVirtual: isVirtual,
		}

		mac := netDeviceMacAddress(filename)
		nic.MacAddress = mac
		nics = append(nics, nic)
	}

	return nics
}

func netDeviceMacAddress(dev string) string {
	// Instead of use udevadm, we can get the device's MAC address by examing
	// the /sys/class/net/$DEVICE/address file in sysfs. However, for devices
	// that have addr_assign_type != 0, return None since the MAC address is
	// random.
	aatPath := filepath.Join(PathSysClassNet, dev, "addr_assign_type")
	contents, err := ioutil.ReadFile(aatPath)
	if err != nil {
		return ""
	}
	if strings.TrimSpace(string(contents)) != "0" {
		return ""
	}
	addrPath := filepath.Join(PathSysClassNet, dev, "address")
	contents, err = ioutil.ReadFile(addrPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}
