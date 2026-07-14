//go:build linux
// +build linux

//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/pkg/linuxpath"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
)

// VPD is a parsed PCI Vital Product Data block. The format is defined in the
// PCI Local Bus Specification (revision 2.2 and later) appendix I.
//
// A VPD block consists of an Identifier String resource, one or more
// VPD-R (read-only) resources, optionally a VPD-W (read-write) resource,
// and an End Tag. Each VPD-R/VPD-W resource contains a sequence of
// 3-byte-keyword-prefixed fields: a two-character ASCII keyword (e.g.
// "PN", "EC", "SN", "V0"-"VZ"), a 1-byte length, and the value bytes.
//
// Keywords with well-known meanings:
//
//	PN  Part Number
//	EC  Engineering Change Level
//	FG  Fabric Geography
//	LC  Location
//	MN  Manufacturer ID
//	PG  Extended Capability
//	SN  Serial Number
//	V0-VZ  Vendor-specific
//	Y0-YZ  System-specific
//
// Vendor-specific keywords are returned in ReadOnly/ReadWrite keyed by
// the two-character keyword name (e.g. "V3").
type VPD struct {
	// Identifier is the ASCII string from the 0x82 Identifier String
	// resource. Typically a vendor-friendly product name.
	Identifier string `json:"identifier,omitempty"`
	// ReadOnly holds the keyword/value pairs from the 0x90 VPD-R section.
	// Values are stored as raw bytes (without trailing padding spaces
	// removed) so callers can decide how to interpret them.
	ReadOnly map[string]string `json:"read_only,omitempty"`
	// ReadWrite holds the keyword/value pairs from the 0x91 VPD-W
	// section, if present.
	ReadWrite map[string]string `json:"read_write,omitempty"`
}

// VPD resource tags from the PCI spec.
const (
	vpdTagSmallEnd   byte = 0x78 // small resource, end tag
	vpdTagLargeIdent byte = 0x82 // large resource, identifier string
	vpdTagLargeRO    byte = 0x90 // large resource, VPD-R (read-only)
	vpdTagLargeRW    byte = 0x91 // large resource, VPD-W (read-write)
)

// ErrVPDTruncated is returned when the input ends before a complete
// resource boundary.
var ErrVPDTruncated = errors.New("vpd: truncated input")

// ErrVPDInvalidTag is returned when an unknown or invalid tag byte
// appears at a resource boundary.
var ErrVPDInvalidTag = errors.New("vpd: invalid resource tag")

// ErrVPDUnavailable is returned by Device.VPD when the device was not
// discovered through sysfs and therefore has no associated VPD file
// to read. Devices constructed via Info.ParseDevice fall into this
// category.
var ErrVPDUnavailable = errors.New("vpd: no sysfs directory associated with device")

// ErrVPDNotPresent is returned by Device.VPD when the sysfs vpd file
// does not exist (the device does not expose VPD).
var ErrVPDNotPresent = errors.New("vpd: not present for device")

// VPD returns the parsed Vital Product Data for the device, reading it
// from the device's sysfs `vpd` file. The sysfs root is taken from ctx,
// so a ctx built via option.WithChroot reads VPD from inside the
// chroot; context.Background() reads the live /sys.
//
// The sysfs VPD file is typically root-readable only. Callers running
// without sufficient privilege will receive a wrapped permission
// error.
//
// VPD returns ErrVPDUnavailable for devices with no resolvable PCI
// address (e.g. those constructed via Info.ParseDevice) and
// ErrVPDNotPresent for devices whose sysfs entry exists but exposes
// no `vpd` file.
func (d *Device) VPD(ctx context.Context) (*VPD, error) {
	if d.Address == "" {
		return nil, ErrVPDUnavailable
	}
	pciAddr := pciaddr.FromString(d.Address)
	if pciAddr == nil {
		return nil, ErrVPDUnavailable
	}
	paths := linuxpath.New(ctx)
	vpdPath := filepath.Join(paths.SysBusPciDevices, pciAddr.String(), "vpd")
	data, err := os.ReadFile(vpdPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVPDNotPresent
		}
		return nil, fmt.Errorf("vpd: read %s: %w", vpdPath, err)
	}
	return ParseVPD(data)
}

// ParseVPD parses a raw VPD byte stream as read from
// /sys/bus/pci/devices/<addr>/vpd (or equivalent PCI config space VPD
// register window).
//
// The parser is lenient: unknown large-resource tags are skipped using
// their declared length, and parsing terminates cleanly on the End tag
// (0x78) or on a zero byte (sometimes used as padding).
func ParseVPD(data []byte) (*VPD, error) {
	v := &VPD{}
	i := 0
	for i < len(data) {
		tag := data[i]
		// Some VPD images are zero-padded after the end tag. Treat a
		// zero byte at a resource boundary as end-of-stream.
		if tag == 0 {
			break
		}
		if tag == vpdTagSmallEnd {
			break
		}
		if i+3 > len(data) {
			return nil, ErrVPDTruncated
		}
		length := int(binary.LittleEndian.Uint16(data[i+1 : i+3]))
		bodyStart := i + 3
		bodyEnd := bodyStart + length
		if bodyEnd > len(data) {
			return nil, ErrVPDTruncated
		}
		body := data[bodyStart:bodyEnd]
		switch tag {
		case vpdTagLargeIdent:
			v.Identifier = strings.TrimSpace(string(body))
		case vpdTagLargeRO:
			kv, err := parseVPDKeywords(body)
			if err != nil {
				return nil, fmt.Errorf("vpd: parsing VPD-R: %w", err)
			}
			v.ReadOnly = kv
		case vpdTagLargeRW:
			kv, err := parseVPDKeywords(body)
			if err != nil {
				return nil, fmt.Errorf("vpd: parsing VPD-W: %w", err)
			}
			v.ReadWrite = kv
		default:
			// Unknown large tag. We skip over its body without
			// erroring because the spec leaves room for future tags.
		}
		i = bodyEnd
	}
	return v, nil
}

// parseVPDKeywords walks a sequence of (2-byte keyword, 1-byte length,
// value) tuples inside a VPD-R or VPD-W resource body.
func parseVPDKeywords(body []byte) (map[string]string, error) {
	out := make(map[string]string)
	i := 0
	for i < len(body) {
		if i+3 > len(body) {
			return nil, ErrVPDTruncated
		}
		key := string(body[i : i+2])
		klen := int(body[i+2])
		valStart := i + 3
		valEnd := valStart + klen
		if valEnd > len(body) {
			return nil, ErrVPDTruncated
		}
		out[key] = string(body[valStart:valEnd])
		i = valEnd
	}
	return out, nil
}
