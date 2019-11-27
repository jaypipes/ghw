// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (ctx *context) dmiFillInfo(info *DMIInfo) error {

	info.Board.AssetTag = ctx.dmiItem("board_asset_tag")
	info.Board.Serial = ctx.dmiItem("board_serial")
	info.Board.Vendor = ctx.dmiItem("board_vendor")
	info.Board.Version = ctx.dmiItem("board_version")

	info.Product.Name = ctx.dmiItem("product_name")
	info.Product.Serial = ctx.dmiItem("product_serial")
	info.Product.UUID = ctx.dmiItem("product_uuid")
	info.Product.Version = ctx.dmiItem("product_version")

	info.System.Vendor = ctx.dmiItem("sys_vendor")

	return nil
}

func (ctx *context) dmiItem(value string) string {
	path := filepath.Join(ctx.pathSysClassDMI(), "id", value)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		warn("Unable to read %s: %s\n", value, err)
		return UNKNOWN
	}

	return strings.TrimSpace(string(b))
}
