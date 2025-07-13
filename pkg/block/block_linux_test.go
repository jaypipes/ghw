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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/util"
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
		{
			line: "loop0",
			expected: entry{
				driveType:         DRIVE_TYPE_VIRTUAL,
				storageController: STORAGE_CONTROLLER_LOOP,
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

func TestDiskPartLabel(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(baseDir)
	opts := &option.Options{}
	opts.Chroot = baseDir
	paths := linuxpath.New(opts)
	partLabel := "TEST_LABEL_GHW"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a disk with one partition with label TEST_LABEL_GHW
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "sda1"), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, "sda", "sda1", "dev"), []byte("259:0\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_PART_ENTRY_NAME=%s\n", partLabel)), 0644)
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

func TestDiskFSLabel(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(baseDir)
	opts := &option.Options{}
	opts.Chroot = baseDir
	paths := linuxpath.New(opts)
	fsLabel := "TEST_LABEL_GHW"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a disk with one partition with label TEST_LABEL_GHW
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "sda1"), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, "sda", "sda1", "dev"), []byte("259:0\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_FS_LABEL=%s\n", fsLabel)), 0644)
	label := diskFSLabel(paths, "sda", "sda1")
	if label != fsLabel {
		t.Fatalf("Got label %s but expected %s", label, fsLabel)
	}

	// Check empty label if not found
	label = diskFSLabel(paths, "sda", "sda2")
	if label != util.UNKNOWN {
		t.Fatalf("Got label %s, but expected %s label", label, util.UNKNOWN)
	}
}

func TestDiskTypeUdev(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(baseDir)
	opts := &option.Options{}
	opts.Chroot = baseDir
	paths := linuxpath.New(opts)
	expectedPartType := "ext4"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a disk with one partition with label TEST_LABEL_GHW
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "sda1"), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, "sda", "sda1", "dev"), []byte("259:0\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_FS_TYPE=%s\n", expectedPartType)), 0644)
	pt := diskPartTypeUdev(paths, "sda", "sda1")
	if pt != expectedPartType {
		t.Fatalf("Got partition type %s but expected %s", pt, expectedPartType)
	}

	// Check empty fs if not found
	pt = diskPartTypeUdev(paths, "sda", "sda2")
	if pt != util.UNKNOWN {
		t.Fatalf("Got partition type %s, but expected %s", pt, util.UNKNOWN)
	}
}

func TestDiskPartUUID(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(baseDir)
	opts := &option.Options{}
	opts.Chroot = baseDir
	paths := linuxpath.New(opts)
	partUUID := "11111111-1111-1111-1111-111111111111"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a disk with one partition with uuid
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, "sda", "sda1"), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, "sda", "sda1", "dev"), []byte("259:0\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_PART_ENTRY_UUID=%s\n", partUUID)), 0644)
	uuid := diskPartUUID(paths, "sda", "sda1")
	if uuid != partUUID {
		t.Fatalf("Got uuid %s but expected %s", uuid, partUUID)
	}

	// Check empty uuid if not found
	uuid = diskPartUUID(paths, "sda", "sda2")
	if uuid != util.UNKNOWN {
		t.Fatalf("Got uuid %s, but expected %s label", uuid, util.UNKNOWN)
	}
}

// TestLoopDevicesWithOption tests to see if we find loop devices when the option is activated
func TestLoopDevicesWithOption(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}
	baseDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(baseDir)
	opts := &option.Options{}
	opts.Chroot = baseDir
	opts.DisableTools = true
	paths := linuxpath.New(opts)
	fsType := "ext4"
	expectedLoopName := "loop0"
	loopNotUsed := "loop1"
	loopPartitionName := "loop0p1"

	_ = os.MkdirAll(paths.SysBlock, 0755)
	_ = os.MkdirAll(paths.RunUdevData, 0755)

	// Emulate a loop device with one partition and another loop deviced not used
	_ = os.Mkdir(filepath.Join(paths.SysBlock, expectedLoopName), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, loopNotUsed), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, expectedLoopName, "queue"), 0755)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, loopNotUsed, "queue"), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, expectedLoopName, "queue", "rotational"), []byte("1\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, expectedLoopName, "size"), []byte("62810112\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, loopNotUsed, "size"), []byte("0\n"), 0644)
	_ = os.Mkdir(filepath.Join(paths.SysBlock, expectedLoopName, loopPartitionName), 0755)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, expectedLoopName, loopPartitionName, "dev"), []byte("259:0\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.SysBlock, expectedLoopName, loopPartitionName, "size"), []byte("102400\n"), 0644)
	_ = os.WriteFile(filepath.Join(paths.RunUdevData, "b259:0"), []byte(fmt.Sprintf("E:ID_FS_TYPE=%s\n", fsType)), 0644)
	d := disks(opts)
	// There should be one disk, the other should be ignored due to 0 size
	if len(d) != 1 {
		t.Fatalf("expected one disk device but the function reported %d", len(d))
	}
	foundDisk := d[0]
	// Should be the one we faked
	if foundDisk.Name != expectedLoopName {
		t.Fatalf("got loop device %s but expected %s", foundDisk.Name, expectedLoopName)
	}
	// Should have only one partition
	if len(foundDisk.Partitions) != 1 {
		t.Fatalf("expected one partition but the function reported %d", len(foundDisk.Partitions))
	}
	// Name should match
	if foundDisk.Partitions[0].Name != loopPartitionName {
		t.Fatalf("got partition %s but expected %s", foundDisk.Partitions[0], loopPartitionName)
	}
}
