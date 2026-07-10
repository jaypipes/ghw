// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

// Package linuxdt reads identity information from the Linux DeviceTree, exposed
// by the kernel under /sys/firmware/devicetree/base. It is the DeviceTree
// counterpart to the linuxdmi package and is used as a fallback on systems (most
// ARM/RISC-V boards) that have no DMI/SMBIOS tables.
//
// All knowledge of DeviceTree property names and layout lives here: callers ask
// for a piece of identity (Model, Vendor, SerialNumber, ...) and never deal with
// raw property paths.
package linuxdt

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/internal/log"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

// Available returns true if the DeviceTree firmware tree is present, i.e. the
// host exposes /sys/firmware/devicetree/base.
func Available(ctx context.Context) bool {
	paths := linuxpath.New(ctx)
	fi, err := os.Stat(paths.SysFirmwareDeviceTree)
	return err == nil && fi.IsDir()
}

// Model returns the DeviceTree "model" string (e.g. "Raspberry Pi 4 Model B
// Rev 1.4"), or util.UNKNOWN if absent.
func Model(ctx context.Context) string {
	return property(ctx, "model")
}

// SerialNumber returns the DeviceTree "serial-number" string, or util.UNKNOWN if
// absent (not all boards provide one).
func SerialNumber(ctx context.Context) string {
	return property(ctx, "serial-number")
}

// Vendor returns the manufacturer/vendor prefix derived from the first entry of
// the "compatible" property (the substring before the first comma, e.g.
// "raspberrypi" from "raspberrypi,4-model-b"), or util.UNKNOWN if absent.
func Vendor(ctx context.Context) string {
	compatible := propertyList(ctx, "compatible")
	if len(compatible) == 0 {
		return util.UNKNOWN
	}
	vendor, _, _ := strings.Cut(compatible[0], ",")
	if vendor == "" {
		return util.UNKNOWN
	}
	return vendor
}

// SoC returns the System-on-Chip identifier taken from the most specific (last)
// entry of the "compatible" list, e.g. "rockchip,rk3576" from
// "seeed,recomputer-rk3576-devkit\0rockchip,rk3576\0". The compatible list runs
// from most specific board to the underlying SoC, so the last entry is the SoC.
// Returns util.UNKNOWN if there is no compatible property.
func SoC(ctx context.Context) string {
	compatible := propertyList(ctx, "compatible")
	if len(compatible) == 0 {
		return util.UNKNOWN
	}
	return compatible[len(compatible)-1]
}

// chassisTypeCodes maps the standardized DeviceTree "chassis-type" property
// values to the equivalent SMBIOS chassis type codes that ghw uses throughout,
// so callers can treat a DeviceTree-sourced chassis type identically to a
// DMI-sourced one.
var chassisTypeCodes = map[string]string{
	"desktop":     "3",  // Desktop
	"laptop":      "9",  // Laptop
	"convertible": "31", // Convertible
	"server":      "17", // Main server chassis
	"tablet":      "30", // Tablet
	"handset":     "11", // Hand held
	"watch":       "1",  // Other (no SMBIOS equivalent)
	"embedded":    "34", // Embedded PC
	"all-in-one":  "13", // All in one
}

// ChassisType returns the SMBIOS chassis type code corresponding to the
// standardized DeviceTree "chassis-type" property (e.g. "34" for "embedded"), or
// util.UNKNOWN if the property is absent or unrecognized. Only some boards
// expose it.
func ChassisType(ctx context.Context) string {
	if code, ok := chassisTypeCodes[property(ctx, "chassis-type")]; ok {
		return code
	}
	return util.UNKNOWN
}

// UBootVersion returns the U-Boot firmware version recorded by the bootloader in
// the "chosen" node, or util.UNKNOWN if absent (boards using another bootloader
// do not provide it).
func UBootVersion(ctx context.Context) string {
	return property(ctx, "chosen/u-boot,version")
}

// property reads a single string property from the DeviceTree base node and
// returns it with trailing NUL byte(s) and surrounding whitespace stripped. name
// may be a relative subpath, e.g. "chosen/u-boot,version". It returns
// util.UNKNOWN if the property is missing or cannot be read.
func property(ctx context.Context, name string) string {
	paths := linuxpath.New(ctx)
	path := filepath.Join(paths.SysFirmwareDeviceTree, name)

	log.Debug(ctx, "reading from %q", path)
	b, err := os.ReadFile(path)
	if err != nil {
		log.Warn(ctx, "Unable to read %s: %s\n", name, err)
		return util.UNKNOWN
	}

	return strings.TrimSpace(strings.Trim(string(b), "\x00"))
}

// propertyList reads a NUL-separated string list property (for example
// "compatible") from the DeviceTree base node. It returns nil if the property is
// missing or cannot be read.
func propertyList(ctx context.Context, name string) []string {
	paths := linuxpath.New(ctx)
	path := filepath.Join(paths.SysFirmwareDeviceTree, name)

	log.Debug(ctx, "reading from %q", path)
	b, err := os.ReadFile(path)
	if err != nil {
		log.Warn(ctx, "Unable to read %s: %s\n", name, err)
		return nil
	}

	var items []string
	for _, item := range strings.Split(string(b), "\x00") {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
