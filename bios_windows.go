// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "github.com/StackExchange/wmi"

type CIM_BIOSElement struct {
	InstallDate  string
	Manufacturer string
	Version      string
}

func (ctx *context) biosFillInfo(info *BIOSInfo) error {
	// Getting disk drives from WMI
	var win32BIOSDescriptions []CIM_BIOSElement
	q1 := wmi.CreateQuery(&win32BIOSDescriptions, "")
	if err := wmi.Query(q1, &win32BIOSDescriptions); err != nil {
		return err
	}
	if len(win32BIOSDescriptions) > 0 {
		info.Vendor = win32BIOSDescriptions[0].Manufacturer
		info.Version = win32BIOSDescriptions[0].Version
		info.Date = win32BIOSDescriptions[0].InstallDate
	}
	return nil
}
