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
			line: "/dev/sda8 /home/Name\040with\040spaces ext4 ro 0 0",
			expected: &mtabEntry{
				Partition:      "/dev/sda6",
				Mountpoint:     "/Name with spaces",
				FilesystemType: "ext4",
				Options: []string{
					"ro",
				},
			},
		},
	}

	for x, test := range tests {
		actual := parseMtabEntry(test.line)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d, expected %v == %v", x, test.expected, actual)
		}
	}
}
