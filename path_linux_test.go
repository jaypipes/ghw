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

	// No environment variable is set for GHW_CHROOT, so pathRoot() should
	// return the default "/"
	path := pathRoot()
	if path != DEFAULT_ROOT_PATH {
		t.Fatalf("Expected pathRoot() to return '/' but got %s", path)
	}

	// Now set the GHW_CHROOT environ variable and verify that pathRoot()
	// returns that value
	os.Setenv("GHW_CHROOT", "/host")
	path = pathRoot()
	if path != "/host" {
		t.Fatalf("Expected pathRoot() to return '/host' but got %s", path)
	}
}
