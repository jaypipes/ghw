//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package net

import (
	"bytes"
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

func TestParseNicAttrEthtool(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_NET"); ok {
		t.Skip("Skipping network tests.")
	}

	tests := []struct {
		input    string
		expected *NIC
	}{
		{
			input: `Settings for eth0:
	Supported ports: [ TP ]
	Supported link modes:   10baseT/Half 10baseT/Full
	                        100baseT/Half 100baseT/Full
	                        1000baseT/Full
	Supported pause frame use: No
	Supports auto-negotiation: Yes
	Supported FEC modes: Not reported
	Advertised link modes:  10baseT/Half 10baseT/Full
	                        100baseT/Half 100baseT/Full
	                        1000baseT/Full
	Advertised pause frame use: No
	Advertised auto-negotiation: Yes
	Advertised FEC modes: Not reported
	Speed: 1000Mb/s
	Duplex: Full
	Auto-negotiation: on
	Port: Twisted Pair
	PHYAD: 1
	Transceiver: internal
	MDI-X: off (auto)
	Supports Wake-on: pumbg
	Wake-on: d
        Current message level: 0x00000007 (7)
                               drv probe link
	Link detected: yes
`,
			expected: &NIC{
				Speed:          "1000Mb/s",
				Duplex:         "Full",
				SupportedPorts: []string{"TP"},
				AdvertisedLinkModes: []string{
					"10baseT/Half",
					"10baseT/Full",
					"100baseT/Half",
					"100baseT/Full",
					"1000baseT/Full",
				},
				SupportedLinkModes: []string{
					"10baseT/Half",
					"10baseT/Full",
					"100baseT/Half",
					"100baseT/Full",
					"1000baseT/Full",
				},
				Capabilities: []*NICCapability{
					{
						Name:      "auto-negotiation",
						IsEnabled: true,
						CanEnable: true,
					},
					{
						Name:      "pause-frame-use",
						IsEnabled: false,
						CanEnable: false,
					},
				},
			},
		},
	}

	for x, test := range tests {
		m := parseNicAttrEthtool(bytes.NewBufferString(test.input))
		actual := &NIC{}
		actual.Speed = strings.Join(m["Speed"], "")
		actual.Duplex = strings.Join(m["Duplex"], "")
		actual.SupportedLinkModes = m["Supported link modes"]
		actual.SupportedPorts = m["Supported ports"]
		actual.SupportedFECModes = m["Supported FEC modes"]
		actual.AdvertisedLinkModes = m["Advertised link modes"]
		actual.AdvertisedFECModes = m["Advertised FEC modes"]
		actual.Capabilities = append(actual.Capabilities, autoNegCap(m))
		actual.Capabilities = append(actual.Capabilities, pauseFrameUseCap(m))
		if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d\nExpected:\n%+v\nActual:\n%+v\n", x, test.expected, actual)
		}
	}
}
