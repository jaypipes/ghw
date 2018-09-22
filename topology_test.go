//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"testing"
)

// nolint: gocyclo
func TestTopology(t *testing.T) {
	info, err := Topology()

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil TopologyInfo, but got nil")
	}

	if len(info.Nodes) == 0 {
		t.Fatalf("Expected >0 nodes but got 0.")
	}

	if info.Architecture == NUMA && len(info.Nodes) == 1 {
		t.Fatalf("Got NUMA architecture but only 1 node.")
	}

	for _, n := range info.Nodes {
		if len(n.Cores) == 0 {
			t.Fatalf("Expected >0 cores but got 0.")
		}
		for _, c := range n.Cores {
			if len(c.LogicalProcessors) == 0 {
				t.Fatalf("Expected >0 logical processors but got 0.")
			}
			if uint32(len(c.LogicalProcessors)) != c.NumThreads {
				t.Fatalf(
					"Expected NumThreads == len(logical procs) but %d != %d",
					c.NumThreads,
					len(c.LogicalProcessors),
				)
			}
		}
		if len(n.Caches) == 0 {
			t.Fatalf("Expected >0 caches but got 0.")
		}
	}
}
