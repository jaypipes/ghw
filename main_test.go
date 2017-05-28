package ghw

import (
    "testing"
)

func TestInfo(t *testing.T) {
    info, err := NewInfo()

    if err != nil {
        t.Errorf("Expected nil error but got %v", err)
    }
    if info == nil {
        t.Errorf("Expected non-nil info but got nil.")
    }

    mem := info.Memory
    if mem == nil {
        t.Errorf("Expected non-nil Memory but got nil.")
    }

    tpb := mem.TotalPhysicalBytes
    if tpb < 1 {
        t.Errorf("Expected >0 total physical memory, but got %d", tpb)
    }
}

