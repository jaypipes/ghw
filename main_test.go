package ghw

import (
    "testing"
)

func TestHost(t *testing.T) {
    host, err := Host()

    if err != nil {
        t.Errorf("Expected nil error but got %v", err)
    }
    if host == nil {
        t.Errorf("Expected non-nil host but got nil.")
    }

    mem := host.Memory
    if mem == nil {
        t.Errorf("Expected non-nil Memory but got nil.")
    }

    tpb := mem.TotalPhysicalBytes
    if tpb < 1 {
        t.Errorf("Expected >0 total physical memory, but got %d", tpb)
    }

    tub := mem.TotalUsableBytes
    if tub < 1 {
        t.Errorf("Expected >0 total usable memory, but got %d", tub)
    }
}
