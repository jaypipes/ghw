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

	"github.com/jaypipes/pcidb"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
	"github.com/jaypipes/ghw/pkg/topology"
	"github.com/jaypipes/ghw/pkg/util"
)

// backward compatibility, to be removed in 1.0.0
type Address pciaddr.Address

// backward compatibility, to be removed in 1.0.0
var AddressFromString = pciaddr.FromString

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
	// Topology node that the PCI device is affined to. Will be nil if the
	// architecture is not NUMA.
	Node *topology.Node `json:"node,omitempty"`
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
	arch topology.Architecture
	ctx  *context.Context
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

// New returns a pointer to an Info struct that contains information about the
// PCI devices on the host system
func New(opts ...*option.Option) (*Info, error) {
	return NewWithContext(context.New(opts...))
}

// NewWithContext returns a pointer to an Info struct that contains information about
// the PCI devices on the host system. Use this function when you want to consume
// the topology package from another package (e.g. gpu)
func NewWithContext(ctx *context.Context) (*Info, error) {
	// by default we don't report NUMA information;
	// we will only if are sure we are running on NUMA architecture
	arch := topology.ARCHITECTURE_SMP
	topo, err := topology.NewWithContext(ctx)
	if err == nil {
		arch = topo.Architecture
	} else {
		ctx.Warn("error detecting system topology: %v", err)
	}
	info := &Info{
		arch: arch,
		ctx:  ctx,
	}
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
