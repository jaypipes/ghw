//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package block

import (
	"fmt"
	"github.com/jaypipes/ghw/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func TestParseMountEntry(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}

	tests := []struct {
		line     string
		expected *mountEntry
	}{
		{
			line: "/dev/sda6 / ext4 rw,relatime,errors=remount-ro,data=ordered 0 0",
			expected: &mountEntry{
				Partition:      "/dev/sda6",
				Mountpoint:     "/",
				FilesystemType: "ext4",
				Options: []string{
					"rw",
					"relatime",
					"errors=remount-ro",
					"data=ordered",
				},
			},
		},
		{
			line: "/dev/sda8 /home/Name\\040with\\040spaces ext4 ro 0 0",
			expected: &mountEntry{
				Partition:      "/dev/sda8",
				Mountpoint:     "/home/Name with spaces",
				FilesystemType: "ext4",
				Options: []string{
					"ro",
				},
			},
		},
		{
			// Whoever might do this in real life should be quarantined and
			// placed in administrative segregation
			line: "/dev/sda8 /home/Name\\011with\\012tab&newline ext4 ro 0 0",
			expected: &mountEntry{
				Partition:      "/dev/sda8",
				Mountpoint:     "/home/Name\twith\ntab&newline",
				FilesystemType: "ext4",
				Options: []string{
					"ro",
				},
			},
		},
		{
			line: "/dev/sda1 /home/Name\\\\withslash ext4 ro 0 0",
			expected: &mountEntry{
				Partition:      "/dev/sda1",
				Mountpoint:     "/home/Name\\withslash",
				FilesystemType: "ext4",
				Options: []string{
					"ro",
				},
			},
		},
		{
			line:     "Indy, bad dates",
			expected: nil,
		},
	}

	for x, test := range tests {
		actual := parseMountEntry(test.line)
		if test.expected == nil {
			if actual != nil {
				t.Fatalf("Expected nil, but got %v", actual)
			}
		} else if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d, expected %v == %v", x, test.expected, actual)
		}
	}
}

func TestDiskTypes(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}

	type entry struct {
		driveType         DriveType
		storageController StorageController
	}

	tests := []struct {
		line     string
		expected entry
	}{
		{
			line: "sda6",
			expected: entry{
				driveType:         DRIVE_TYPE_HDD,
				storageController: STORAGE_CONTROLLER_SCSI,
			},
		},
		{
			line: "nvme0n1",
			expected: entry{
				driveType:         DRIVE_TYPE_SSD,
				storageController: STORAGE_CONTROLLER_NVME,
			},
		},
		{
			line: "vda1",
			expected: entry{
				driveType:         DRIVE_TYPE_HDD,
				storageController: STORAGE_CONTROLLER_VIRTIO,
			},
		},
		{
			line: "xvda1",
			expected: entry{
				driveType:         DRIVE_TYPE_HDD,
				storageController: STORAGE_CONTROLLER_SCSI,
			},
		},
		{
			line: "fda1",
			expected: entry{
				driveType:         DRIVE_TYPE_FDD,
				storageController: STORAGE_CONTROLLER_UNKNOWN,
			},
		},
		{
			line: "sr0",
			expected: entry{
				driveType:         DRIVE_TYPE_ODD,
				storageController: STORAGE_CONTROLLER_SCSI,
			},
		},
		{
			line: "mmcblk0",
			expected: entry{
				driveType:         DRIVE_TYPE_SSD,
				storageController: STORAGE_CONTROLLER_MMC,
			},
		},
		{
			line: "Indy, bad dates",
			expected: entry{
				driveType:         DRIVE_TYPE_UNKNOWN,
				storageController: STORAGE_CONTROLLER_UNKNOWN,
			},
		},
	}

	for _, test := range tests {
		gotDriveType, gotStorageController := diskTypes(test.line)
		if test.expected.driveType != gotDriveType {
			t.Fatalf(
				"For %s, expected drive type %s, but got %s",
				test.line, test.expected.driveType, gotDriveType,
			)
		}
		if test.expected.storageController != gotStorageController {
			t.Fatalf(
				"For %s, expected storage controller %s, but got %s",
				test.line, test.expected.storageController, gotStorageController,
			)
		}
	}
}

func TestISCSI(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}

	baseDir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(baseDir)
	ctx := context.New()
	ctx.Chroot = baseDir
	paths := linuxpath.New(ctx)

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate an iSCSI device
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "size"), []byte("500118192\n"), 0644)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "dev"), []byte("259:0\n"), 0644)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "queue"), 0755)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "queue", "rotational"), []byte("0\n"), 0644)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "queue", "physical_block_size"), []byte("512\n"), 0644)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "device"), 0755)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "device", "vendor"), []byte("LIO-ORG\n"), 0644)
	udevData := "E:ID_MODEL=disk0\nE:ID_SERIAL=6001405961d8b6f55cf48beb0de296b2\n" +
		"E:ID_PATH=ip-192.168.130.10:3260-iscsi-iqn.2022-01.com.redhat.foo:disk0-lun-0\n" +
		"E:ID_WWN=0x6001405961d8b6f55cf48beb0de296b2\n"
	_ = ioutil.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(udevData), 0644)

	diskInventory := disks(ctx, paths)
	if diskInventory[0].DriveType != DRIVE_TYPE_ISCSI {
		t.Fatalf("Got drive type %s, but expected ISCSI", diskInventory[0].DriveType)
	}
}

func TestDiskPartLabel(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(baseDir)
	ctx := context.New()
	ctx.Chroot = baseDir
	paths := linuxpath.New(ctx)
	partLabel := "TEST_LABEL_GHW"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a disk with one partition with label TEST_LABEL_GHW
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "sda1"), 0755)
	_ = ioutil.WriteFile(filepath.Join(paths.SysBlock, "sda", "sda1", "dev"), []byte("259:0\n"), 0644)
	_ = ioutil.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_FS_LABEL=%s\n", partLabel)), 0644)
	label := diskPartLabel(paths, "sda", "sda1")
	if label != partLabel {
		t.Fatalf("Got label %s but expected %s", label, partLabel)
	}

	// Check empty label if not found
	label = diskPartLabel(paths, "sda", "sda2")
	if label != util.UNKNOWN {
		t.Fatalf("Got label %s, but expected %s label", label, util.UNKNOWN)
	}
}
