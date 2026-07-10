// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/baseboard"
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

func TestBaseboardDeviceTree(t *testing.T) {
	// Mixtile Blade 3: compatible prefix "rockchip" despite model "Mixtile Blade 3".
	root := deviceTreeChroot(t, map[string]string{
		"model":         "Mixtile Blade 3 v1.0.1\x00",
		"compatible":    "rockchip,rk3588-blade3-v101-linux\x00rockchip,rk3588\x00",
		"serial-number": "6bb18c7ac387ae18\x00",
	})

	info, err := baseboard.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Product, "Mixtile Blade 3 v1.0.1"; got != want {
		t.Errorf("product: got %q, want %q", got, want)
	}
	// Vendor is the compatible prefix, taken verbatim (a downstream-kernel quirk).
	if got, want := info.Vendor, "rockchip"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got, want := info.SerialNumber, "6bb18c7ac387ae18"; got != want {
		t.Errorf("serial: got %q, want %q", got, want)
	}
	if got := info.AssetTag; got != util.UNKNOWN {
		t.Errorf("asset tag: got %q, want %q", got, util.UNKNOWN)
	}
	if got := info.Version; got != util.UNKNOWN {
		t.Errorf("version: got %q, want %q", got, util.UNKNOWN)
	}
}
