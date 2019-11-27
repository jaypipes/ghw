//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

// DMIInfo describes all the information for the hardware
type DMIInfo struct {
	BIOS    BIOSInfo    `json:"bios_info"`
	Board   BoardInfo   `json:"board"`
	Chassis ChassisInfo `json:"chassis"`
	Product ProductInfo `json:"product"`
	System  SystemInfo  `json:"system"`
}

// BIOSInfo defines BIOS release information
type BIOSInfo struct {
	Date    string `json:"date"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

// BoardInfo defines motherboard release information
type BoardInfo struct {
	AssetTag string `json:"asset_tag"`
	Serial   string `json:"serial"`
	Vendor   string `json:"vendor"`
	Version  string `json:"version"`
}

// ProductInfo defines product information
type ProductInfo struct {
	Name    string `json:"name"`
	Serial  string `json:"serial"`
	UUID    string `json:"uuid"`
	Version string `json:"version"`
}

// SystemInfo defines system information
type SystemInfo struct {
	Vendor string `json:"vendor"`
}

func (info *DMIInfo) String() string {
	return fmt.Sprintf(
		"dmi\n  bios: %+v\n  board: %+v\n  chassis: %+v\n  product: %+v\n  system: %+v",
		info.BIOS,
		info.Board,
		info.Chassis,
		info.Product,
		info.System,
	)
}

// func (productInfo *ProductInfo) String() string {
// 	return fmt.Sprintf(
// 		"%s serial='%s' uuid:'%s' version:'%s'",
// 		productInfo.Name,
// 		productInfo.Serial,
// 		productInfo.UUID,
// 		productInfo.Version,
// 	)
// }

// func (systemInfo *SystemInfo) String() string {
// 	return fmt.Sprintf(
// 		"vendor='%s'",
// 		systemInfo.Vendor,
// 	)
// }

// DMI provides motherboard, chassis and BIOS information
func DMI(opts ...*WithOption) (*DMIInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &DMIInfo{}
	if err := ctx.dmiFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate DMI information in a top-level
// "dmi" YAML/JSON map/object key
type dmiPrinter struct {
	Info *DMIInfo `json:"dmi"`
}

// YAMLString returns a string with the DMI information formatted as YAML
// under a top-level "dmi:" key
func (info *DMIInfo) YAMLString() string {
	return safeYAML(dmiPrinter{info})
}

// JSONString returns a string with the DMI information formatted as JSON
// under a top-level "dmi:" key
func (info *DMIInfo) JSONString(indent bool) string {
	return safeJSON(dmiPrinter{info}, indent)
}
