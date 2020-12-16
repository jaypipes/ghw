//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/jaypipes/pcidb"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/util"
)

var (
	regexAddress *regexp.Regexp = regexp.MustCompile(
		`^(([0-9a-f]{0,4}):)?([0-9a-f]{2}):([0-9a-f]{2})\.([0-9a-f]{1})$`,
	)
)

type Device struct {
	// The PCI address of the device
	Address   string         `json:"address"`
	Vendor    *pcidb.Vendor  `json:"vendor"`
	Product   *pcidb.Product `json:"product"`
	Subsystem *pcidb.Product `json:"subsystem"`
	// optional subvendor/sub-device information
	Class *pcidb.Class `json:"class"`
	// optional sub-class for the device
	Subclass *pcidb.Subclass `json:"subclass"`
	// optional programming interface
	ProgrammingInterface *pcidb.ProgrammingInterface `json:"programming_interface"`
}

// NOTE(jaypipes) Device has a custom JSON marshaller because we don't want
// to serialize the entire PCIDB information for the Vendor (which includes all
// of the vendor's products, etc). Instead, we simply serialize the ID and
// human-readable name of the vendor, product, class, etc.
func (d *Device) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString("{")
	b.WriteString(fmt.Sprintf("\"address\":\"%s\"", d.Address))
	b.WriteString(",\"vendor\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.Vendor.ID,
			d.Vendor.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"product\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.Product.ID,
			d.Product.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"subsystem\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.Subsystem.ID,
			d.Subsystem.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"class\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.Class.ID,
			d.Class.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"subclass\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.Subclass.ID,
			d.Subclass.Name,
		),
	)
	b.WriteString("},")
	b.WriteString("\"programming_interface\": {")
	b.WriteString(
		fmt.Sprintf(
			"\"id\":\"%s\",\"name\":\"%s\"",
			d.ProgrammingInterface.ID,
			d.ProgrammingInterface.Name,
		),
	)
	b.WriteString("}")
	b.WriteString("}")
	return b.Bytes(), nil
}

func (d *Device) String() string {
	vendorName := util.UNKNOWN
	if d.Vendor != nil {
		vendorName = d.Vendor.Name
	}
	productName := util.UNKNOWN
	if d.Product != nil {
		productName = d.Product.Name
	}
	className := util.UNKNOWN
	if d.Class != nil {
		className = d.Class.Name
	}
	return fmt.Sprintf(
		"%s -> class: '%s' vendor: '%s' product: '%s'",
		d.Address,
		className,
		vendorName,
		productName,
	)
}

type Info struct {
	ctx *context.Context
	// hash of class ID -> class information
	Classes map[string]*pcidb.Class
	// hash of vendor ID -> vendor information
	Vendors map[string]*pcidb.Vendor
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*pcidb.Product
}

type Address struct {
	Domain   string
	Bus      string
	Slot     string
	Function string
}

// Given a string address, returns a complete Address struct, filled in with
// domain, bus, slot and function components. The address string may either
// be in $BUS:$SLOT.$FUNCTION (BSF) format or it can be a full PCI address
// that includes the 4-digit $DOMAIN information as well:
// $DOMAIN:$BUS:$SLOT.$FUNCTION.
//
// Returns "" if the address string wasn't a valid PCI address.
func AddressFromString(address string) *Address {
	addrLowered := strings.ToLower(address)
	matches := regexAddress.FindStringSubmatch(addrLowered)
	if len(matches) == 6 {
		dom := "0000"
		if matches[1] != "" {
			dom = matches[2]
		}
		return &Address{
			Domain:   dom,
			Bus:      matches[3],
			Slot:     matches[4],
			Function: matches[5],
		}
	}
	return nil
}

// New returns a pointer to an Info struct that contains information about the
// PCI devices on the host system
func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}
	return info, nil
}
