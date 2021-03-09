//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/pci"
)

func TestPCI(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := pci.New()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	// Since we can't count on a specific device being present on the machine
	// being tested (and we haven't built in fixtures/mocks for things yet)
	// about all we can do is verify that the returned list of pointers to
	// PCIDevice structs is non-empty
	devs := info.ListDevices()
	if len(devs) == 0 {
		t.Fatalf("Expected to find >0 PCI devices from PCIInfo.ListDevices() but got 0.")
	}

	// Ensure that the data fields are at least populated, even if we don't yet
	// check for data accuracy
	for _, dev := range devs {
		if dev.Class == nil {
			t.Fatalf("Expected device class for %s to be non-nil", dev.Address)
		}
		if dev.Product == nil {
			t.Fatalf("Expected device product for %s to be non-nil", dev.Address)
		}
		if dev.Vendor == nil {
			t.Fatalf("Expected device vendor for %s to be non-nil", dev.Address)
		}
		if dev.Revision == "" {
			t.Fatalf("Expected device revision for %s to be non-empty", dev.Address)
		}
		if dev.Subclass == nil {
			t.Fatalf("Expected device subclass for %s to be non-nil", dev.Address)
		}
		if dev.Subsystem == nil {
			t.Fatalf("Expected device subsystem for %s to be non-nil", dev.Address)
		}
		if dev.ProgrammingInterface == nil {
			t.Fatalf("Expected device programming interface for %s to be non-nil", dev.Address)
		}
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
