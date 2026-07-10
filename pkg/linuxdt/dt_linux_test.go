// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package linuxdt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/linuxdt"
	"github.com/jaypipes/ghw/pkg/util"
)

// chrootWith creates a fake /sys/firmware/devicetree/base tree under a temp dir
// and returns a context chrooted to it. props maps a (possibly nested, e.g.
// "chosen/u-boot,version") property name to its raw bytes.
func chrootWith(t *testing.T, props map[string][]byte) context.Context {
	t.Helper()
	root := t.TempDir()
	base := filepath.Join(root, "sys", "firmware", "devicetree", "base")
	for name, val := range props {
		path := filepath.Join(base, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, val, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return ghw.WithChroot(root)(context.TODO())
}

func TestAvailable(t *testing.T) {
	ctx := chrootWith(t, map[string][]byte{
		"model": []byte("Raspberry Pi 4 Model B Rev 1.4\x00"),
	})
	if !linuxdt.Available(ctx) {
		t.Error("expected DeviceTree to be available")
	}

	emptyCtx := ghw.WithChroot(t.TempDir())(context.TODO())
	if linuxdt.Available(emptyCtx) {
		t.Error("expected DeviceTree to be unavailable when base dir is absent")
	}
}

func TestModelAndSerial(t *testing.T) {
	// Trailing NUL must be stripped.
	ctx := chrootWith(t, map[string][]byte{
		"model":         []byte("Raspberry Pi 4 Model B Rev 1.4\x00"),
		"serial-number": []byte("10000000196f8c53\x00"),
	})
	if got, want := linuxdt.Model(ctx), "Raspberry Pi 4 Model B Rev 1.4"; got != want {
		t.Errorf("model: got %q, want %q", got, want)
	}
	if got, want := linuxdt.SerialNumber(ctx), "10000000196f8c53"; got != want {
		t.Errorf("serial: got %q, want %q", got, want)
	}

	// Absent properties return util.UNKNOWN.
	bare := chrootWith(t, map[string][]byte{"model": []byte("Some Board\x00")})
	if got := linuxdt.SerialNumber(bare); got != util.UNKNOWN {
		t.Errorf("missing serial: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestVendor(t *testing.T) {
	// Vendor is the prefix of the first NUL-separated "compatible" entry.
	ctx := chrootWith(t, map[string][]byte{
		"compatible": []byte("raspberrypi,4-model-b\x00brcm,bcm2711\x00"),
	})
	if got, want := linuxdt.Vendor(ctx), "raspberrypi"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}

	noCompat := chrootWith(t, map[string][]byte{"model": []byte("Some Board\x00")})
	if got := linuxdt.Vendor(noCompat); got != util.UNKNOWN {
		t.Errorf("vendor without compatible: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestSoC(t *testing.T) {
	// The SoC is the last (most specific) entry of the "compatible" list.
	ctx := chrootWith(t, map[string][]byte{
		"compatible": []byte("seeed,recomputer-rk3576-devkit\x00rockchip,rk3576\x00"),
	})
	if got, want := linuxdt.SoC(ctx), "rockchip,rk3576"; got != want {
		t.Errorf("soc: got %q, want %q", got, want)
	}

	noCompat := chrootWith(t, map[string][]byte{"model": []byte("Some Board\x00")})
	if got := linuxdt.SoC(noCompat); got != util.UNKNOWN {
		t.Errorf("soc without compatible: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestChassisType(t *testing.T) {
	// "embedded" maps to the SMBIOS "Embedded PC" code (34).
	ctx := chrootWith(t, map[string][]byte{
		"chassis-type": []byte("embedded\x00"),
	})
	if got, want := linuxdt.ChassisType(ctx), "34"; got != want {
		t.Errorf("chassis-type: got %q, want %q", got, want)
	}

	noType := chrootWith(t, map[string][]byte{"model": []byte("Some Board\x00")})
	if got := linuxdt.ChassisType(noType); got != util.UNKNOWN {
		t.Errorf("missing chassis-type: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestUBootVersion(t *testing.T) {
	// Lives in the nested "chosen" node, with a comma in the property name.
	ctx := chrootWith(t, map[string][]byte{
		"chosen/u-boot,version": []byte("2026.04_armbian-2026.04\x00"),
	})
	if got, want := linuxdt.UBootVersion(ctx), "2026.04_armbian-2026.04"; got != want {
		t.Errorf("u-boot version: got %q, want %q", got, want)
	}

	noUBoot := chrootWith(t, map[string][]byte{"model": []byte("Some Board\x00")})
	if got := linuxdt.UBootVersion(noUBoot); got != util.UNKNOWN {
		t.Errorf("missing u-boot version: got %q, want %q", got, util.UNKNOWN)
	}
}
