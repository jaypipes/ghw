// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package product_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/product"
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

func TestProductDeviceTree(t *testing.T) {
	root := deviceTreeChroot(t, map[string]string{
		"model":         "Raspberry Pi 4 Model B Rev 1.4\x00",
		"compatible":    "raspberrypi,4-model-b\x00brcm,bcm2711\x00",
		"serial-number": "10000000196f8c53\x00",
	})

	info, err := product.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Name, "Raspberry Pi 4 Model B Rev 1.4"; got != want {
		t.Errorf("name: got %q, want %q", got, want)
	}
	if got, want := info.Vendor, "raspberrypi"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got, want := info.SerialNumber, "10000000196f8c53"; got != want {
		t.Errorf("serial: got %q, want %q", got, want)
	}
	for name, got := range map[string]string{
		"family": info.Family, "uuid": info.UUID, "sku": info.SKU, "version": info.Version,
	} {
		if got != util.UNKNOWN {
			t.Errorf("%s: got %q, want %q", name, got, util.UNKNOWN)
		}
	}
}

func TestProductDeviceTreeNoSerial(t *testing.T) {
	// amlogic odroid: no serial-number property.
	root := deviceTreeChroot(t, map[string]string{
		"model":      "Hardkernel ODROID-HC4\x00",
		"compatible": "hardkernel,odroid-hc4\x00amlogic,sm1\x00",
	})

	info, err := product.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Vendor, "hardkernel"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got := info.SerialNumber; got != util.UNKNOWN {
		t.Errorf("serial: got %q, want %q", got, util.UNKNOWN)
	}
}
