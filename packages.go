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

// PackagesInfo defines installed package information
type PackagesInfo struct {
	Installed []*PackageInfo `json:"installed"`
}

// PackageInfo defines installed package information
type PackageInfo struct {
	Label       string `json:"label"`
	Version     string `json:"version"`
	InstallDate string `json:"install_date"`
}

func (info *PackagesInfo) String() string {
	var result strings.Builder
	result.WriteString("Installed Packages:")
	if info != nil && info.Installed != nil {
		for _, packageInfo := range info.Installed {
			result.WriteString(fmt.Sprintf("\t[%s]: version %s installed on: %s\n", packageInfo.Label, packageInfo.Version, packageInfo.InstallDate))
		}
		return result.String()
	}
	return "No packages found"
}

// Packages returns a pointer to a PackageInfo collection containing information
// about the host's installed packages
func Packages(opts ...*WithOption) (*PackagesInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &PackagesInfo{}
	if err := ctx.packagesFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate product information in a top-level
// "product" YAML/JSON map/object key
type packagePrinter struct {
	Info *PackagesInfo `json:"product"`
}

// YAMLString returns a string with the product information formatted as YAML
// under a top-level "dmi:" key
func (info *PackagesInfo) YAMLString() string {
	return safeYAML(packagePrinter{info})
}

// JSONString returns a string with the product information formatted as JSON
// under a top-level "product:" key
func (info *PackagesInfo) JSONString(indent bool) string {
	return safeJSON(packagePrinter{info}, indent)
}
