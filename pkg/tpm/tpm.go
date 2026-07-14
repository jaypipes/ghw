//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package tpm

import (
	"fmt"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/marshal"
)

// Device describes a single Trusted Platform Module (TPM) present on the host
// system.
type Device struct {
	// ManufacturerName is the human-readable name of the TPM manufacturer,
	// decoded from the four-character short name the TPM reports (e.g. "AMD",
	// "IFX" for Infineon, "STM" for STMicroelectronics).
	ManufacturerName string `json:"manufacturer_name"`
	// ManufacturerVendorID is the 16-bit TCG Vendor ID assigned to the TPM
	// manufacturer by the Trusted Computing Group. Note that this is a
	// separate namespace from the PCI vendor ID. It may be empty when the
	// host only exposes the manufacturer's short name.
	// See: https://trustedcomputinggroup.org/resource/vendor-id-registry/
	ManufacturerVendorID string `json:"manufacturer_vendor_id"`
	// SpecVersion is the TCG TPM library (specification) version the device
	// implements, i.e. the ISO/IEC 11889 version, e.g. "1.2" or "2.0".
	SpecVersion string `json:"spec_version"`
	// FirmwareVersion is the TPM firmware version, when the host exposes it.
	FirmwareVersion string `json:"firmware_version"`
}

// String returns a short string describing the TPM device.
func (d *Device) String() string {
	return fmt.Sprintf(
		"manufacturer_name=%s manufacturer_vendor_id=%s firmware_version=%s spec_version=%s",
		d.ManufacturerName, d.ManufacturerVendorID, d.FirmwareVersion, d.SpecVersion,
	)
}

// Info describes the Trusted Platform Module(s) on the host system.
type Info struct {
	// Devices are the TPM devices detected under /sys/class/tpm. A host can,
	// in principle, expose more than one TPM device.
	Devices []*Device `json:"devices"`
}

// String returns a short string with summary information about the TPM(s) on
// the host system.
func (i *Info) String() string {
	switch len(i.Devices) {
	case 0:
		return "tpm (no devices)"
	case 1:
		return "tpm " + i.Devices[0].String()
	default:
		return fmt.Sprintf("tpm (%d devices)", len(i.Devices))
	}
}

// New returns a pointer to an Info struct that contains information about
// the TPM(s) on the host system.
func New(args ...any) (*Info, error) {
	ctx := config.ContextFromArgs(args...)
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate TPM information in a top-level
// "tpm" YAML/JSON map/object key
type tpmPrinter struct {
	Info *Info `json:"tpm"`
}

// YAMLString returns a string with the TPM information formatted as YAML
// under a top-level "tpm:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(tpmPrinter{i})
}

// JSONString returns a string with the TPM information formatted as JSON
// under a top-level "tpm:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(tpmPrinter{i}, indent)
}
