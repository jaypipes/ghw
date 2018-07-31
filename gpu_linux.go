// +build linux
//
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
	"strconv"
	"strings"
)

const (
	PATH_SYSFS_CLASS_DRM = "/sys/class/drm"
)

func gpuFillInfo(info *GPUInfo) error {
	// In Linux, each graphics card is listed under the /sys/class/drm
	// directory as a symbolic link named "cardN", where N is a zero-based
	// index of the card in the system. "DRM" stands for Direct Rendering
	// Manager and is the Linux subsystem that is responsible for graphics I/O
	//
	// Each card may have multiple symbolic
	// links in this directory representing the interfaces from the graphics
	// card over a particular wire protocol (HDMI, DisplayPort, etc). These
	// symbolic links are named cardN-<INTERFACE_TYPE>-<DISPLAY_ID>. For
	// instance, on one of my local workstations with an NVIDIA GTX 1050ti
	// graphics card with one HDMI, one DisplayPort, and one DVI interface to
	// the card, I see the following in /sys/class/drm:
	//
	// $ ll /sys/class/drm/
	// total 0
	// drwxr-xr-x  2 root root    0 Jul 16 11:50 ./
	// drwxr-xr-x 75 root root    0 Jul 16 11:50 ../
	// lrwxrwxrwx  1 root root    0 Jul 16 11:50 card0 -> ../../devices/pci0000:00/0000:00:03.0/0000:03:00.0/drm/card0/
	// lrwxrwxrwx  1 root root    0 Jul 16 11:50 card0-DP-1 -> ../../devices/pci0000:00/0000:00:03.0/0000:03:00.0/drm/card0/card0-DP-1/
	// lrwxrwxrwx  1 root root    0 Jul 16 11:50 card0-DVI-D-1 -> ../../devices/pci0000:00/0000:00:03.0/0000:03:00.0/drm/card0/card0-DVI-D-1/
	// lrwxrwxrwx  1 root root    0 Jul 16 11:50 card0-HDMI-A-1 -> ../../devices/pci0000:00/0000:00:03.0/0000:03:00.0/drm/card0/card0-HDMI-A-1/
	//
	// In this routine, we are only interested in the first link (card0), which
	// we follow to gather information about the actual device from the PCI
	// subsystem (we query the modalias file of the PCI device's sysfs
	// directory using the `ghw.PCIInfo.GetDevice()` function.
	links, err := ioutil.ReadDir(PATH_SYSFS_CLASS_DRM)
	if err != nil {
		fmt.Fprintf(os.Stderr, `************************ WARNING ***********************************
/sys/class/drm does not exist on this system (likely the host system is a
virtual machine or container with no graphics). Therefore,
GPUInfo.GraphicsCards will be an empty array.
********************************************************************
`,
		)
		return nil
	}
	cards := make([]*GraphicsCard, 0)
	for _, link := range links {
		lname := link.Name()
		if !strings.HasPrefix(lname, "card") {
			continue
		}
		if strings.ContainsRune(lname, '-') {
			continue
		}
		// Grab the card's zero-based integer index
		lnameBytes := []byte(lname)
		cardIdx, err := strconv.Atoi(string(lnameBytes[4:]))
		if err != nil {
			cardIdx = -1
		}

		// Calculate the card's PCI address by looking at the symbolic link's
		// target
		lpath := filepath.Join(PATH_SYSFS_CLASS_DRM, lname)
		dest, err := os.Readlink(lpath)
		if err != nil {
			continue
		}
		pathParts := strings.Split(dest, "/")
		numParts := len(pathParts)
		pciAddress := pathParts[numParts-3]
		card := &GraphicsCard{
			Address: pciAddress,
			Index:   cardIdx,
		}
		cards = append(cards, card)
	}
	gpuFillNUMANodes(cards)
	gpuFillPCIDevice(cards)
	info.GraphicsCards = cards
	return nil
}

// Loops through each GraphicsCard struct and attempts to fill the DeviceInfo
// attribute with PCI device information
func gpuFillPCIDevice(cards []*GraphicsCard) {
	pci, err := PCI()
	if err != nil {
		return
	}
	for _, card := range cards {
		if card.DeviceInfo == nil {
			card.DeviceInfo = pci.GetDevice(card.Address)
		}
	}
}

// Loops through each GraphicsCard struct and find which NUMA nodes the card is
// affined to, setting the GraphicsCard.Nodes field accordingly. If the host
// system is not a NUMA system, the Nodes field will be set to an empty array
// of Node pointers.
func gpuFillNUMANodes(cards []*GraphicsCard) {
	topo, err := Topology()
	if err != nil {
		for _, card := range cards {
			if topo.Architecture != NUMA {
				card.Nodes = make([]*TopologyNode, 0)
			}
		}
		return
	}
	for _, card := range cards {
		if topo.Architecture != NUMA {
			card.Nodes = make([]*TopologyNode, 0)
			continue
		}
		// Each graphics card on a NUMA system will have a pseudo-file
		// called /sys/class/drm/card$CARD_INDEX/device/numa_node which
		// contains a comma-separated list of NUMA nodes that the card is
		// affined to
		cardIndexStr := strconv.Itoa(card.Index)
		fpath := filepath.Join(
			PATH_SYSFS_CLASS_DRM,
			"card"+cardIndexStr,
			"device",
			"numa_node",
		)
		numaContents, err := ioutil.ReadFile(fpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, `************************ WARNING ***********************************
Unable to read numa_nodes descriptor file on this system.
Setting graphics card's Nodes attribute to empty array.
********************************************************************
`,
			)
			card.Nodes = make([]*TopologyNode, 0)
			continue
		}
		cardNodes := make([]*TopologyNode, 0)
		nodeIndexes := strings.Split(string(numaContents), ",")
		for _, nodeIndex := range nodeIndexes {
			for _, node := range topo.Nodes {
				nodeIndexInt, _ := strconv.Atoi(nodeIndex)
				if nodeIndexInt == int(node.Id) {
					cardNodes = append(cardNodes, node)
				}
			}
		}
		card.Nodes = cardNodes
	}
}
