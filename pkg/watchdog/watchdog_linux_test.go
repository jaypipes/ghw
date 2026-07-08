//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package watchdog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/watchdog"
)

func TestWatchdogSysfs(t *testing.T) {
	root := t.TempDir()

	wdDir := filepath.Join(root, "sys", "class", "watchdog", "watchdog0")
	if err := os.MkdirAll(wdDir, 0700); err != nil {
		t.Fatalf("failed to create watchdog sysfs dir: %v", err)
	}

	info, err := watchdog.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if !info.Present {
		t.Error("expected watchdog to be present with a /sys/class/watchdog entry")
	}
}

func TestWatchdogDevNode(t *testing.T) {
	root := t.TempDir()

	devDir := filepath.Join(root, "dev")
	if err := os.MkdirAll(devDir, 0700); err != nil {
		t.Fatalf("failed to create dev dir: %v", err)
	}
	// A regular file stands in for the character device; detection only
	// checks for existence.
	if err := os.WriteFile(filepath.Join(devDir, "watchdog"), nil, 0600); err != nil {
		t.Fatalf("failed to create dev/watchdog: %v", err)
	}

	info, err := watchdog.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if !info.Present {
		t.Error("expected watchdog to be present with a /dev/watchdog node")
	}
}

func TestWatchdogAbsent(t *testing.T) {
	root := t.TempDir()

	// Empty watchdog class directory must not count as present.
	if err := os.MkdirAll(filepath.Join(root, "sys", "class", "watchdog"), 0700); err != nil {
		t.Fatalf("failed to create watchdog sysfs dir: %v", err)
	}

	info, err := watchdog.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if info.Present {
		t.Error("expected watchdog to be absent")
	}
}
