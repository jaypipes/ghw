//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/pci"

	"github.com/jaypipes/ghw/testdata"
)

type pciTestCase struct {
	addr     string
	node     int
	revision string
	driver   string
}

// nolint: gocyclo
func TestPCINUMANode(t *testing.T) {
	info := pciTestSetup(t)

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
	info := pciTestSetup(t)

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
		t.Run(fmt.Sprintf("%s", tCase.addr), func(t *testing.T) {
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
func TestPCIDriver(t *testing.T) {
	info := pciTestSetup(t)

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

	dev := info.ParseDevice("0000:3c:00.0", "pci:v0000144Dd0000A804sv0000144Dsd0000A801bc01sc08i02")
	s := marshal.SafeJSON(context.FromEnv(), dev, true)
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
		t.Fatalf("Parsed succesfully junk data")
	}

	dev = info.ParseDevice("0000:00:01.0", "pci:v00008086d00005916sv000017AAsd0000224Bbc03sc00i00extrajunkextradataextraextra")
	if dev == nil {
		t.Fatalf("Failed to parse valid modalias with extra data")
	}
}

func pciTestSetup(t *testing.T) *pci.Info {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNumaSnapshot := filepath.Join(testdataPath, "linux-amd64-intel-xeon-L5640.tar.gz")
	// from now on we use constants reflecting the content of the snapshot we requested,
	// which we reviewed beforehand. IOW, you need to know the content of the
	// snapshot to fully understand this test. Inspect it using
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc topology

	info, err := pci.New(option.WithSnapshot(option.SnapshotOptions{
		Path: multiNumaSnapshot,
	}))

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil PCIInfo, but got nil")
	}
	return info
}
