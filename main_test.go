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
func TestHost(t *testing.T) {
	host, err := Host()

	if err != nil {
		t.Fatalf("Expected nil error but got %v", err)
	}
	if host == nil {
		t.Fatalf("Expected non-nil host but got nil.")
	}

	mem := host.Memory
	if mem == nil {
		t.Fatalf("Expected non-nil Memory but got nil.")
	}

	tpb := mem.TotalPhysicalBytes
	if tpb < 1 {
		t.Fatalf("Expected >0 total physical memory, but got %d", tpb)
	}

	tub := mem.TotalUsableBytes
	if tub < 1 {
		t.Fatalf("Expected >0 total usable memory, but got %d", tub)
	}

	cpu := host.CPU
	if cpu == nil {
		t.Fatalf("Expected non-nil CPU, but got nil")
	}

	cores := cpu.TotalCores
	if cores < 1 {
		t.Fatalf("Expected >0 total cores, but got %d", cores)
	}

	threads := cpu.TotalThreads
	if threads < 1 {
		t.Fatalf("Expected >0 total threads, but got %d", threads)
	}

	block := host.Block
	if block == nil {
		t.Fatalf("Expected non-nil Block but got nil.")
	}

	blockTpb := block.TotalPhysicalBytes
	if blockTpb < 1 {
		t.Fatalf("Expected >0 total physical block bytes, but got %d", blockTpb)
	}

	topology := host.Topology
	if topology == nil {
		t.Fatalf("Expected non-nil Topology but got nil.")
	}

	if len(topology.Nodes) < 1 {
		t.Fatalf("Expected >0 nodes , but got %d", len(topology.Nodes))
	}

	gpu := host.GPU
	if gpu == nil {
		t.Fatalf("Expected non-nil GPU but got nil.")
	}
}
