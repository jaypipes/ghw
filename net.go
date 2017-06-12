package ghw

import (
    "fmt"
    "strings"
)

type NIC struct {
    Name string
    BusType string
    Driver string
    MacAddress string
    Model string
    Vendor string
    IsVirtual bool
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
        "NIC %s%s%s%s",
        n.Name,
        vendorStr,
        modelStr,
        isVirtualStr,
    )
}

type NetInfo struct {
    NICs []*NIC
}

func Net() (*NetInfo, error) {
    info := &NetInfo{}
    err := netFillInfo(info)
    if err != nil {
        return nil, err
    }
    return info, nil
}

func (i *NetInfo) String() string {
    return fmt.Sprintf(
        "net (%d NICs)",
        len(i.NICs),
    )
}
