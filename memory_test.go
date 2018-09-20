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
func TestMemory(t *testing.T) {
	mem, err := Memory()
	if err != nil {
		t.Fatalf("Expected nil error, but got %v", err)
	}

	tpb := mem.TotalPhysicalBytes
	tub := mem.TotalUsableBytes

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
