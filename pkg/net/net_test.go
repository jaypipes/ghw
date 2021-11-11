//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package net_test

import (
	"os"
	"testing"

	"github.com/adumandix/ghw/pkg/net"
)

func TestNet(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_NET"); ok {
		t.Skip("Skipping network tests.")
	}

	info, err := net.New()

	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if info == nil {
		t.Fatalf("Expected non-nil NetworkInfo, but got nil")
	}

	if len(info.NICs) > 0 {
		for _, n := range info.NICs {
			if n.Name == "" {
				t.Fatalf("Expected a NIC name but got \"\".")
			}
		}
	}
}
