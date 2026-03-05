//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package baseboard

import (
	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/util"
)

// Info defines baseboard release information
type Info struct {
	AssetTag     string `json:"asset_tag"`
	SerialNumber string `json:"serial_number"`
	Vendor       string `json:"vendor"`
	Version      string `json:"version"`
	Product      string `json:"product"`
}

func (i *Info) String() string {
	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	serialStr := ""
	if i.SerialNumber != "" && i.SerialNumber != util.UNKNOWN {
		serialStr = " serial=" + i.SerialNumber
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}

	productStr := ""
	if i.Product != "" {
		productStr = " product=" + i.Product
	}

	return "baseboard" + util.ConcatStrings(
		vendorStr,
		serialStr,
		versionStr,
		productStr,
	)
}

// New returns a pointer to an Info struct containing information about the
// host's baseboard
func New(args ...any) (*Info, error) {
	ctx := config.ContextFromArgs(args...)
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate baseboard information in a top-level
// "baseboard" YAML/JSON map/object key
type baseboardPrinter struct {
	Info *Info `json:"baseboard"`
}

// YAMLString returns a string with the baseboard information formatted as YAML
// under a top-level "dmi:" key
func (info *Info) YAMLString() string {
	return marshal.SafeYAML(baseboardPrinter{info})
}

// JSONString returns a string with the baseboard information formatted as JSON
// under a top-level "baseboard:" key
func (info *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(baseboardPrinter{info}, indent)
}
