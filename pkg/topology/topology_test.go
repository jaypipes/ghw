//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package topology_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/topology"
)

// nolint: gocyclo
func TestTopology(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_TOPOLOGY"); ok {
		t.Skip("Skipping topology tests.")
	}

	info, err := topology.New()

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil TopologyInfo, but got nil")
	}

	if len(info.Nodes) == 0 {
		t.Fatalf("Expected >0 nodes but got 0.")
	}

	if info.Architecture == topology.ARCHITECTURE_NUMA && len(info.Nodes) == 1 {
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

func TestTopologyMarshalUnmarshal(t *testing.T) {
	data, err := topology.New(option.WithNullAlerter())
	if err != nil {
		t.Fatalf("Expected no error creating topology.Info, but got %v", err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Expected no error marshaling topology.Info, but got %v", err)
	}

	var topo *topology.Info

	err = json.Unmarshal(jdata, &topo)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling topology.Info, but got %v", err)
	}
}
