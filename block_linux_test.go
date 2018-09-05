//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

// +build linux

package ghw

import (
	"os"
	"reflect"
	"testing"
)

func TestParseMtabEntry(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_BLOCK"); ok {
		t.Skip("Skipping block tests.")
	}

	tests := []struct {
		line     string
		expected *mtabEntry
	}{
		{
			line: "/dev/sda6 / ext4 rw,relatime,errors=remount-ro,data=ordered 0 0",
			expected: &mtabEntry{
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
			expected: &mtabEntry{
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
			expected: &mtabEntry{
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
			expected: &mtabEntry{
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
		actual := parseMtabEntry(test.line)
		if test.expected == nil {
			if actual != nil {
				t.Fatalf("Expected nil, but got %v", actual)
			}
		} else if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d, expected %v == %v", x, test.expected, actual)
		}
	}
}
