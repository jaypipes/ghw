//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jaypipes/pcidb"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
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
	Revision  string         `json:"revision"`
	Subsystem *pcidb.Product `json:"subsystem"`
	// optional subvendor/sub-device information
	Class *pcidb.Class `json:"class"`
	// optional sub-class for the device
	Subclass *pcidb.Subclass `json:"subclass"`
	// optional programming interface
	ProgrammingInterface *pcidb.ProgrammingInterface `json:"programming_interface"`
}

type devIdent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type devMarshallable struct {
	Address   string   `json:"address"`
	Vendor    devIdent `json:"vendor"`
	Product   devIdent `json:"product"`
	Revision  string   `json:"revision"`
	Subsystem devIdent `json:"subsystem"`
	Class     devIdent `json:"class"`
	Subclass  devIdent `json:"subclass"`
	Interface devIdent `json:"programming_interface"`
}

// NOTE(jaypipes) Device has a custom JSON marshaller because we don't want
// to serialize the entire PCIDB information for the Vendor (which includes all
// of the vendor's products, etc). Instead, we simply serialize the ID and
// human-readable name of the vendor, product, class, etc.
func (d *Device) MarshalJSON() ([]byte, error) {
	dm := devMarshallable{
		Address: d.Address,
		Vendor: devIdent{
			ID:   d.Vendor.ID,
			Name: d.Vendor.Name,
		},
		Product: devIdent{
			ID:   d.Product.ID,
			Name: d.Product.Name,
		},
		Revision: d.Revision,
		Subsystem: devIdent{
			ID:   d.Subsystem.ID,
			Name: d.Subsystem.Name,
		},
		Class: devIdent{
			ID:   d.Class.ID,
			Name: d.Class.Name,
		},
		Subclass: devIdent{
			ID:   d.Subclass.ID,
			Name: d.Subclass.Name,
		},
		Interface: devIdent{
			ID:   d.ProgrammingInterface.ID,
			Name: d.ProgrammingInterface.Name,
		},
	}
	return json.Marshal(dm)
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
	// All PCI devices on the host system
	Devices []*Device
	// hash of class ID -> class information
	// DEPRECATED. Will be removed in v1.0. Please use
	// github.com/jaypipes/pcidb to explore PCIDB information
	Classes map[string]*pcidb.Class `json:"-"`
	// hash of vendor ID -> vendor information
	// DEPRECATED. Will be removed in v1.0. Please use
	// github.com/jaypipes/pcidb to explore PCIDB information
	Vendors map[string]*pcidb.Vendor `json:"-"`
	// hash of vendor ID + product/device ID -> product information
	// DEPRECATED. Will be removed in v1.0. Please use
	// github.com/jaypipes/pcidb to explore PCIDB information
	Products map[string]*pcidb.Product `json:"-"`
}

func (i *Info) String() string {
	return fmt.Sprintf("PCI (%d devices)", len(i.Devices))
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

// simple private struct used to encapsulate PCI information in a top-level
// "pci" YAML/JSON map/object key
type pciPrinter struct {
	Info *Info `json:"pci"`
}

// YAMLString returns a string with the PCI information formatted as YAML
// under a top-level "pci:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, pciPrinter{i})
}

// JSONString returns a string with the PCI information formatted as JSON
// under a top-level "pci:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, pciPrinter{i}, indent)
}
