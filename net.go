//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"strings"
)

type NIC struct {
	Name            string
	BusType         string
	Driver          string
	MacAddress      string
	Model           string
	Vendor          string
	IsVirtual       bool
	EnabledFeatures []string
}

func (n *NIC) String() string {
	vendorStr := ""
	if n.Vendor != "" {
		vendorStr = " [" + strings.TrimSpace(n.Vendor) + "]"
	}
	modelStr := ""
	if n.Model != "" {
		modelStr = " - " + strings.TrimSpace(n.Model)
	}
	isVirtualStr := ""
	if n.IsVirtual {
		isVirtualStr = " (virtual)"
	}
	return fmt.Sprintf(
		"%s%s%s%s",
		n.Name,
		vendorStr,
		modelStr,
		isVirtualStr,
	)
}

type NetworkInfo struct {
	NICs []*NIC
}

func Network() (*NetworkInfo, error) {
	info := &NetworkInfo{}
	err := netFillInfo(info)
	if err != nil {
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
