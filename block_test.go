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

    disks := info.Disks
    if len(disks) == 0 {
        t.Errorf("Expected >0 disks. Got %d", len(disks))
    }

    d1 := disks[0]
    if d1.Name == "" {
        t.Errorf("Expected disk name, but got \"\"")
    }
    if d1.SerialNumber == "unknown" {
        t.Errorf("Got unknown serial number.")
    }
    if d1.SizeBytes <= 0 {
        t.Errorf("Expected >0 disk size, but got %d", d1.SizeBytes)
    }
}
