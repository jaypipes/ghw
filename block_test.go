package ghw

import (
    "testing"
)

func TestBlockTotalPhysicalBytes(t *testing.T) {
    info, err  := Block()
    if err != nil {
        t.Errorf("Expected no error creating BlockInfo, but got %v", err)
    }
    tpb := info.TotalPhysicalBytes

    if tpb < 1 {
        t.Errorf("Expected >0 total physical bytes, got %d", tpb)
    }
}
