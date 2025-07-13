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
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/jaypipes/ghw/pkg/gpu"
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

	opts := option.FromEnv()
	paths := linuxpath.New(opts)

	// No environment variable is set for GHW_CHROOT, so pathProcCpuinfo() should
	// return the default "/proc/cpuinfo"
	path := paths.ProcCpuinfo
	if path != "/proc/cpuinfo" {
		t.Fatalf("Expected pathProcCpuInfo() to return '/proc/cpuinfo' but got %s", path)
	}

	// Now set the GHW_CHROOT environ variable and verify that pathRoot()
	// returns that value
	os.Setenv("GHW_CHROOT", "/host")

	opts = option.FromEnv()
	paths = linuxpath.New(opts)

	path = paths.ProcCpuinfo
	if path != "/host/proc/cpuinfo" {
		t.Fatalf("Expected path.ProcCpuinfo to return '/host/proc/cpuinfo' but got %s", path)
	}
}

func TestPathSpecificRoots(t *testing.T) {
	opts := option.FromEnv()
	opts.PathOverrides = option.PathOverrides{
		"/proc": "/host-proc",
		"/sys":  "/host-sys",
	}
	paths := linuxpath.New(opts)

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
	opts := option.FromEnv()
	opts.PathOverrides = option.PathOverrides{
		"/proc": "/host2-proc",
		"/sys":  "/host2-sys",
	}
	opts.Chroot = "/redirect"

	paths := linuxpath.New(opts)

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

func TestGpuPathRegexp(t *testing.T) {
	tmp := t.TempDir()

	// Make sure that last element of path is unique across other paths.
	var paths = []string{
		"../../devices/pci0000:00/0000:00:03.1/0000:07:00.0/drm/card0",
		"../../devices/pci0000:00/0000:00:0d.0/0000:01:00.1/drm/card1",
		"../../devices/pci0000:25/0000:25:01.0/0000:26:00.0/drm/card2",
		"../../devices/pci0000:89/0000:89:01.0/0000:8a:00.0/drm/card3",
	}

	// Expecting third element from paths above.
	var expectedAddrs = []string{
		"0000:07:00.0", "0000:01:00.1", "0000:26:00.0", "0000:8a:00.0",
	}

	drmPath := filepath.Join(tmp, "/sys/class/drm")
	err := os.MkdirAll(drmPath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	for _, target := range paths {
		linkname := filepath.Join(drmPath, filepath.Base(target))
		err := os.Symlink(target, linkname)
		if err != nil {
			t.Fatal(err)
		}
	}

	info, err := gpu.New(option.WithChroot(tmp))
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	if len(info.GraphicsCards) != len(expectedAddrs) {
		t.Fatalf("Expected %d graphics cards, got %d", len(expectedAddrs), len(info.GraphicsCards))
	}
	foundAddrs := make([]string, 0)
	for _, card := range info.GraphicsCards {
		foundAddrs = append(foundAddrs, card.Address)
	}

	sort.Strings(expectedAddrs)
	sort.Strings(foundAddrs)

	if !slices.Equal(expectedAddrs, foundAddrs) {
		t.Fatalf("Some cards not found")
	}
}
