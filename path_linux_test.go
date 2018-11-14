//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

// +build linux

package ghw

import (
	"os"
	"testing"
)

func TestPathRoot(t *testing.T) {
	orig, origExists := os.LookupEnv("GHW_CHROOT")
	if origExists {
		// For tests, save the original, test an override and then at the end
		// of the test, reset to the original
		defer os.Setenv("GHW_CHROOT", orig)
		os.Unsetenv("GHW_CHROOT")
	} else {
		defer os.Unsetenv("GHW_CHROOT")
	}

	ctx := contextFromEnv()

	// No environment variable is set for GHW_CHROOT, so pathProcCpuinfo() should
	// return the default "/proc/cpuinfo"
	path := ctx.pathProcCpuinfo()
	if path != "/proc/cpuinfo" {
		t.Fatalf("Expected pathProcCpuInfo() to return '/proc/cpuinfo' but got %s", path)
	}

	// Now set the GHW_CHROOT environ variable and verify that pathRoot()
	// returns that value
	os.Setenv("GHW_CHROOT", "/host")

	ctx = contextFromEnv()

	path = ctx.pathProcCpuinfo()
	if path != "/host/proc/cpuinfo" {
		t.Fatalf("Expected pathProcCpuinfo() to return '/host/proc/cpuinfo' but got %s", path)
	}
}
