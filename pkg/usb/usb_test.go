//go:build linux
// +build linux

//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUSB(t *testing.T) {
	usbDir, err := os.MkdirTemp("", "TestUSB")
	if err != nil {
		t.Fatalf("could not create temp directory: %v", err)
	}

	data := `
DEVTYPE=usb_interface
DRIVER=usbhid
PRODUCT=46a/a087/101
TYPE=0/0/0
INTERFACE=3/1/2
MODALIAS=usb:v046ApA087d0101dc00dsc00dp00ic03isc01ip02in00
	`

	err = os.WriteFile(filepath.Join(usbDir, "uevent"), []byte(data), 0600)
	if err != nil {
		t.Fatalf("could not write uevent file in %s: %+v", usbDir, err)
	}

	var u USB
	err = fillUSBFromUevent(usbDir, &u)
	if err != nil {
		t.Fatalf("could not fill USB info from uevent file: %v", err)
	}

	usbExpected := USB{
		Driver:     "usbhid",
		Type:       "0/0/0",
		VendorID:   "46a",
		ProductID:  "a087",
		RevisionID: "101",
	}

	if !reflect.DeepEqual(u, usbExpected) {
		t.Fatalf("expected: %+v, but got %+v", usbExpected, u)
	}

}
