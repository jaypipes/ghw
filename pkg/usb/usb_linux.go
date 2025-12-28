// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func (i *Info) load(ctx context.Context) error {
	var errs []error

	i.Devices, errs = usbs(ctx)

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("error(s) happened during reading usb info: %+v", errs)
}

func fillUSBFromUevent(dir string, dev *Device) (err error) {
	ueventFp, err := os.Open(filepath.Join(dir, "uevent"))
	if err != nil {
		return
	}
	defer func() {
		err = ueventFp.Close()
	}()

	sc := bufio.NewScanner(ueventFp)
	for sc.Scan() {
		line := sc.Text()

		splits := strings.SplitN(line, "=", 2)
		if len(splits) != 2 {
			continue
		}

		key := strings.ToUpper(splits[0])
		val := splits[1]

		switch key {
		case "DRIVER":
			dev.Driver = val
		case "TYPE":
			dev.Type = val
		case "PRODUCT":
			splits := strings.SplitN(val, "/", 3)
			if len(splits) != 3 {
				continue
			}
			dev.VendorID = splits[0]
			dev.ProductID = splits[1]
			dev.RevisionID = splits[2]
		}
	}
	return nil
}

func slurp(path string) string {
	bs, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(bytes.TrimSpace(bs))
}

func usbs(ctx context.Context) ([]*Device, []error) {
	paths := linuxpath.New(ctx)
	devs := make([]*Device, 0)
	errs := []error{}

	usbDevicesDirs, err := os.ReadDir(paths.SysBusUsbDevices)
	if err != nil {
		return devs, []error{err}
	}

	for _, dir := range usbDevicesDirs {
		linkPath := filepath.Join(paths.SysBusUsbDevices, dir.Name())
		fullDir, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}
		if !filepath.IsAbs(fullDir) {
			fullDir, err = filepath.Abs(filepath.Join(paths.SysBusUsbDevices, fullDir))
			if err != nil {
				continue
			}
		}

		dev := Device{}

		err = fillUSBFromUevent(fullDir, &dev)
		if err != nil {
			errs = append(errs, err)
		}

		dev.Interface = slurp(filepath.Join(fullDir, "interface"))
		dev.Product = slurp(filepath.Join(fullDir, "product"))

		devs = append(devs, &dev)
	}

	return devs, errs
}
