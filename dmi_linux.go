// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"strings"
)

func (ctx *context) dmiFillInfo(info *DMIInfo) error {
	info.BIOS.Date = readDMI("bios_date")
	info.BIOS.Vendor = readDMI("bios_vendor")
	info.BIOS.Version = readDMI("bios_version")

	info.Board.AssetTag = readDMI("board_asset_tag")
	info.Board.Serial = readDMI("board_serial")
	info.Board.Vendor = readDMI("board_vendor")
	info.Board.Version = readDMI("board_version")

	info.Chassis.AssetTag = readDMI("chassis_asset_tag")
	info.Chassis.Serial = readDMI("chassis_serial")
	info.Chassis.Type = readDMI("chassis_type")
	info.Chassis.Vendor = readDMI("chassis_vendor")
	info.Chassis.Version = readDMI("chassis_version")

	info.Product.Name = readDMI("product_name")
	info.Product.Serial = readDMI("product_serial")
	info.Product.UUID = readDMI("product_uuid")
	info.Product.Version = readDMI("product_version")

	info.System.Vendor = readDMI("sys_vendor")

	return nil
}

func readDMI(value string) string {
	path := "/sys/class/dmi/id/" + value

	b, err := ioutil.ReadFile(path)
	if err != nil {
		warn("Unable to read " + value + " because of " + err.Error() + "\n")
		return UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
