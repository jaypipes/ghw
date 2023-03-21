//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package net

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

func TestParseEthtoolLinkInfo(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_NET"); ok {
		t.Skip("Skipping network tests.")
	}

	truePtr		:= true
	falsePtr	:= false
	tests := []struct {
		input		string
		expected	*NICLinkInfo
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
			expected: &NICLinkInfo{
				Speed:						"1000Mb/s",
				Duplex:						"Full",
				AutoNegotiation:			&truePtr,
				Port:						"Twisted Pair",
				PHYAD:						"1",
				Transceiver:				"internal",
				MDIX:						[]string{"off", "(auto)"},
				SupportsWakeOn:				"pumbg",
				WakeOn:						"d",
				LinkDetected:				&truePtr,
				SupportedPorts:				[]string{"[", "TP", "]"},
				SupportedLinkModes:			[]string{
												"10baseT/Half",
												"10baseT/Full",
												"100baseT/Half",
												"100baseT/Full",
												"1000baseT/Full",
											},
				SupportedPauseFrameUse:		&falsePtr,
				SupportsAutoNegotiation:	&truePtr,
				SupportedFECModes:			[]string{"Not", "reported"},
				AdvertisedLinkModes:		[]string{
												"10baseT/Half",
												"10baseT/Full",
												"100baseT/Half",
												"100baseT/Full",
												"1000baseT/Full",
											},
				AdvertisedPauseFrameUse:	&falsePtr,
				AdvertisedAutoNegotiation:	&truePtr,
				AdvertisedFECModes:			[]string{"Not", "reported"},
				NETIFMsgLevel:				[]string{"0x00000007", "(7)", "drv", "probe", "link"},
			},
		},
	}

	for x, test := range tests {
		actual := netParseEthtoolLinkInfo(test.input)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Fatalf("In test %d, expected %v == %v", x, test.expected, actual)
		}
	}
}
