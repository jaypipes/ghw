// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product

import (
	"testing"

	"github.com/jaypipes/ghw/pkg/util"
)

const ioregAppleSilicon = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<array>
	<dict>
		<key>IOObjectClass</key>
		<string>IOPlatformExpertDevice</string>
		<key>IOPlatformSerialNumber</key>
		<string>C02ABC123DEF</string>
		<key>IOPlatformUUID</key>
		<string>12345678-1234-5678-9ABC-DEF012345678</string>
		<key>manufacturer</key>
		<data>QXBwbGUgSW5jLgA=</data>
		<key>model</key>
		<data>TWFjMTQsNwA=</data>
		<key>model-number</key>
		<data>WjFBQjAwMDFDAAAAAAAAAAAAAAAAAA==</data>
		<key>region-info</key>
		<data>TEwvQQAAAAAAAAAAAAAAAA==</data>
	</dict>
</array>
</plist>`

const systemProfilerOutput = `{
  "SPHardwareDataType" : [
    {
      "_name" : "hardware_overview",
      "chip_type" : "Apple M2",
      "machine_model" : "Mac14,7",
      "machine_name" : "MacBook Pro",
      "model_number" : "Z1AB0001CLL/A",
      "platform_UUID" : "12345678-1234-5678-9ABC-DEF012345678",
      "serial_number" : "C02ABC123DEF"
    }
  ]
}`

func TestParsePlatformExpertDevice(t *testing.T) {
	dev, err := parsePlatformExpertDevice([]byte(ioregAppleSilicon))
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	for key, want := range map[string]string{
		"IOPlatformSerialNumber": "C02ABC123DEF",
		"model":                  "Mac14,7",
		"manufacturer":           "Apple Inc.",
		"model-number":           "Z1AB0001C",
		"region-info":            "LL/A",
	} {
		if got := property(dev, key); got != want {
			t.Errorf("%s: got %q, want %q", key, got, want)
		}
	}
}

func TestParsePlatformExpertDeviceErrors(t *testing.T) {
	for name, out := range map[string]string{
		"empty":       "",
		"not a plist": "ioreg: not found",
		"no node":     `<?xml version="1.0"?><plist version="1.0"><array/></plist>`,
	} {
		if _, err := parsePlatformExpertDevice([]byte(out)); err == nil {
			t.Errorf("%s: expected an error, but got none", name)
		}
	}
}

func TestParseHardwareInfo(t *testing.T) {
	hw, err := parseHardwareInfo([]byte(systemProfilerOutput))
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if got, want := hw.MachineName, "MacBook Pro"; got != want {
		t.Errorf("machine_name: got %q, want %q", got, want)
	}
	if got, want := hw.MachineModel, "Mac14,7"; got != want {
		t.Errorf("machine_model: got %q, want %q", got, want)
	}
	if got, want := hw.ModelNumber, "Z1AB0001CLL/A"; got != want {
		t.Errorf("model_number: got %q, want %q", got, want)
	}
	if got, want := hw.SerialNumber, "C02ABC123DEF"; got != want {
		t.Errorf("serial_number: got %q, want %q", got, want)
	}
	if got, want := hw.PlatformUUID, "12345678-1234-5678-9ABC-DEF012345678"; got != want {
		t.Errorf("platform_UUID: got %q, want %q", got, want)
	}
}

func TestParseHardwareInfoErrors(t *testing.T) {
	for name, out := range map[string]string{
		"empty":      "",
		"not json":   "system_profiler: not found",
		"no items":   `{"SPHardwareDataType": []}`,
		"other type": `{"SPSoftwareDataType": [{"_name": "os_overview"}]}`,
	} {
		if _, err := parseHardwareInfo([]byte(out)); err == nil {
			t.Errorf("%s: expected an error, but got none", name)
		}
	}
}

func TestProperty(t *testing.T) {
	node := map[string]any{
		"IOPlatformSerialNumber": "C02ABC123DEF",
		"model":                  []byte("Mac14,7\x00"),
		"model-number":           []byte("Z1AB0001C\x00\x00\x00\x00\x00"),
		"IORegistryEntryID":      uint64(4294968043),
	}

	for _, tc := range []struct {
		key  string
		want string
	}{
		{"IOPlatformSerialNumber", "C02ABC123DEF"},
		{"model", "Mac14,7"},
		{"model-number", "Z1AB0001C"},
		{"IORegistryEntryID", ""}, // not a string property
		{"no-such-key", ""},
	} {
		if got := property(node, tc.key); got != tc.want {
			t.Errorf("property(%q): got %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestFill(t *testing.T) {
	appleSilicon, err := parsePlatformExpertDevice([]byte(ioregAppleSilicon))
	if err != nil {
		t.Fatal(err)
	}
	hw, err := parseHardwareInfo([]byte(systemProfilerOutput))
	if err != nil {
		t.Fatal(err)
	}

	// An Intel Mac carries a product version and no Apple part number, and
	// some machines only report their serial as a device tree property.
	intel := map[string]any{
		"IOPlatformUUID": "564D0C1F-9B4A-1234-ABCD-0123456789AB",
		"manufacturer":   []byte("Apple Inc.\x00"),
		"model":          []byte("MacBookPro16,1\x00"),
		"serial-number":  []byte("C02XYZ456GHI\x00\x00\x00"),
		"version":        []byte("1.0\x00"),
	}

	for _, tc := range []struct {
		name string
		dev  map[string]any
		hw   *hardwareOverview
		want Info
	}{
		{
			name: "apple silicon",
			dev:  appleSilicon,
			hw:   hw,
			want: Info{
				Family:       "MacBook Pro",
				Name:         "Mac14,7",
				Vendor:       "Apple Inc.",
				SerialNumber: "C02ABC123DEF",
				UUID:         "12345678-1234-5678-9ABC-DEF012345678",
				SKU:          "Z1AB0001CLL/A",
				Version:      util.UNKNOWN,
			},
		},
		{
			// system_profiler is what supplies the family, so without it that
			// is the only field we lose.
			name: "apple silicon without system_profiler",
			dev:  appleSilicon,
			want: Info{
				Family:       util.UNKNOWN,
				Name:         "Mac14,7",
				Vendor:       "Apple Inc.",
				SerialNumber: "C02ABC123DEF",
				UUID:         "12345678-1234-5678-9ABC-DEF012345678",
				SKU:          "Z1AB0001CLL/A",
				Version:      util.UNKNOWN,
			},
		},
		{
			name: "intel",
			dev:  intel,
			want: Info{
				Family:       util.UNKNOWN,
				Name:         "MacBookPro16,1",
				Vendor:       "Apple Inc.",
				SerialNumber: "C02XYZ456GHI",
				UUID:         "564D0C1F-9B4A-1234-ABCD-0123456789AB",
				SKU:          util.UNKNOWN,
				Version:      "1.0",
			},
		},
		{
			// Nothing in the I/O Registry: system_profiler backfills all but
			// the vendor, which it does not report.
			name: "system_profiler only",
			dev:  map[string]any{},
			hw:   hw,
			want: Info{
				Family:       "MacBook Pro",
				Name:         "Mac14,7",
				Vendor:       util.UNKNOWN,
				SerialNumber: "C02ABC123DEF",
				UUID:         "12345678-1234-5678-9ABC-DEF012345678",
				SKU:          "Z1AB0001CLL/A",
				Version:      util.UNKNOWN,
			},
		},
		{
			name: "no sources",
			dev:  map[string]any{},
			want: Info{
				Family:       util.UNKNOWN,
				Name:         util.UNKNOWN,
				Vendor:       util.UNKNOWN,
				SerialNumber: util.UNKNOWN,
				UUID:         util.UNKNOWN,
				SKU:          util.UNKNOWN,
				Version:      util.UNKNOWN,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got Info
			got.fill(tc.dev, tc.hw)
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestProductDisableTools(t *testing.T) {
	t.Setenv("GHW_DISABLE_TOOLS", "1")

	if _, err := New(); err == nil {
		t.Fatal("Expected an error with external tools disabled, but got none")
	}
}
