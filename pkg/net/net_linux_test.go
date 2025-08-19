//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package net

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/safchain/ethtool"
)

func TestAlwaysFailUntilWeFigureOutTheOutputDifferences(t *testing.T) {
	t.Fail()
}

func TestNetDeviceCapabilitiesFromEthHandle(t *testing.T) {
	type testCase struct {
		name string
		dev  string
		// dev(string) -> map[string]FeatureState
		feats map[string]map[string]ethtool.FeatureState
		// dev(string) -> error
		errs          map[string]error
		expected      []*NICCapability
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "nil data",
			dev:           "foodev",
			expectedError: fmt.Errorf("unsupported device: foodev"),
		},
		{
			name: "empty data",
			dev:  "foodev",
			feats: map[string]map[string]ethtool.FeatureState{
				"foodev": map[string]ethtool.FeatureState{},
			},
			expected: []*NICCapability{},
		},
		{
			name: "minimal data",
			dev:  "foodev",
			feats: map[string]map[string]ethtool.FeatureState{
				"foodev": map[string]ethtool.FeatureState{
					"foo": ethtool.FeatureState{
						Available:    true,
						NeverChanged: true,
					},
					"bar": ethtool.FeatureState{
						Available: true,
						Requested: true,
					},
				},
			},
			expected: []*NICCapability{
				{
					Name:      "bar",
					CanEnable: true,
				},
				{
					Name:      "foo",
					CanEnable: true,
				},
			},
		},
		{
			name: "minimal data and errors",
			dev:  "foodev",
			feats: map[string]map[string]ethtool.FeatureState{
				"foodev": map[string]ethtool.FeatureState{
					"foo": ethtool.FeatureState{
						Available:    true,
						NeverChanged: true,
					},
					"bar": ethtool.FeatureState{
						Available: true,
						Requested: true,
					},
				},
			},
			errs: map[string]error{
				"foodev": fmt.Errorf("fake error"),
			},
			expected:      nil,
			expectedError: fmt.Errorf("fake error"),
		},
		{
			name:  "only errors",
			dev:   "foodev",
			feats: map[string]map[string]ethtool.FeatureState{},
			errs: map[string]error{
				"foodev": fmt.Errorf("fake error"),
			},
			expected:      nil,
			expectedError: fmt.Errorf("fake error"),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			fc := fakeCollector{
				feats: tCase.feats,
				errs:  tCase.errs,
			}
			got, err := netDeviceCapabilitiesFromEthHandle(fc, tCase.dev)
			if (err != nil) != (tCase.expectedError != nil) {
				t.Fatalf("got error %v expected error %v", err, tCase.expectedError)
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].Name < got[j].Name
			})
			if !reflect.DeepEqual(got, tCase.expected) {
				t.Errorf("got %v expected %v", got, tCase.expected)
			}
		})
	}

}

type fakeCollector struct {
	// dev(string) -> map[string]FeatureState
	feats map[string]map[string]ethtool.FeatureState
	errs  map[string]error
}

func (fc fakeCollector) FeaturesWithState(intf string) (map[string]ethtool.FeatureState, error) {
	if err, ok := fc.errs[intf]; ok {
		return nil, err
	}
	if feats, ok := fc.feats[intf]; ok {
		return feats, nil
	}
	return nil, fmt.Errorf("unsupported device: %s", intf)
}
