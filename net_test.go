package ghw

import (
    "testing"
)

func TestNet(t *testing.T) {
    info, err := Net()

    if err != nil {
        t.Fatalf("Expected nil err, but got %v", err)
    }
    if info == nil {
        t.Fatalf("Expected non-nil NetInfo, but got nil")
    }

    if len(info.NICs) == 0 {
        t.Fatalf("Expected >0 NICs but got 0.")
    }

    for _, n := range info.NICs {
        if n.Name == "" {
            t.Fatalf("Expected a NIC name but got \"\".")
        }
    }
}
