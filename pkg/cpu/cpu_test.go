//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cpu_test

import (
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/cpu"
)

// nolint: gocyclo
func TestCPU(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_CPU"); ok {
		t.Skip("Skipping CPU tests.")
	}

	info, err := cpu.New()

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
