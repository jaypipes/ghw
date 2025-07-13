//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package topology_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/memory"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/snapshot"
	"github.com/jaypipes/ghw/pkg/topology"

	"github.com/jaypipes/ghw/testdata"
)

// nolint: gocyclo
func TestTopologyNUMADistances(t *testing.T) {
	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNumaSnapshot := filepath.Join(testdataPath, "linux-amd64-intel-xeon-L5640.tar.gz")
	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(multiNumaSnapshot, unpackDir)
	if err != nil {
		t.Fatal(err)
	}
	// from now on we use constants reflecting the content of the snapshot we requested,
	// which we reviewed beforehand. IOW, you need to know the content of the
	// snapshot to fully understand this test. Inspect it using
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc topology

	info, err := topology.New(option.WithChroot(unpackDir))

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil TopologyInfo, but got nil")
	}

	if len(info.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes but got 0.")
	}

	for _, n := range info.Nodes {
		if len(n.Distances) != len(info.Nodes) {
			t.Fatalf("Expected distances to all known nodes")
		}
	}

	if info.Nodes[0].Distances[0] != info.Nodes[1].Distances[1] {
		t.Fatalf("Expected symmetric distance to self, got %v and %v", info.Nodes[0].Distances, info.Nodes[1].Distances)
	}

	if info.Nodes[0].Distances[1] != info.Nodes[1].Distances[0] {
		t.Fatalf("Expected symmetric distance to the other node, got %v and %v", info.Nodes[0].Distances, info.Nodes[1].Distances)
	}
}

// we have this test in topology_linux_test.go (and not in topology_test.go) because `topologyFillInfo`
// is not implemented on darwin; so having it in the platform-independent tests would lead to false negatives.
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

// nolint: gocyclo
func TestTopologyPerNUMAMemory(t *testing.T) {
	testdataPath, err := testdata.SnapshotsDirectory()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	multiNumaSnapshot := filepath.Join(testdataPath, "linux-amd64-intel-xeon-L5640.tar.gz")
	unpackDir := t.TempDir()
	err = snapshot.UnpackInto(multiNumaSnapshot, unpackDir)
	if err != nil {
		t.Fatal(err)
	}
	// from now on we use constants reflecting the content of the snapshot we requested,
	// which we reviewed beforehand. IOW, you need to know the content of the
	// snapshot to fully understand this test. Inspect it using
	// GHW_SNAPSHOT_PATH="/path/to/linux-amd64-intel-xeon-L5640.tar.gz" ghwc topology
	memInfo, err := memory.New(option.WithChroot(unpackDir))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if memInfo == nil {
		t.Fatalf("Expected non-nil MemoryInfo, but got nil")
	}

	info, err := topology.New(option.WithChroot(unpackDir))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil TopologyInfo, but got nil")
	}

	if len(info.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes but got 0.")
	}

	for _, node := range info.Nodes {
		if node.Memory == nil {
			t.Fatalf("missing memory information for node %d", node.ID)
		}

		if node.Memory.TotalPhysicalBytes <= 0 {
			t.Fatalf("negative physical size for node %d", node.ID)
		}
		if node.Memory.TotalPhysicalBytes > memInfo.TotalPhysicalBytes {
			t.Fatalf("physical size for node %d exceeds system's", node.ID)
		}
		if node.Memory.TotalUsableBytes <= 0 {
			t.Fatalf("negative usable size for node %d", node.ID)
		}
		if node.Memory.TotalUsableBytes > memInfo.TotalUsableBytes {
			t.Fatalf("usable size for node %d exceeds system's", node.ID)
		}
		if node.Memory.TotalUsableBytes > node.Memory.TotalPhysicalBytes {
			t.Fatalf("excessive usable size for node %d", node.ID)
		}
		if node.Memory.DefaultHugePageSize == 0 {
			t.Fatalf("unexpected default HP size for node %d", node.ID)
		}
		if len(node.Memory.HugePageAmountsBySize) != 2 {
			t.Fatalf("expected 2 huge page info records, but got '%d' for node %d", len(node.Memory.HugePageAmountsBySize), node.ID)
		}
	}
}
