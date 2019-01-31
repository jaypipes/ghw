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

func TestParseEthtoolFeature(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_NET"); ok {
		t.Skip("Skipping network tests.")
	}

	tests := []struct {
		line     string
		expected *NICCapability
	}{
		{
			line: "scatter-gather: off",
			expected: &NICCapability{
				Name:      "scatter-gather",
				IsEnabled: false,
				CanEnable: true,
			},
		},
		{
			line: "scatter-gather: on",
			expected: &NICCapability{
				Name:      "scatter-gather",
				IsEnabled: true,
				CanEnable: true,
			},
		},
		{
			line: "scatter-gather: off [fixed]",
			expected: &NICCapability{
				Name:      "scatter-gather",
				IsEnabled: false,
				CanEnable: false,
			},
		},
	}

	for x, test := range tests {
		actual := netParseEthtoolFeature(test.line)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d, expected %v == %v", x, test.expected, actual)
		}
	}
}
