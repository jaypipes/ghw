//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package tpm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/tpm"
)

// makeTPMDir creates <root>/sys/class/tpm/<name> and returns its path.
func makeTPMDir(t *testing.T, root, name string) string {
	t.Helper()
	tpmDir := filepath.Join(root, "sys", "class", "tpm", name)
	if err := os.MkdirAll(tpmDir, 0700); err != nil {
		t.Fatalf("failed to create tpm sysfs dir: %v", err)
	}
	return tpmDir
}

func TestTPM12Caps(t *testing.T) {
	root := t.TempDir()
	tpmDir := makeTPMDir(t, root, "tpm0")

	// 0x49465800 is the hex-encoded ASCII short name "IFX" (Infineon).
	caps := `Manufacturer: 0x49465800
TCG version: 1.2
Firmware version: 6.40
`
	if err := os.WriteFile(filepath.Join(tpmDir, "caps"), []byte(caps), 0600); err != nil {
		t.Fatalf("failed to write caps: %v", err)
	}

	info, err := tpm.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if len(info.Devices) != 1 {
		t.Fatalf("expected 1 TPM device, but got %d", len(info.Devices))
	}
	dev := info.Devices[0]
	if dev.ManufacturerName != "IFX" {
		t.Errorf("unexpected manufacturer name: %q", dev.ManufacturerName)
	}
	if dev.SpecVersion != "1.2" {
		t.Errorf("unexpected spec version: %q", dev.SpecVersion)
	}
	if dev.FirmwareVersion != "6.40" {
		t.Errorf("unexpected firmware version: %q", dev.FirmwareVersion)
	}
}

func TestTPM20VersionMajor(t *testing.T) {
	root := t.TempDir()
	tpmDir := makeTPMDir(t, root, "tpm0")

	// TPM 2.0 exposes no caps file; the manufacturer comes from the device
	// vendor attribute and the spec version from tpm_version_major.
	deviceDir := filepath.Join(tpmDir, "device")
	if err := os.MkdirAll(deviceDir, 0700); err != nil {
		t.Fatalf("failed to create device dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(deviceDir, "vendor"), []byte("STM\n"), 0600); err != nil {
		t.Fatalf("failed to write vendor: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tpmDir, "tpm_version_major"), []byte("2\n"), 0600); err != nil {
		t.Fatalf("failed to write tpm_version_major: %v", err)
	}

	info, err := tpm.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if len(info.Devices) != 1 {
		t.Fatalf("expected 1 TPM device, but got %d", len(info.Devices))
	}
	dev := info.Devices[0]
	if dev.ManufacturerName != "STM" {
		t.Errorf("unexpected manufacturer name: %q", dev.ManufacturerName)
	}
	if dev.SpecVersion != "2.0" {
		t.Errorf("unexpected spec version: %q", dev.SpecVersion)
	}
}

func TestTPM20VendorID(t *testing.T) {
	root := t.TempDir()
	tpmDir := makeTPMDir(t, root, "tpm0")

	// A 16-bit hex-encoded vendor attribute is the TCG Vendor ID (0x1414 is
	// Microsoft), not a manufacturer short name.
	deviceDir := filepath.Join(tpmDir, "device")
	if err := os.MkdirAll(deviceDir, 0700); err != nil {
		t.Fatalf("failed to create device dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(deviceDir, "vendor"), []byte("0x1414\n"), 0600); err != nil {
		t.Fatalf("failed to write vendor: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tpmDir, "tpm_version_major"), []byte("2\n"), 0600); err != nil {
		t.Fatalf("failed to write tpm_version_major: %v", err)
	}

	info, err := tpm.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if len(info.Devices) != 1 {
		t.Fatalf("expected 1 TPM device, but got %d", len(info.Devices))
	}
	dev := info.Devices[0]
	if dev.ManufacturerVendorID != "0x1414" {
		t.Errorf("unexpected manufacturer vendor id: %q", dev.ManufacturerVendorID)
	}
	if dev.ManufacturerName != "" {
		t.Errorf("expected empty manufacturer name, but got %q", dev.ManufacturerName)
	}
}

func TestMultipleTPMDevices(t *testing.T) {
	root := t.TempDir()
	// A tpm0 with a version file plus a bare tpm1: both should be enumerated.
	tpm0 := makeTPMDir(t, root, "tpm0")
	if err := os.WriteFile(filepath.Join(tpm0, "tpm_version_major"), []byte("2\n"), 0600); err != nil {
		t.Fatalf("failed to write tpm_version_major: %v", err)
	}
	makeTPMDir(t, root, "tpm1")
	// A resource-manager entry sharing the tpm prefix must be ignored.
	makeTPMDir(t, root, "tpmrm0")

	info, err := tpm.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if len(info.Devices) != 2 {
		t.Fatalf("expected 2 TPM devices, but got %d", len(info.Devices))
	}
}

func TestTPMAbsent(t *testing.T) {
	root := t.TempDir()

	// An empty tpm class directory without any tpmN entry means no TPM.
	if err := os.MkdirAll(filepath.Join(root, "sys", "class", "tpm"), 0700); err != nil {
		t.Fatalf("failed to create tpm sysfs dir: %v", err)
	}

	info, err := tpm.New(ghw.WithChroot(root))
	if err != nil {
		t.Fatalf("expected nil err, but got %v", err)
	}
	if len(info.Devices) != 0 {
		t.Errorf("expected no TPM devices, but got %d", len(info.Devices))
	}
}
