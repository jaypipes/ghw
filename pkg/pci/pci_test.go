//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"os"
	"testing"

	"github.com/jaypipes/ghw"
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
	devs := info.Devices
	if len(devs) == 0 {
		t.Fatalf("Expected to find >0 PCI devices in PCIInfo.Devices but got 0.")
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

func TestPCIWithDisableTopology(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}

	// Test with DisableTopology option enabled
	info, err := pci.New(ghw.WithDisableTopology())
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo with DisableTopology, but got %v", err)
	}

	devs := info.Devices
	if len(devs) == 0 {
		t.Fatalf("Expected to find >0 PCI devices in PCIInfo.Devices but got 0.")
	}

	// When DisableTopology is enabled, all devices should have nil Node field
	for _, dev := range devs {
		if dev.Node != nil {
			t.Errorf("Expected device %s to have nil Node when DisableTopology is enabled, but got %v", dev.Address, dev.Node)
		}
		// Verify other fields are still populated correctly
		if dev.Class == nil {
			t.Fatalf("Expected device class for %s to be non-nil", dev.Address)
		}
		if dev.Product == nil {
			t.Fatalf("Expected device product for %s to be non-nil", dev.Address)
		}
		if dev.Vendor == nil {
			t.Fatalf("Expected device vendor for %s to be non-nil", dev.Address)
		}
	}
}

func BenchmarkPCIMemoryComparison(b *testing.B) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		b.Skip("Skipping PCI benchmarks.")
	}

	b.Run("WithTopologyDetection", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := pci.New()
			if err != nil {
				b.Fatalf("Error getting PCI info: %v", err)
			}
		}
	})

	b.Run("WithDisableTopology", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := pci.New(ghw.WithDisableTopology())
			if err != nil {
				b.Fatalf("Error getting PCI info with DisableTopology: %v", err)
			}
		}
	})
}
