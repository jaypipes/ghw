//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package tpm

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

func (i *Info) load(ctx context.Context) error {
	paths := linuxpath.New(ctx)

	entries, err := os.ReadDir(paths.SysClassTPM)
	if err != nil {
		// No /sys/class/tpm directory at all means no TPM devices.
		return nil
	}

	for _, entry := range entries {
		// TPM device directories are named tpm0, tpm1, ... A system can
		// expose more than one TPM device, so we enumerate all of them
		// rather than assuming a single "tpm0".
		if !isTPMDeviceName(entry.Name()) {
			continue
		}
		dev := &Device{}
		dev.load(ctx, filepath.Join(paths.SysClassTPM, entry.Name()))
		i.Devices = append(i.Devices, dev)
	}
	return nil
}

// isTPMDeviceName reports whether name looks like a TPM device directory
// (tpm0, tpm1, ...). It deliberately excludes the resource-manager entries
// (e.g. tpmrm0) that share the tpm prefix.
func isTPMDeviceName(name string) bool {
	digits, ok := strings.CutPrefix(name, "tpm")
	if !ok || digits == "" {
		return false
	}
	_, err := strconv.Atoi(digits)
	return err == nil
}

// load populates the Device from the sysfs directory of a single TPM device
// (e.g. /sys/class/tpm/tpm0). sysfs has never been standardized in this area,
// so every source consulted here is best-effort.
func (d *Device) load(ctx context.Context, tpmPath string) {
	// Both TPM 1.2 and (some) TPM 2.0 drivers may expose a `caps` file with
	// manufacturer, spec and firmware version information. On TPM 2.0, when a
	// `caps` file is present, it does not report the spec version; that is
	// filled in separately below.
	if file, err := os.Open(filepath.Join(tpmPath, "caps")); err == nil {
		defer util.SafeClose(file)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			switch key {
			case "Manufacturer":
				d.setManufacturer(val)
			case "TCG version":
				d.SpecVersion = val
			case "Firmware version":
				d.FirmwareVersion = val
			}
		}
	}

	if d.ManufacturerName == "" && d.ManufacturerVendorID == "" {
		if b, err := os.ReadFile(filepath.Join(tpmPath, "device", "vendor")); err == nil {
			d.setManufacturer(strings.TrimSpace(string(b)))
		}
	}

	// TPM 2.0 devices expose the major spec version separately (kernel >= 5.5)
	// and generally do not report it via the `caps` file.
	if d.SpecVersion == "" {
		versionMajor := util.SafeIntFromFile(
			ctx, filepath.Join(tpmPath, "tpm_version_major"),
		)
		if versionMajor > 0 {
			d.SpecVersion = strconv.Itoa(versionMajor) + ".0"
		}
	}
}

// setManufacturer parses a raw manufacturer value read from sysfs and records
// it as either a human-readable name or a TCG Vendor ID.
//
// The kernel exposes this value inconsistently. It may be:
//   - a hex-encoded 32-bit value that is really a four-character ASCII short
//     name, e.g. "0x414D4400" -> "AMD" or "0x49465800" -> "IFX" (Infineon);
//   - a hex-encoded 16-bit value that is the TCG Vendor ID, e.g. "0x1414"
//     (Microsoft) or "0x15D1" (Infineon);
//   - a plain string name, e.g. "STM".
func (d *Device) setManufacturer(raw string) {
	if raw == "" {
		return
	}
	hexStr := strings.TrimPrefix(strings.TrimPrefix(raw, "0x"), "0X")
	if n, err := strconv.ParseUint(hexStr, 16, 64); err == nil {
		switch len(hexStr) {
		case 4: // 16 bits: the TCG Vendor ID
			d.ManufacturerVendorID = "0x" + strings.ToUpper(hexStr)
			return
		case 8: // 32 bits: a four-character ASCII short name
			if name := fourCC(uint32(n)); name != "" {
				d.ManufacturerName = name
				return
			}
		}
	}
	// Fall back to treating the value as a plain manufacturer name.
	d.ManufacturerName = raw
}

// fourCC decodes a 32-bit value as a four-character ASCII string, trimming
// trailing NUL and space padding. It returns "" if the bytes are not printable
// ASCII, in which case the value is not a short name.
func fourCC(v uint32) string {
	b := []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	for _, c := range b {
		if c != 0 && (c < 0x20 || c > 0x7e) {
			return ""
		}
	}
	return strings.TrimRight(string(b), " \x00")
}
