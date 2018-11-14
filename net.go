//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

type NICCapability struct {
	Name      string
	IsEnabled bool
	CanChange bool
}

type NIC struct {
	Name         string
	MacAddress   string
	IsVirtual    bool
	Capabilities []*NICCapability
}

func (n *NIC) String() string {
	isVirtualStr := ""
	if n.IsVirtual {
		isVirtualStr = " (virtual)"
	}
	return fmt.Sprintf(
		"%s%s",
		n.Name,
		isVirtualStr,
	)
}

type NetworkInfo struct {
	NICs []*NIC
}

func Network(opts ...*WithOption) (*NetworkInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &NetworkInfo{}
	if err := ctx.netFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *NetworkInfo) String() string {
	return fmt.Sprintf(
		"net (%d NICs)",
		len(i.NICs),
	)
}
