// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package chassis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/chassis"
	"github.com/jaypipes/ghw/pkg/util"
)

// writeFiles writes name->content files (names may be nested subpaths) rooted at
// dir, creating parent directories as needed.
func writeFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// deviceTreeChroot builds a chroot whose DeviceTree base holds the given
// properties and which has no DMI, so the DeviceTree fallback is exercised.
func deviceTreeChroot(t *testing.T, props map[string]string) string {
	t.Helper()
	root := t.TempDir()
	writeFiles(t, filepath.Join(root, "sys", "firmware", "devicetree", "base"), props)
	return root
}

func TestChassisDeviceTreeRPi(t *testing.T) {
	root := deviceTreeChroot(t, map[string]string{
		"model":         "Raspberry Pi 4 Model B Rev 1.4\x00",
		"compatible":    "raspberrypi,4-model-b\x00brcm,bcm2711\x00",
		"serial-number": "10000000196f8c53\x00",
	})

	info, err := chassis.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Vendor, "raspberrypi"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got, want := info.Version, "Raspberry Pi 4 Model B Rev 1.4"; got != want {
		t.Errorf("version: got %q, want %q", got, want)
	}
	if got, want := info.SerialNumber, "10000000196f8c53"; got != want {
		t.Errorf("serial: got %q, want %q", got, want)
	}
	if got := info.AssetTag; got != util.UNKNOWN {
		t.Errorf("asset tag: got %q, want %q", got, util.UNKNOWN)
	}
	// No chassis-type on the Pi -> Type/TypeDescription unknown.
	if got := info.Type; got != util.UNKNOWN {
		t.Errorf("type: got %q, want %q", got, util.UNKNOWN)
	}
	if got := info.TypeDescription; got != util.UNKNOWN {
		t.Errorf("type description: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestChassisDeviceTreeChassisType(t *testing.T) {
	// radxa NIO 12L: has chassis-type, no serial-number.
	root := deviceTreeChroot(t, map[string]string{
		"model":        "Radxa NIO 12L\x00",
		"compatible":   "radxa,nio-12l\x00mediatek,mt8395\x00",
		"chassis-type": "embedded\x00",
	})

	info, err := chassis.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Vendor, "radxa"; got != want {
		t.Errorf("vendor: got %q, want %q", got, want)
	}
	if got, want := info.Type, "34"; got != want {
		t.Errorf("type: got %q, want %q", got, want)
	}
	if got, want := info.TypeDescription, "Embedded PC"; got != want {
		t.Errorf("type description: got %q, want %q", got, want)
	}
	if got := info.SerialNumber; got != util.UNKNOWN {
		t.Errorf("serial: got %q, want %q", got, util.UNKNOWN)
	}
}

func TestChassisDMITakesPrecedence(t *testing.T) {
	root := t.TempDir()
	// Both DMI and DeviceTree present: DMI must win.
	writeFiles(t, filepath.Join(root, "sys", "class", "dmi", "id"), map[string]string{
		"chassis_vendor":  "Acme Corp\n",
		"chassis_type":    "3\n",
		"chassis_serial":  "DMI-SERIAL\n",
		"chassis_version": "1.0\n",
	})
	writeFiles(t, filepath.Join(root, "sys", "firmware", "devicetree", "base"), map[string]string{
		"compatible":    "raspberrypi,4-model-b\x00",
		"serial-number": "DT-SERIAL\x00",
	})

	info, err := chassis.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Vendor, "Acme Corp"; got != want {
		t.Errorf("vendor: got %q, want %q (DMI should take precedence)", got, want)
	}
	if got, want := info.SerialNumber, "DMI-SERIAL"; got != want {
		t.Errorf("serial: got %q, want %q (DMI should take precedence)", got, want)
	}
	if got, want := info.TypeDescription, "Desktop"; got != want {
		t.Errorf("type description: got %q, want %q", got, want)
	}
}
