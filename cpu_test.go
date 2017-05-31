package ghw

import (
    "testing"
)

func TestCPU(t *testing.T) {
    info, err := CPU()

    if err != nil {
        t.Errorf("Expected nil err, but got %v", err)
    }
    if info == nil {
        t.Errorf("Expected non-nil CPUInfo, but got nil")
    }
}
