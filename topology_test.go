package ghw

import (
    "testing"
)

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
}
