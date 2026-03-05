//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"os"
	"path/filepath"
)

// ExpectedCloneUSBContent returns a slice of strings pertaning to the USB interfaces
func ExpectedCloneUSBContent() []string {
	const sysBusUSB = "/sys/bus/usb/devices/"

	paths := []string{sysBusUSB}
	usbDevicesDirs, err := os.ReadDir(sysBusUSB)
	if err != nil {
		return []string{}
	}

	for _, dir := range usbDevicesDirs {
		susBusUSBLink := filepath.Join(sysBusUSB, dir.Name())
		paths = append(paths, susBusUSBLink)

		fullDir, err := os.Readlink(susBusUSBLink)
		if err != nil {
			continue
		}
		if !filepath.IsAbs(fullDir) {
			fullDir, err = filepath.Abs(filepath.Join(sysBusUSB, fullDir))
			if err != nil {
				continue
			}
		}
		for _, fileName := range []string{"uevent", "interface", "product"} {
			paths = append(paths, filepath.Join(fullDir, fileName))
		}

	}

	return paths
}
