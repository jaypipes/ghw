// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/bios"
	"github.com/jaypipes/ghw/pkg/util"
)

func deviceTreeChroot(t *testing.T, props map[string]string) string {
	t.Helper()
	root := t.TempDir()
	base := filepath.Join(root, "sys", "firmware", "devicetree", "base")
	for name, content := range props {
		path := filepath.Join(base, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestBIOSDeviceTreeUBoot(t *testing.T) {
	root := deviceTreeChroot(t, map[string]string{
		"compatible":            "radxa,nio-12l\x00mediatek,mt8395\x00",
		"chosen/u-boot,version": "2022.10_armbian-2022.10-S40c5\x00",
	})

	info, err := bios.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Vendor, "U-Boot"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got, want := info.Version, "2022.10_armbian-2022.10-S40c5"; got != want {
		t.Errorf("version: got %q, want %q", got, want)
	}
	if got := info.Date; got != util.UNKNOWN {
		t.Errorf("date: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestBIOSDeviceTreeNoUBoot(t *testing.T) {
	// Raspberry Pi: no chosen/u-boot,version (uses its own bootloader).
	root := deviceTreeChroot(t, map[string]string{
		"compatible": "raspberrypi,4-model-b\x00brcm,bcm2711\x00",
	})

	info, err := bios.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	for name, got := range map[string]string{
		"vendor": info.Vendor, "version": info.Version, "date": info.Date,
	} {
		if got != util.UNKNOWN {
			t.Errorf("%s: got %q, want %q", name, got, util.UNKNOWN)
		}
	}
}
