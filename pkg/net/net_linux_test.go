//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package net

import "testing"

func TestEthtool(t *testing.T) {
	// intentionally always fail to make CI red so we don't merge by mistake
	t.Fail()
}
