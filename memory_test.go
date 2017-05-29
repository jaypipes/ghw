package ghw

import (
    "testing"
)

func TestTotalPhysicalBytes(t *testing.T) {
    tpb := memTotalPhysicalBytes()

    if tpb < 1 {
        t.Errorf("Expected >0 total physical memory, got %d", tpb)
    }
}

func TestTotalUsableBytes(t *testing.T) {
    tpb := memTotalUsableBytes()

    if tpb < 1 {
        t.Errorf("Expected >0 total usable memory, got %d", tpb)
    }
}

func TestSupportedPageSizes(t *testing.T) {
    sps := memSupportedPageSizes()

    if sps == nil {
        t.Errorf("Expected non-nil supported page sizes, but got nil")
    }
    if len(sps) == 0 {
        t.Errorf("Expected >0 supported page sizes, but got 0.")
    }
}
