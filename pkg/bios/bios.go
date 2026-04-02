//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package bios

import (
	"fmt"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/util"
)

// Info defines BIOS release information
type Info struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Date    string `json:"date"`
}

func (i *Info) String() string {

	vendorStr := ""
	if i.Vendor != "" {
		vendorStr = " vendor=" + i.Vendor
	}
	versionStr := ""
	if i.Version != "" {
		versionStr = " version=" + i.Version
	}
	dateStr := ""
	if i.Date != "" && i.Date != util.UNKNOWN {
		dateStr = " date=" + i.Date
	}

	res := fmt.Sprintf(
		"bios%s%s%s",
		vendorStr,
		versionStr,
		dateStr,
	)
	return res
}

// New returns a pointer to a Info struct containing information
// about the host's BIOS
func New(args ...any) (*Info, error) {
	ctx := config.ContextFromArgs(args...)
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}
