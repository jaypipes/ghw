//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package block_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaypipes/ghw/pkg/block"

	"github.com/jaypipes/ghw/testdata"
)

// nolint: gocyclo
func TestBlock(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}

	info, err := block.New()

	if err != nil {
		t.Fatalf("Expected no error creating block.Info, but got %v", err)
	}
	tpb := info.TotalPhysicalBytes

	if tpb < 1 {
		t.Fatalf("Expected >0 total physical bytes, got %d", tpb)
	}

	disks := info.Disks
	if len(disks) == 0 {
		t.Fatalf("Expected >0 disks. Got %d", len(disks))
	}

	var d0 *block.Disk
	// Skip loop devices on generic tests as we don't know what the underlying system is going to have
	// And loop devices don't have Serial Numbers for example.
	for _, d := range disks {
		if d.StorageController != block.STORAGE_CONTROLLER_LOOP {
			d0 = d
			break
		}
	}
	if d0.Name == "" {
		t.Fatalf("Expected disk name, but got \"\"")
	}
	if d0.SizeBytes <= 0 {
		t.Fatalf("Expected >0 disk size, but got %d", d0.SizeBytes)
	}
	if d0.Partitions == nil {
		t.Fatalf("Expected non-nil partitions, but got nil.")
	}
	if d0.PhysicalBlockSizeBytes <= 0 {
		t.Fatalf("Expected >0 sector size, but got %d", d0.PhysicalBlockSizeBytes)
	}

	if len(d0.Partitions) > 0 {
		p0 := d0.Partitions[0]
		if p0 == nil {
			t.Fatalf("Expected non-nil partition, but got nil.")
		}
		if !strings.HasPrefix(p0.Name, d0.Name) {
			t.Fatalf("Expected partition name to begin with disk name but "+
				"got %s does not begin with %s", p0.Name, d0.Name)
		}
	}

	for _, p := range d0.Partitions {
		if p.SizeBytes <= 0 {
			t.Fatalf("Expected >0 partition size, but got %d", p.SizeBytes)
		}
		if p.Disk != d0 {
			t.Fatalf("Expected disk to be the same as d0 but got %v", p.Disk)
		}
	}
}

func TestBlockMarshalUnmarshal(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	blocks, err := block.New()
	if err != nil {
		t.Fatalf("Expected no error creating block.Info, but got %v", err)
	}

	data, err := json.Marshal(blocks)
	if err != nil {
		t.Fatalf("Expected no error marshaling block.Info, but got %v", err)
	}

	var bi *block.Info
	err = json.Unmarshal(data, &bi)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling block.Info, but got %v", err)
	}
}

type blockData struct {
	Block block.Info `json:"block"`
}

func TestBlockUnmarshal(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	testdataPath, err := testdata.SamplesDirectory()
	if err != nil {
		t.Fatalf("Expected nil err when detecting the samples directory, but got %v", err)
	}

	data, err := os.ReadFile(filepath.Join(testdataPath, "dell-r610-block.json"))
	if err != nil {
		t.Fatalf("Expected nil err when reading the sample data, but got %v", err)
	}

	var bd blockData
	err = json.Unmarshal(data, &bd)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling block.Info, but got %v", err)
	}

	// to learn why we check these values, please review the "dell-r610-block.json" sample
	sda := findDiskByName(bd.Block.Disks, "sda")
	if sda == nil {
		t.Fatalf("unexpected error: can't find 'sda' in the test data")
	}
	if sda.DriveType != block.DRIVE_TYPE_HDD || sda.StorageController != block.STORAGE_CONTROLLER_SCSI {
		t.Fatalf("inconsistent data for sda: %s", sda)
	}

	zram0 := findDiskByName(bd.Block.Disks, "zram0")
	if zram0 == nil {
		t.Fatalf("unexpected error: can't find 'zram0' in the test data")
	}
	if zram0.DriveType != block.DRIVE_TYPE_SSD {
		t.Fatalf("inconsistent data for zram0: %s", zram0)
	}
}

func findDiskByName(disks []*block.Disk, name string) *block.Disk {
	for _, disk := range disks {
		if disk.Name == name {
			return disk
		}
	}
	return nil
}
