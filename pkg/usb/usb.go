// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package usb

import (
	"fmt"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

type Device struct {
	Driver     string `json:"driver"`
	Type       string `json:"type"`
	VendorID   string `json:"vendor_id"`
	ProductID  string `json:"product_id"`
	Product    string `json:"product"`
	RevisionID string `json:"revision_id"`
	Interface  string `json:"interface"`
}

func (d Device) String() string {
	kvs := []struct {
		name  string
		value string
	}{
		{"driver", d.Driver},
		{"type", d.Type},
		{"vendorID", d.VendorID},
		{"productID", d.ProductID},
		{"product", d.Product},
		{"revisionID", d.RevisionID},
		{"interface", d.Interface},
	}

	var str strings.Builder

	i := 0
	for _, s := range kvs {
		k := s.name
		v := s.value

		if v == "" {
			continue
		}
		needsQuotationMarks := strings.ContainsAny(v, " \t")

		if i > 0 {
			str.WriteString(" ")
		}
		i++
		str.WriteString(k)
		str.WriteString("=")
		if needsQuotationMarks {
			str.WriteString("\"")
		}
		str.WriteString(v)
		if needsQuotationMarks {
			str.WriteString("\"")
		}

	}

	return str.String()
}

// Info describes all network interface controllers (NICs) in the host system.
type Info struct {
	ctx     *context.Context
	Devices []*Device `json:"devices"`
}

// String returns a short string with information about the networking on the
// host system.
func (i *Info) String() string {
	return fmt.Sprintf(
		"USB (%d USBs)",
		len(i.Devices),
	)
}

// New returns a pointer to an Info struct that contains information about the
// network interface controllers (NICs) on the host system
func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}

	return info, nil
}

// simple private struct used to encapsulate usb information in a
// top-level "usb" YAML/JSON map/object key
type usbPrinter struct {
	Info *Info `json:"usb"`
}

// YAMLString returns a string with the net information formatted as YAML
// under a top-level "net:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, usbPrinter{i})
}

// JSONString returns a string with the net information formatted as JSON
// under a top-level "net:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, usbPrinter{i}, indent)
}
