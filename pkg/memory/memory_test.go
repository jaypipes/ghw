//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/memory"
	"github.com/jaypipes/ghw/pkg/option"
)

// nolint: gocyclo
func TestMemory(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_MEMORY"); ok {
		t.Skip("Skipping MEMORY tests.")
	}

	mem, err := memory.New()
	if err != nil {
		t.Fatalf("Expected nil error, but got %v", err)
	}

	tpb := mem.TotalPhysicalBytes
	tub := mem.TotalUsableBytes
	if tpb == 0 {
		t.Fatalf("Total physical bytes reported zero")
	}
	if tub == 0 {
		t.Fatalf("Total usable bytes reported zero")
	}
	if tpb < tub {
		t.Fatalf("Total physical bytes < total usable bytes. %d < %d",
			tpb, tub)
	}

	sps := mem.SupportedPageSizes

	if sps == nil {
		t.Fatalf("Expected non-nil supported page sizes, but got nil")
	}
	if len(sps) == 0 {
		t.Fatalf("Expected >0 supported page sizes, but got 0.")
	}
}

func TestMemoryMarshalUnmarshal(t *testing.T) {
	data, err := memory.New(option.WithNullAlerter())
	if err != nil {
		t.Fatalf("Expected no error creating memory.Info, but got %v", err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Expected no error marshaling memory.Info, but got %v", err)
	}

	var topo *memory.Info

	err = json.Unmarshal(jdata, &topo)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling memory.Info, but got %v", err)
	}
}
