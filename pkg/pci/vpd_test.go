//go:build linux
// +build linux

//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

// buildVPD assembles a synthetic VPD byte stream from a sequence of
// (tag, body) records terminated by an end tag. Length is encoded as
// little-endian uint16 after the tag byte.
func buildVPD(records []struct {
	tag  byte
	body []byte
}) []byte {
	var out []byte
	for _, r := range records {
		out = append(out, r.tag)
		var l [2]byte
		binary.LittleEndian.PutUint16(l[:], uint16(len(r.body)))
		out = append(out, l[:]...)
		out = append(out, r.body...)
	}
	out = append(out, vpdTagSmallEnd)
	return out
}

// buildKeywords assembles a sequence of (2-byte name, 1-byte length,
// value) tuples suitable for use as the body of a VPD-R or VPD-W
// resource.
func buildKeywords(entries []struct {
	key string
	val string
}) []byte {
	var out []byte
	for _, e := range entries {
		if len(e.key) != 2 {
			panic("vpd test: keyword must be exactly 2 chars")
		}
		out = append(out, e.key[0], e.key[1])
		out = append(out, byte(len(e.val)))
		out = append(out, []byte(e.val)...)
	}
	return out
}

func TestParseVPDEmpty(t *testing.T) {
	v, err := ParseVPD([]byte{vpdTagSmallEnd})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Identifier != "" {
		t.Errorf("expected empty identifier, got %q", v.Identifier)
	}
	if len(v.ReadOnly) != 0 {
		t.Errorf("expected empty ReadOnly map, got %v", v.ReadOnly)
	}
}

func TestParseVPDIdentifierOnly(t *testing.T) {
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{vpdTagLargeIdent, []byte("ConnectX-8 NIC")},
	})
	v, err := ParseVPD(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Identifier != "ConnectX-8 NIC" {
		t.Errorf("got identifier %q, want %q", v.Identifier, "ConnectX-8 NIC")
	}
}

func TestParseVPDReadOnlyKeywords(t *testing.T) {
	roBody := buildKeywords([]struct{ key, val string }{
		{"PN", "MCX755106AS-HEAT"},
		{"EC", "A1"},
		{"SN", "MT2317X12345"},
		{"V3", "0123456789abcdef0123456789abcdef"},
	})
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{vpdTagLargeIdent, []byte("ConnectX-7")},
		{vpdTagLargeRO, roBody},
	})
	v, err := ParseVPD(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{
		"PN": "MCX755106AS-HEAT",
		"EC": "A1",
		"SN": "MT2317X12345",
		"V3": "0123456789abcdef0123456789abcdef",
	}
	if !reflect.DeepEqual(v.ReadOnly, want) {
		t.Errorf("ReadOnly mismatch:\n got  %v\n want %v", v.ReadOnly, want)
	}
	if v.Identifier != "ConnectX-7" {
		t.Errorf("got identifier %q, want %q", v.Identifier, "ConnectX-7")
	}
}

func TestParseVPDReadWriteSection(t *testing.T) {
	rwBody := buildKeywords([]struct{ key, val string }{
		{"YA", "asset-1234"},
	})
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{vpdTagLargeRW, rwBody},
	})
	v, err := ParseVPD(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := v.ReadWrite["YA"], "asset-1234"; got != want {
		t.Errorf("YA = %q, want %q", got, want)
	}
}

func TestParseVPDUnknownTagSkipped(t *testing.T) {
	// 0x83 is unused in the VPD vocabulary we care about; the parser
	// must skip it by length without bailing out.
	roBody := buildKeywords([]struct{ key, val string }{
		{"PN", "MCX755106AS-HEAT"},
	})
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{0x83, []byte("future-tag-payload")},
		{vpdTagLargeRO, roBody},
	})
	v, err := ParseVPD(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ReadOnly["PN"] != "MCX755106AS-HEAT" {
		t.Errorf("PN not parsed past unknown tag, got %q", v.ReadOnly["PN"])
	}
}

func TestParseVPDTruncatedHeader(t *testing.T) {
	// Tag byte with only one length byte: should report truncated.
	_, err := ParseVPD([]byte{vpdTagLargeIdent, 0x10})
	if !errors.Is(err, ErrVPDTruncated) {
		t.Errorf("got err %v, want ErrVPDTruncated", err)
	}
}

func TestParseVPDTruncatedBody(t *testing.T) {
	// Declares 16 bytes of body but supplies only 4.
	_, err := ParseVPD([]byte{vpdTagLargeIdent, 0x10, 0x00, 'a', 'b', 'c', 'd'})
	if !errors.Is(err, ErrVPDTruncated) {
		t.Errorf("got err %v, want ErrVPDTruncated", err)
	}
}

func TestParseVPDTruncatedKeyword(t *testing.T) {
	// VPD-R section with a keyword header but no value bytes.
	bad := []byte{'P', 'N', 0x10} // declares 16-byte PN but nothing follows
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{vpdTagLargeRO, bad},
	})
	_, err := ParseVPD(data)
	if !errors.Is(err, ErrVPDTruncated) {
		t.Errorf("got err %v, want ErrVPDTruncated", err)
	}
}

func TestParseVPDZeroPaddingTerminates(t *testing.T) {
	// Some firmware zero-pads after the end tag. Parsing should stop
	// at a zero byte at a resource boundary without erroring.
	data := buildVPD([]struct {
		tag  byte
		body []byte
	}{
		{vpdTagLargeIdent, []byte("CX")},
	})
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	v, err := ParseVPD(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Identifier != "CX" {
		t.Errorf("got identifier %q, want %q", v.Identifier, "CX")
	}
}

func TestDeviceVPDUnavailableForParsedDevice(t *testing.T) {
	// A Device with no sysdir (e.g. one constructed via
	// Info.ParseDevice) must surface ErrVPDUnavailable rather than
	// trying to read a path-less file.
	d := &Device{}
	if _, err := d.VPD(); !errors.Is(err, ErrVPDUnavailable) {
		t.Errorf("got %v, want ErrVPDUnavailable", err)
	}
}
