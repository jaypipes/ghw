//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package watchdog

import (
	"fmt"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/marshal"
)

// Info describes the hardware watchdog on the host system.
type Info struct {
	Present bool `json:"present"`
}

// String returns a short string indicating whether a hardware watchdog is
// present on the host system.
func (i *Info) String() string {
	return fmt.Sprintf("watchdog present: %v", i.Present)
}

// New returns a pointer to an Info struct that contains information about
// the hardware watchdog on the host system.
func New(args ...any) (*Info, error) {
	ctx := config.ContextFromArgs(args...)
	info := &Info{}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

// simple private struct used to encapsulate watchdog information in a
// top-level "watchdog" YAML/JSON map/object key
type watchdogPrinter struct {
	Info *Info `json:"watchdog"`
}

// YAMLString returns a string with the watchdog information formatted as YAML
// under a top-level "watchdog:" key
func (i *Info) YAMLString() string {
	return marshal.SafeYAML(watchdogPrinter{i})
}

// JSONString returns a string with the watchdog information formatted as JSON
// under a top-level "watchdog:" key
func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(watchdogPrinter{i}, indent)
}
