//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

// NOTE(jaypipes): This entire package is deprecated and will be removed in the
// 1.0 release of ghw. Please use the aliased definitions of WithXXX functions
// in the main `ghw` package.
package option

import (
	"io"
	"log"
	"os"

	"github.com/jaypipes/pcidb"
)

const (
	DefaultChroot = "/"
)

const (
	envKeyChroot            = "GHW_CHROOT"
	envKeyDisableWarnings   = "GHW_DISABLE_WARNINGS"
	envKeyDisableTools      = "GHW_DISABLE_TOOLS"
	envKeySnapshotPath      = "GHW_SNAPSHOT_PATH"
	envKeySnapshotRoot      = "GHW_SNAPSHOT_ROOT"
	envKeySnapshotExclusive = "GHW_SNAPSHOT_EXCLUSIVE"
	envKeySnapshotPreserve  = "GHW_SNAPSHOT_PRESERVE"
)

// Alerter emits warnings about undesirable but recoverable errors.
// We use a subset of a logger interface only to emit warnings, and
// `Warninger` sounded ugly.
type Alerter interface {
	Printf(format string, v ...interface{})
}

var (
	NullAlerter = log.New(io.Discard, "", 0)
)

// EnvOrDefaultAlerter returns the default instance ghw will use to emit
// its warnings. ghw will emit warnings to stderr by default unless the
// environs variable GHW_DISABLE_WARNINGS is specified; in the latter case
// all warning will be suppressed.
func EnvOrDefaultAlerter() Alerter {
	var dest io.Writer
	if _, exists := os.LookupEnv(envKeyDisableWarnings); exists {
		dest = io.Discard
	} else {
		// default
		dest = os.Stderr
	}
	return log.New(dest, "", 0)
}

// EnvOrDefaultChroot returns the value of the GHW_CHROOT environs variable or
// the default value of "/" if not set
func EnvOrDefaultChroot() string {
	// Grab options from the environs by default
	if val, exists := os.LookupEnv(envKeyChroot); exists {
		return val
	}
	return DefaultChroot
}

// EnvOrDefaultDisableTools return true if ghw should use external tools to
// augment the data collected from sysfs. Most users want to do this most of
// time, so this is enabled by default.  Users consuming snapshots may want to
// opt out, thus they can set the GHW_DISABLE_TOOLS environs variable to any
// value to make ghw skip calling external tools even if they are available.
func EnvOrDefaultDisableTools() bool {
	if _, exists := os.LookupEnv(envKeyDisableTools); exists {
		return false
	}
	return true
}

// Options is used to represent optionally-configured settings. Each field is a
// pointer to some concrete value so that we can tell when something has been
// set or left unset.
type Options struct {
	// To facilitate querying of sysfs filesystems that are bind-mounted to a
	// non-default root mountpoint, we allow users to set the GHW_CHROOT
	// environ variable to an alternate mountpoint. For instance, assume that
	// the user of ghw is a Golang binary being executed from an application
	// container that has certain host filesystems bind-mounted into the
	// container at /host. The user would ensure the GHW_CHROOT environ
	// variable is set to "/host" and ghw will build its paths from that
	// location instead of /
	Chroot string

	// Alerter contains the target for ghw warnings
	Alerter Alerter

	// DisableTools optionally request ghw to not call any external program to learn
	// about the hardware. The default is to use such tools if available.
	DisableTools bool

	// PathOverrides optionally allows to override the default paths ghw uses
	// internally to learn about the system resources.
	PathOverrides PathOverrides

	// PCIDB allows users to provide a custom instance of the PCI database
	// (pcidb.PCIDB) to be used by ghw. This can be useful for testing,
	// supplying a preloaded database, or providing an instance created with
	// custom pcidb.WithOption settings, instead of letting ghw load the PCI
	// database automatically.
	PCIDB *pcidb.PCIDB
}

func (o *Options) Warn(msg string, args ...interface{}) {
	if o.Alerter != nil {
		o.Alerter.Printf("WARNING: "+msg, args...)
	}
}

type Option func(opts *Options)

// WithChroot allows to override the root directory ghw uses.
func WithChroot(dir string) Option {
	return func(opts *Options) {
		opts.Chroot = dir
	}
}

// WithAlerter sets alerting options for ghw
//
// DEPRECATED. Use `pkg/context.WithDisableWarnings`
func WithAlerter(alerter Alerter) Option {
	return func(opts *Options) {
		opts.Alerter = alerter
	}
}

// WithNullAlerter sets No-op alerting options for ghw
//
// DEPRECATED. Use `pkg/context.WithDisableWarnings`
func WithNullAlerter() Option {
	return func(opts *Options) {
		opts.Alerter = NullAlerter
	}
}

// WithDisableTools revents ghw from calling external tools to discover
// hardware capabilities.
//
// DEPRECATED. Use `pkg/context.WithDisableTools`
func WithDisableTools() Option {
	return func(opts *Options) {
		opts.DisableTools = true
	}
}

// WithPCIDB allows you to provide a custom instance of the PCI database (pcidb.PCIDB)
// to ghw. This is useful if you want to use a preloaded or specially configured
// PCI database, such as one created with custom pcidb.WithOption settings, instead
// of letting ghw load the PCI database automatically.
//
// DEPRECATED. Use `pkg/context.WithPCIDB`
func WithPCIDB(pcidb *pcidb.PCIDB) Option {
	return func(opts *Options) {
		opts.PCIDB = pcidb
	}
}

// PathOverrides is a map, keyed by the string name of a mount path, of override paths
type PathOverrides map[string]string

// WithPathOverrides supplies path-specific overrides for the context
//
// DEPRECATED. Use `pkg/context.WithPathOverrides`
func WithPathOverrides(overrides PathOverrides) Option {
	return func(opts *Options) {
		opts.PathOverrides = overrides
	}
}

// FromEnv returns an Options populated from the environs or default option
// values
//
// DEPRECATED. Use `pkg/context.FromEnv`
func FromEnv() *Options {
	return &Options{
		Chroot:       EnvOrDefaultChroot(),
		DisableTools: EnvOrDefaultDisableTools(),
	}
}
