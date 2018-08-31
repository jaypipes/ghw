//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"os"
	"reflect"
	"testing"
)

func TestPCIAddressFromString(t *testing.T) {

	tests := []struct {
		addrStr  string
		expected *PCIAddress
	}{
		{
			addrStr: "00:00.0",
			expected: &PCIAddress{
				Domain:   "0000",
				Bus:      "00",
				Slot:     "00",
				Function: "0",
			},
		},
		{
			addrStr: "0000:00:00.0",
			expected: &PCIAddress{
				Domain:   "0000",
				Bus:      "00",
				Slot:     "00",
				Function: "0",
			},
		},
		{
			addrStr: "0000:03:00.0",
			expected: &PCIAddress{
				Domain:   "0000",
				Bus:      "03",
				Slot:     "00",
				Function: "0",
			},
		},
		{
			addrStr: "0000:03:00.A",
			expected: &PCIAddress{
				Domain:   "0000",
				Bus:      "03",
				Slot:     "00",
				Function: "a",
			},
		},
	}
	for x, test := range tests {
		got := PCIAddressFromString(test.addrStr)
		if !reflect.DeepEqual(got, test.expected) {
			t.Fatalf("Test #%d failed. Expected %v but got %v", x, test.expected, got)
		}
	}
}

func TestPCI(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := PCI()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	// Since we can't count on a specific device being present on the machine
	// being tested (and we haven't built in fixtures/mocks for things yet)
	// about all we can do is verify that the returned list of pointers to
	// PCIDevice structs is non-empty
	devs := info.ListDevices()
	if len(devs) == 0 {
		t.Fatalf("Expected to find >0 PCI devices from PCIInfo.ListDevices() but got 0.")
	}
}
