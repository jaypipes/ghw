//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cpu_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/snapshot"
	"github.com/jaypipes/ghw/pkg/topology"
	"github.com/jaypipes/ghw/testdata"
)

// nolint: gocyclo
func TestArmCPU(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_CPU"); ok {
		t.Skip("Skipping CPU tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNumaSnapshot := filepath.Join(testdataPath, "linux-arm64-c288e0776090cd558ef793b2a4e61939.tar.gz")

	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(multiNumaSnapshot, unpackDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := cpu.New(option.WithChroot(unpackDir))

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil CPUInfo, but got nil")
	}

	if len(info.Processors) == 0 {
		t.Fatalf("Expected >0 processors but got 0.")
	}

	for _, p := range info.Processors {
		if p.Vendor == "" {
			t.Fatalf("Expected not empty vendor field.")
		}
		if p.TotalCores == 0 {
			t.Fatalf("Expected >0 cores but got 0.")
		}
		if p.TotalHardwareThreads == 0 {
			t.Fatalf("Expected >0 threads but got 0.")
		}
		if len(p.Capabilities) == 0 {
			t.Fatalf("Expected >0 capabilities but got 0.")
		}
		if !p.HasCapability(p.Capabilities[0]) {
			t.Fatalf("Expected p to have capability %s, but did not.",
				p.Capabilities[0])
		}
		if len(p.Cores) == 0 {
			t.Fatalf("Expected >0 cores in processor, but got 0.")
		}
		for _, c := range p.Cores {
			if c.TotalHardwareThreads == 0 {
				t.Fatalf("Expected >0 threads but got 0.")
			}
			if len(c.LogicalProcessors) == 0 {
				t.Fatalf("Expected >0 logical processors but got 0.")
			}
		}
	}
}

func TestCheckCPUTopologyFilesForOfflineCPU(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_CPU"); ok {
		t.Skip("Skipping CPU tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	offlineCPUSnapshot := filepath.Join(testdataPath, "linux-amd64-offlineCPUs.tar.gz")

	// Capture stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("Cannot pipe StdErr. %v", err)
	}
	os.Stderr = wErr

	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(offlineCPUSnapshot, unpackDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := cpu.New(option.WithChroot(unpackDir))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil CPUInfo, but got nil")
	}

	if len(info.Processors) == 0 {
		t.Fatalf("Expected >0 processors but got 0.")
	}
	wErr.Close()
	var bufErr bytes.Buffer
	if _, err := io.Copy(&bufErr, rErr); err != nil {
		t.Fatalf("Failed to copy data to buffer: %v", err)
	}
	errorOutput := bufErr.String()
	if strings.Contains(errorOutput, "WARNING: failed to read int from file:") {
		t.Fatalf("Unexpected warning related to missing files under topology directory was reported")
	}
}

func TestNumCoresAmongOfflineCPUs(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_CPU"); ok {
		t.Skip("Skipping CPU tests.")
	}

	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	offlineCPUSnapshot := filepath.Join(testdataPath, "linux-amd64-offlineCPUs.tar.gz")

	// Capture stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("Cannot pipe the StdErr. %v", err)
	}
	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(offlineCPUSnapshot, unpackDir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := topology.New(option.WithChroot(unpackDir))
	if err != nil {
		t.Fatalf("Error determining node topology. %v", err)
	}

	if len(info.Nodes) < 1 {
		t.Fatal("No nodes found. Must contain one or more nodes")
	}
	for _, node := range info.Nodes {
		if len(node.Cores) < 1 {
			t.Fatal("No cores found. Must contain one or more cores")
		}
	}
	wErr.Close()
	var bufErr bytes.Buffer
	if _, err := io.Copy(&bufErr, rErr); err != nil {
		t.Fatalf("Failed to copy data to buffer: %v", err)
	}
	errorOutput := bufErr.String()
	if strings.Contains(errorOutput, "WARNING: failed to read int from file:") {
		t.Fatalf("Unexpected warnings related to missing files under topology directory was raised")
	}
}
