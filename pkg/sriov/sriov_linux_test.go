//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sriov_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaypipes/ghw/pkg/option"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
	"github.com/jaypipes/ghw/pkg/sriov"

	"github.com/jaypipes/ghw/testdata"
)

// nolint: gocyclo
func TestStringify(t *testing.T) {
	info := sriovTestSetup(t)

	for _, physFn := range info.PhysicalFunctions {
		s := physFn.String()
		if s == "" || !strings.Contains(s, "function") || !strings.Contains(s, "physical") {
			t.Errorf("Wrong string representation %q", s)
		}
	}

	for _, virtFn := range info.VirtualFunctions {
		s := virtFn.String()
		if s == "" || !strings.Contains(s, "function") || !strings.Contains(s, "virtual") {
			t.Errorf("Wrong string representation %q", s)
		}
	}

}

// nolint: gocyclo
func TestCountDevices(t *testing.T) {
	info := sriovTestSetup(t)

	// Check the content of
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc sriov
	// to verify these magic numbers
	expectedPhysDevs := 2
	expectedVirtDevsPerPhys := 4
	numPhysDevs := len(info.PhysicalFunctions)
	if numPhysDevs != expectedPhysDevs {
		t.Errorf("Expected %d physical devices found %d", expectedPhysDevs, numPhysDevs)
	}
	numVirtDevs := len(info.VirtualFunctions)
	if numPhysDevs*expectedVirtDevsPerPhys != numVirtDevs {
		t.Errorf("Expected %d=(%d*%d) virtual devices found %d", numPhysDevs*expectedVirtDevsPerPhys, numPhysDevs, expectedVirtDevsPerPhys, numVirtDevs)
	}

	for _, physDev := range info.PhysicalFunctions {
		numVFs := len(physDev.VFs)
		if numVFs != expectedVirtDevsPerPhys {
			t.Errorf("Expected %d virtual devices for PF %s found %d", expectedVirtDevsPerPhys, physDev.Address.String(), numVFs)
		}
	}
}

type pfTestCase struct {
	addr    string
	netname string
}

// nolint: gocyclo
func TestMatchPhysicalFunction(t *testing.T) {
	info := sriovTestSetup(t)

	// Check the content of
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc sriov
	// to verify these magic numbers
	for _, pfTC := range []pfTestCase{
		{
			addr:    "0000:05:00.0",
			netname: "enp5s0f0",
		},
		{
			addr:    "0000:05:00.1",
			netname: "enp5s0f1",
		},
	} {
		addr := pciaddr.FromString(pfTC.addr)
		pf := findPF(info.PhysicalFunctions, addr)
		if pf == nil {
			t.Fatalf("missing PF at addr %q", addr.String())
		}
		if pf.PCI == nil {
			t.Errorf("missing PCI device for %q", addr.String())
		}
		if pf.PCI.Driver != "igb" {
			t.Errorf("unexpected driver for %#v: %q", pf, pf.PCI.Driver)
		}
		if len(pf.Interfaces) != 1 || pf.Interfaces[0] != pfTC.netname {
			t.Errorf("unexpected interfaces for %#v: %v", pf, pf.Interfaces)
		}
		if pf.MaxVFNum != 7 {
			t.Errorf("unexpected MaxVFNum for %#v: %d", pf, pf.MaxVFNum)
		}
		if len(pf.VFs) != 4 {
			t.Errorf("unexpected VF count for %#v: %d", pf, len(pf.VFs))
		}
		for _, vfInst := range pf.VFs {
			vf := findVF(info.VirtualFunctions, vfInst.Address)
			if vf == nil {
				t.Errorf("VF %#v from %#v not found among info.VirtualFunctions", vfInst, pf)
			}
		}
	}
}

func TestMatchVirtualFunction(t *testing.T) {
	info := sriovTestSetup(t)

	// Check the content of
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc sriov
	// to verify these magic numbers

	for _, vf := range info.VirtualFunctions {
		if vf.PCI == nil {
			t.Errorf("missing PCI device for %q", vf.Address.String())
		}
		if vf.PCI.Driver != "igbvf" {
			t.Errorf("unexpected driver for %#v: %q", vf, vf.PCI.Driver)
		}

		pf := findPF(info.PhysicalFunctions, vf.ParentAddress)
		if pf == nil {
			t.Fatalf("missing parent device for %q", vf.Address.String())
		}
		if vf2 := findVFInst(pf.VFs, vf.Address); vf2 == nil {
			t.Errorf("VF %#v not included in parent %#v VFs", vf, pf)
		}
	}
}

func findPF(pfs []*sriov.PhysicalFunction, addr *pciaddr.Address) *sriov.PhysicalFunction {
	for _, pf := range pfs {
		if pf.Address.Equal(addr) {
			return pf
		}
	}
	return nil
}

func findVF(vfs []*sriov.VirtualFunction, addr *pciaddr.Address) *sriov.VirtualFunction {
	for _, vf := range vfs {
		if vf.Address.Equal(addr) {
			return vf
		}
	}
	return nil
}

func findVFInst(vfs []sriov.VirtualFunction, addr *pciaddr.Address) *sriov.VirtualFunction {
	for idx := 0; idx < len(vfs); idx++ {
		if vfs[idx].Address.Equal(addr) {
			return &vfs[idx]
		}
	}
	return nil
}

func sriovTestSetup(t *testing.T) *sriov.Info {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_SRIOV"); ok {
		t.Skip("Skipping SRIOV tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNumaSnapshot := filepath.Join(testdataPath, "linux-amd64-intel-xeon-L5640.tar.gz")
	// from now on we use constants reflecting the content of the snapshot we requested,
	// which we reviewed beforehand. IOW, you need to know the content of the
	// snapshot to fully understand this test. Inspect it using
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc sriov

	info, err := sriov.New(option.WithSnapshot(option.SnapshotOptions{
		Path: multiNumaSnapshot,
	}))

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil SRIOVInfo, but got nil")
	}
	return info
}
