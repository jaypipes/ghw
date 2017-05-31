package ghw

import (
    "os"
    "strings"
    "testing"
)

func TestBlock(t *testing.T) {
    if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
        t.Skip("Skipping block tests.")
    }
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

    d0 := disks[0]
    if d0.Name == "" {
        t.Errorf("Expected disk name, but got \"\"")
    }
    if d0.SerialNumber == "unknown" {
        t.Errorf("Got unknown serial number.")
    }
    if d0.SizeBytes <= 0 {
        t.Errorf("Expected >0 disk size, but got %d", d0.SizeBytes)
    }
    if d0.Partitions == nil {
        t.Errorf("Expected non-nil partitions, but got nil.")
    }
    if d0.SectorSizeBytes <= 0 {
        t.Errorf("Expected >0 sector size, but got %d", d0.SectorSizeBytes)
    }

    if len(d0.Partitions) > 0 {
        p0 := d0.Partitions[0]
        if p0 == nil {
            t.Errorf("Expected non-nil partition, but got nil.")
        }
        if ! strings.HasPrefix(p0.Name, d0.Name) {
            t.Errorf("Expected partition name to begin with disk name but " +
                     "got %s does not begin with %s", p0.Name, d0.Name)
        }
    }

    for _, p := range d0.Partitions {
        // Check that all the singular functions return the same information as
        // the information constructed by ghw.Block()
        ps := PartitionSizeBytes(p.Name)
        if ps != p.SizeBytes {
            t.Errorf("Expected matching size, but got %d != %d",
                     ps, p.SizeBytes)
        }
        pmp := PartitionMountPoint(p.Name)
        if pmp != p.MountPoint {
            t.Errorf("Expected matching mountpoints, but got %s != %s",
                     pmp, p.MountPoint)
        }
        pt := PartitionType(p.Name)
        if pt != p.Type {
            t.Errorf("Expected matching types, but got %s != %s",
                     pt, p.Type)
        }
        pro := PartitionIsReadOnly(p.Name)
        if pro != p.IsReadOnly {
            t.Errorf("Expected matching readonly, but got %v != %v",
                     pro, p.IsReadOnly)
        }
        if p.Disk != d0 {
            t.Errorf("Expected disk to be the same as d0 but got %v", p.Disk)
        }
    }
}
