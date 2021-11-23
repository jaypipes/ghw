//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build linux
// +build linux

package linuxpath_test

import (
	"os"

	"testing"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/option"
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

	ctx := context.FromEnv()
	paths := linuxpath.New(ctx)

	// No environment variable is set for GHW_CHROOT, so pathProcCpuinfo() should
	// return the default "/proc/cpuinfo"
	path := paths.ProcCpuinfo
	if path != "/proc/cpuinfo" {
		t.Fatalf("Expected pathProcCpuInfo() to return '/proc/cpuinfo' but got %s", path)
	}

	// Now set the GHW_CHROOT environ variable and verify that pathRoot()
	// returns that value
	os.Setenv("GHW_CHROOT", "/host")

	ctx = context.FromEnv()
	paths = linuxpath.New(ctx)

	path = paths.ProcCpuinfo
	if path != "/host/proc/cpuinfo" {
		t.Fatalf("Expected path.ProcCpuinfo to return '/host/proc/cpuinfo' but got %s", path)
	}
}

func TestPathSpecificRoots(t *testing.T) {
	ctx := context.New(option.WithPathOverrides(option.PathOverrides{
		"/proc": "/host-proc",
		"/sys":  "/host-sys",
	}))

	paths := linuxpath.New(ctx)

	path := paths.ProcCpuinfo
	expectedPath := "/host-proc/cpuinfo"
	if path != expectedPath {
		t.Fatalf("Expected path.ProcCpuInfo to return %q but got %q", expectedPath, path)
	}

	path = paths.SysBusPciDevices
	expectedPath = "/host-sys/bus/pci/devices"
	if path != expectedPath {
		t.Fatalf("Expected path.SysBusPciDevices to return %q but got %q", expectedPath, path)
	}
}

func TestPathChrootAndSpecifics(t *testing.T) {
	ctx := context.New(
		option.WithPathOverrides(option.PathOverrides{
			"/proc": "/host2-proc",
			"/sys":  "/host2-sys",
		}),
		option.WithChroot("/redirect"),
	)

	paths := linuxpath.New(ctx)

	path := paths.ProcCpuinfo
	expectedPath := "/redirect/host2-proc/cpuinfo"
	if path != expectedPath {
		t.Fatalf("Expected path.ProcCpuInfo to return %q but got %q", expectedPath, path)
	}

	path = paths.SysBusPciDevices
	expectedPath = "/redirect/host2-sys/bus/pci/devices"
	if path != expectedPath {
		t.Fatalf("Expected path.SysBusPciDevices to return %q but got %q", expectedPath, path)
	}
}
