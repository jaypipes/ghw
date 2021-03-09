//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package address_test

import (
	"reflect"
	"strings"
	"testing"

	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
)

func TestPCIAddressFromString(t *testing.T) {

	tests := []struct {
		addrStr  string
		expected *pciaddr.Address
		// AddressFromString is more flexible than String() and wants
		// to accept addresses not in full canonical form, as long as
		// it can do the right thing - e.g. a sane default Domain exists.
		// Thus we need to sometimes skip the Address -> string check.
		skipStringTest bool
	}{
		{
			addrStr: "00:00.0",
			expected: &pciaddr.Address{
				Domain:   "0000",
				Bus:      "00",
				Slot:     "00",
				Function: "0",
			},
			skipStringTest: true,
		},
		{
			addrStr: "0000:00:00.0",
			expected: &pciaddr.Address{
				Domain:   "0000",
				Bus:      "00",
				Slot:     "00",
				Function: "0",
			},
		},
		{
			addrStr: "0000:03:00.0",
			expected: &pciaddr.Address{
				Domain:   "0000",
				Bus:      "03",
				Slot:     "00",
				Function: "0",
			},
		},
		{
			addrStr: "0000:03:00.A",
			expected: &pciaddr.Address{
				Domain:   "0000",
				Bus:      "03",
				Slot:     "00",
				Function: "a",
			},
		},
	}
	for x, test := range tests {
		got := pciaddr.FromString(test.addrStr)
		if !reflect.DeepEqual(got, test.expected) {
			t.Fatalf("Test #%d failed. Expected %v but got %v", x, test.expected, got)
		}

		if test.skipStringTest {
			continue
		}

		addrStr := got.String()
		// addresses are case insensitive
		if !strings.EqualFold(addrStr, test.addrStr) {
			t.Fatalf("Test #%d failed. Expected %q but got %q (case insensitive match)", x, test.addrStr, addrStr)
		}
	}
}
