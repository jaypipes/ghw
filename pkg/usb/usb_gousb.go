//go:build (!linux || gousb) && cgo
// +build !linux gousb
// +build cgo

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
)

/*
#cgo pkg-config: libusb-1.0
*/

func (i *Info) load() error {
	var err error
	i.USBs, err = usbDevices()
	return err
}

func usbDevices() ([]*USB, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true
	})

	if err != nil {
		return []*USB{}, err
	}

	ret := make([]*USB, 0, len(devs))

	for _, dev := range devs {
		usbType := fmt.Sprintf("%d/%d/%d", int(dev.Desc.Class), int(dev.Desc.SubClass), int(dev.Desc.Protocol))
		var vendorName string
		if usbid.Vendors[dev.Desc.Vendor] != nil {
			vendorName = usbid.Vendors[dev.Desc.Vendor].String()
		}
		u := USB{
			Driver:     "",
			Type:       usbType,
			VendorID:   dev.Desc.Vendor.String(),
			ProductID:  dev.Desc.Product.String(),
			Product:    vendorName,
			RevisionID: "",
			Interface:  "",
		}

		ret = append(ret, &u)
	}

	return ret, nil
}
