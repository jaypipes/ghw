package ghw

import (
    "testing"
)

func TestMemTotalPhysicalBytes(t *testing.T) {
    tpb := memTotalPhysicalBytes()

    if tpb < 1 {
        t.Fatalf("Expected >0 total physical memory, got %d", tpb)
    }
}

func TestMemTotalUsableBytes(t *testing.T) {
    tpb := memTotalUsableBytes()

    if tpb < 1 {
        t.Fatalf("Expected >0 total usable memory, got %d", tpb)
    }
}

func TestMemSupportedPageSizes(t *testing.T) {
    sps := memSupportedPageSizes()

    if sps == nil {
        t.Fatalf("Expected non-nil supported page sizes, but got nil")
    }
    if len(sps) == 0 {
        t.Fatalf("Expected >0 supported page sizes, but got 0.")
    }
}
