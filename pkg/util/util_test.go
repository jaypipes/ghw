//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package util_test

import (
	"testing"

	"github.com/jaypipes/ghw/pkg/util"
)

// nolint: gocyclo
func TestConcatStrings(t *testing.T) {
	type testCase struct {
		items    []string
		expected string
	}

	testCases := []testCase{
		{
			items:    []string{},
			expected: "",
		},
		{
			items:    []string{"simple"},
			expected: "simple",
		},
		{
			items: []string{
				"foo",
				"bar",
				"baz",
			},
			expected: "foobarbaz",
		},
		{
			items: []string{
				"foo ",
				" bar ",
				" baz",
			},
			expected: "foo  bar  baz",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.expected, func(t *testing.T) {
			got := util.ConcatStrings(tCase.items...)
			if got != tCase.expected {
				t.Errorf("expected %q got %q", tCase.expected, got)
			}
		})
	}
}
