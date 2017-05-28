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
