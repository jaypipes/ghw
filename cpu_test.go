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

    if len(info.Processors) == 0 {
        t.Errorf("Expected >0 processors but got 0.")
    }

    for _, p := range info.Processors {
        if p.NumCores == 0 {
            t.Errorf("Expected >0 cores but got 0.")
        }
        if p.NumThreads == 0 {
            t.Errorf("Expected >0 threads but got 0.")
        }
        if len(p.Capabilities) == 0 {
            t.Errorf("Expected >0 capabilities but got 0.")
        }
    }
}
