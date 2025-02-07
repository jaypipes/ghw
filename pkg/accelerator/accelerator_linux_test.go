//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package accelerator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/accelerator"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/snapshot"

	"github.com/jaypipes/ghw/testdata"
)

func testScenario(t *testing.T, filename string, hwFilter []string, expectedDevs int) {
	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	t.Setenv("PCIDB_PATH", testdata.PCIDBChroot())

	workstationSnapshot := filepath.Join(testdataPath, filename)

	tmpRoot, err := os.MkdirTemp("", "ghw-accelerator-testing-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}

	_, err = snapshot.UnpackInto(workstationSnapshot, tmpRoot, 0)
	if err != nil {
		t.Fatalf("Unable to unpack %q into %q: %v", workstationSnapshot, tmpRoot, err)
	}

	defer func() {
		_ = snapshot.Cleanup(tmpRoot)
	}()

	info, err := accelerator.New(hwFilter, option.WithChroot(tmpRoot))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil AcceleratorInfo, but got nil")
	}
	if len(info.Devices) != expectedDevs {
		t.Fatalf("Expected %d processing accelerator devices, but found %d.", expectedDevs, len(info.Devices))
	}
}

func TestAcceleratorDefault(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_ACCELERATOR"); ok {
		t.Skip("Skipping PCI tests.")
	}

	// In this scenario we have 1 processing accelerator device
	testScenario(t, "linux-amd64-accel.tar.gz", []string{}, 1)

}

func TestAcceleratorNvidia(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_ACCELERATOR"); ok {
		t.Skip("Skipping PCI tests.")
	}

	// In this scenario we have 1 Nvidia 3D controller device
	testScenario(t, "linux-amd64-accel-nvidia.tar.gz", []string{}, 1)
}

func TestAcceleratorFilter(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_ACCELERATOR"); ok {
		t.Skip("Skipping PCI tests.")
	}

	// Set the filter to detect only processing accelerators (Nvidia not included)
	discoveryFilter := make([]string, 0)
	discoveryFilter = append(discoveryFilter, "::1200")

	// In this scenario we have 1 Nvidia 3D controller device
	testScenario(t, "linux-amd64-accel-nvidia.tar.gz", discoveryFilter, 0)
}
