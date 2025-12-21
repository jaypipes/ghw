//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/pci"
	"github.com/jaypipes/ghw/pkg/snapshot"
	"github.com/jaypipes/ghw/pkg/util"

	"github.com/jaypipes/ghw/testdata"
)

type pciTestCase struct {
	addr       string
	parentAddr string
	node       int
	revision   string
	driver     string
	iommuGroup string
}

// nolint: gocyclo
func TestPCINUMANode(t *testing.T) {
	info := pciTestSetupXeon(t)

	tCases := []pciTestCase{
		{
			addr: "0000:07:03.0",
			// -1 is actually what we get out of the box on the snapshotted box
			node: -1,
		},
		{
			addr: "0000:05:11.0",
			node: 0,
		},
		{
			addr: "0000:05:00.1",
			node: 1,
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%d)", tCase.addr, tCase.node), func(t *testing.T) {
			dev := info.GetDevice(tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Node == nil {
				if tCase.node != -1 {
					addrs := []string{}
					for _, d := range info.Devices {
						addrs = append(addrs, d.Address)
					}
					msg := fmt.Sprintf("address: %q device addresses: %v", tCase.addr, addrs)
					t.Error(msg)
					t.Fatalf("got nil numa NODE for address %q", tCase.addr)
				}
			} else {
				if dev.Node.ID != tCase.node {
					t.Errorf("got NUMA node info %#v, expected on node %d", dev.Node, tCase.node)
				}
			}
		})
	}
}

// nolint: gocyclo
func TestPCIDeviceRevision(t *testing.T) {
	info := pciTestSetupXeon(t)

	var tCases []pciTestCase = []pciTestCase{
		{
			addr:     "0000:07:03.0",
			revision: "0x0a",
		},
		{
			addr:     "0000:05:00.0",
			revision: "0x01",
		},
	}
	for _, tCase := range tCases {
		t.Run(tCase.addr, func(t *testing.T) {
			dev := info.GetDevice(tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Revision != tCase.revision {
				t.Errorf("device %q got revision %q expected %q", tCase.addr, dev.Revision, tCase.revision)
			}
		})
	}
}

// nolint: gocyclo
func TestPCIParent(t *testing.T) {
	info := pciTestSetupI7(t)
	tCases := []pciTestCase{
		{
			addr:       "0000:04:00.0",
			parentAddr: "0000:00:06.0",
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%s)", tCase.addr, tCase.parentAddr), func(t *testing.T) {
			dev := info.GetDevice(tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.ParentAddress != tCase.parentAddr {
				t.Errorf("got parent %q expected %q", dev.ParentAddress, tCase.parentAddr)
			}
		})
	}
}

// nolint: gocyclo
func TestPCIIommuGroup(t *testing.T) {
	info := pciTestSetupI7(t)
	tCases := []pciTestCase{
		{
			addr:       "0000:00:1f.0",
			iommuGroup: "13",
		},
		{
			addr:       "0000:00:1f.5",
			iommuGroup: "13",
		},
		{
			addr:       "0000:04:00.0",
			iommuGroup: "14",
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%s)", tCase.addr, tCase.iommuGroup), func(t *testing.T) {
			dev := info.GetDevice(tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.IOMMUGroup != tCase.iommuGroup {
				t.Errorf("got iommu_group %q expected %q", dev.IOMMUGroup, tCase.iommuGroup)
			}
		})
	}
}

// nolint: gocyclo
func TestPCIDriver(t *testing.T) {
	info := pciTestSetupXeon(t)

	tCases := []pciTestCase{
		{
			addr:   "0000:07:03.0",
			driver: "mgag200",
		},
		{
			addr:   "0000:05:11.0",
			driver: "igbvf",
		},
		{
			addr:   "0000:05:00.1",
			driver: "igb",
		},
	}
	for _, tCase := range tCases {
		t.Run(fmt.Sprintf("%s (%s)", tCase.addr, tCase.driver), func(t *testing.T) {
			dev := info.GetDevice(tCase.addr)
			if dev == nil {
				t.Fatalf("got nil device for address %q", tCase.addr)
			}
			if dev.Driver != tCase.driver {
				t.Errorf("got driver %q expected %q", dev.Driver, tCase.driver)
			}
		})
	}
}

func TestPCIMarshalJSON(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	dev := info.ParseDevice("0000:3c:00.0", "pci:v0000144Dd0000A804sv0000144Dsd0000A801bc01sc08i02\n")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias")
	}
	s := marshal.SafeJSON(dev, true)
	if s == "" {
		t.Fatalf("Error marshalling device: %v", dev)
	}
}

// the sriov-device-plugin code has a test like this
func TestPCIMalformedModalias(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	var dev *pci.Device
	dev = info.ParseDevice("0000:00:01.0", "pci:junk")
	if dev != nil {
		t.Fatalf("Parsed successfully junk data")
	}

	dev = info.ParseDevice("0000:00:01.0", "pci:v00008086d00005916sv000017AAsd0000224Bbc03sc00i00extrajunkextradataextraextra")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias with extra data")
	}
}

func pciTestSetupXeon(t *testing.T) *pci.Info {
	const snapshotFilename = "linux-amd64-intel-xeon-L5640.tar.gz"
	return pciTestSetup(t, snapshotFilename)
}

func pciTestSetupI7(t *testing.T) *pci.Info {
	const snapshotFilename = "linux-amd64-intel-i7-1270P.tar.gz"
	return pciTestSetup(t, snapshotFilename)
}

func pciTestSetup(t *testing.T, snapshotFilename string) *pci.Info {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	snapshotPath := filepath.Join(testdataPath, snapshotFilename)
	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(snapshotPath, unpackDir)
	if err != nil {
		t.Fatal(err)
	}

	// from now on we use constants reflecting the content of the snapshot we requested,
	// which we reviewed beforehand. IOW, you need to know the content of the
	// snapshot to fully understand this test. Inspect it using
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc topology

	info, err := pci.New(option.WithChroot(unpackDir))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil PCIInfo, but got nil")
	}
	return info
}

// we have this test in pci_linux_test.go (and not in pci_test.go) because `pciFillInfo` is implemented
// only on linux; so having it in the platform-independent tests would lead to false negatives.
func TestPCIMarshalUnmarshal(t *testing.T) {
	data, err := pci.New(option.WithNullAlerter())
	if err != nil {
		t.Fatalf("Expected no error creating pci.Info, but got %v", err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Expected no error marshaling pci.Info, but got %v", err)
	}

	var topo *pci.Info

	err = json.Unmarshal(jdata, &topo)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling pci.Info, but got %v", err)
	}
}

func TestPCIModaliasWithUpperCaseClassID(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	dev := info.ParseDevice("0000:00:1f.4", "pci:v00008086d00009D23sv00001028sd000007EAbc0Csc05i00\n")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias")
	}
	if dev.Class.Name == util.UNKNOWN {
		t.Fatalf("Failed to lookup class name")
	}
}
