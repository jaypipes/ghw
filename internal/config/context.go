//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package config

import (
	"context"
	"os"

	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/pcidb"
)

const (
	envKeyChroot          = "GHW_CHROOT"
	envKeyDisableWarnings = "GHW_DISABLE_WARNINGS"
	envKeyDisableTools    = "GHW_DISABLE_TOOLS"
	envKeyDisableTopology = "GHW_DISABLE_TOPOLOGY"
)

type Key string

var (
	defaultChroot          = "/"
	chrootKey              = Key("ghw.chroot")
	defaultToolsEnabled    = true
	toolsEnabledKey        = Key("ghw.tools.enabled")
	defaultTopologyEnabled = true
	topologyEnabledKey     = Key("ghw.topology.enabled")
	pcidbKey               = Key("ghw.pcidb")
	pathOverridesKey       = Key("ghw.path.overrides")
)

// Modifier sets some value on the context
type Modifier func(context.Context) context.Context

// WithChroot allows overriding the root directory ghw examines.
func WithChroot(path string) Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, chrootKey, path)
	}
}

// Chroot gets a context's chroot override or the default if none is set.
func Chroot(ctx context.Context) string {
	if ctx == nil {
		return defaultChroot
	}
	if v := ctx.Value(chrootKey); v != nil {
		return v.(string)
	}
	return defaultChroot
}

// EnvOrDefaultChroot returns the value of the GHW_CHROOT environs variable or
// the default value of "/" if not set
func EnvOrDefaultChroot() string {
	// Grab options from the environs by default
	if val, exists := os.LookupEnv(envKeyChroot); exists {
		return val
	}
	return defaultChroot
}

// WithDisableTools prevents ghw from calling external tools to discover
// hardware capabilities.
func WithDisableTools() Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, toolsEnabledKey, false)
	}
}

// ToolsEnabled returns true if external tools have been disabled.
func ToolsEnabled(ctx context.Context) bool {
	if ctx == nil {
		return defaultToolsEnabled
	}
	if v := ctx.Value(toolsEnabledKey); v != nil {
		return v.(bool)
	}
	return defaultToolsEnabled
}

// EnvOrDefaultDisableTools return true if ghw should use external tools to
// augment the data collected from sysfs. Most users want to do this most of
// time, so this is enabled by default.  Users consuming snapshots may want to
// opt out, thus they can set the GHW_DISABLE_TOOLS environs variable to any
// value to make ghw skip calling external tools even if they are available.
func EnvOrDefaultDisableTools() bool {
	if _, exists := os.LookupEnv(envKeyDisableTools); exists {
		return true
	}
	return false
}

// WithDisableTopology disables system topology detection to reduce memory
// consumption.  When using this option, ghw will skip scanning NUMA topology,
// CPU cores, memory caches, and node distances. This is useful when you only
// need basic PCI or GPU information and want to minimize memory overhead. The
// system architecture will be assumed to be SMP, and device Node fields will
// be nil.
func WithDisableTopology() Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, topologyEnabledKey, false)
	}
}

// EnvOrDefaultEnableTopology return true if ghw should detect system topology.
func EnvOrDefaultEnableTopology() bool {
	if _, exists := os.LookupEnv(envKeyDisableTopology); exists {
		return false
	}
	return defaultTopologyEnabled
}

// TopologyEnabled returns true if the detection of system topology is enabled.
func TopologyEnabled(ctx context.Context) bool {
	if ctx == nil {
		return defaultTopologyEnabled
	}
	if v := ctx.Value(topologyEnabledKey); v != nil {
		return v.(bool)
	}
	return defaultTopologyEnabled
}

// WithPCIDB allows you to provide a custom instance of the PCI database
// (pcidb.PCIDB) to ghw. This is useful if you want to use a preloaded or
// specially configured PCI database, such as one created with custom
// pcidb.WithOption settings, instead of letting ghw load the PCI database
// automatically.
func WithPCIDB(db *pcidb.PCIDB) Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pcidbKey, db)
	}
}

// PCIDB returns any PCIDB pointer set in the supplied context.
func PCIDB(ctx context.Context) *pcidb.PCIDB {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(pcidbKey); v != nil {
		return v.(*pcidb.PCIDB)
	}
	return nil
}

// WithPathOverrides supplies path-specific overrides for the context
func WithPathOverrides(overrides map[string]string) Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pathOverridesKey, overrides)
	}
}

// PathOverrides returns any path overrides set in the supplied context.
func PathOverrides(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(pathOverridesKey); v != nil {
		return v.(map[string]string)
	}
	return nil
}

// ContextFromEnv returns a new context.Context populated from the environs or
// default option values
func ContextFromEnv() context.Context {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, chrootKey, EnvOrDefaultChroot())
	ctx = context.WithValue(ctx, toolsEnabledKey, EnvOrDefaultDisableTools())
	ll := EnvOrDefaultLogLevel()
	logLevelVar.Set(ll)
	ctx = context.WithValue(ctx, logLevelKey, ll)
	disableWarn := EnvOrDefaultDisableWarnings()
	if disableWarn {
		ctx = WithDisableWarnings()(ctx)
	}
	useLogfmt := EnvOrDefaultLogLogfmt()
	if useLogfmt {
		ctx = WithLogLogfmt()(ctx)
	}
	return ctx
}

// fromOptions converts old-style pkg/option.Options by setting any Options
// fields on the supplied context.
//
// TODO(jaypipes): Remove this when we fully deprecate the old-style
// pkg/options stuff.
func fromOptions(ctx context.Context, opts *option.Options) context.Context {
	if opts == nil {
		return ctx
	}
	if opts.Chroot != "" {
		ctx = context.WithValue(ctx, chrootKey, opts.Chroot)
	}
	if opts.DisableTools {
		ctx = context.WithValue(ctx, toolsEnabledKey, false)
	}
	if opts.PCIDB != nil {
		ctx = context.WithValue(ctx, pcidbKey, opts.PCIDB)
	}
	if opts.PathOverrides != nil {
		ctx = context.WithValue(ctx, pathOverridesKey, opts.PathOverrides)
	}
	return ctx
}

// ContextFromArgs returns a context.Context populated with any old-style
// options or new-style arguments.
func ContextFromArgs(args ...any) context.Context {
	ctx := context.TODO()
	optsUsed := false
	opts := &option.Options{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case context.Context:
			ctx = arg
		case Modifier:
			ctx = arg(ctx)
		case option.Option:
			arg(opts)
			optsUsed = true
		}
	}
	if optsUsed {
		ctx = fromOptions(ctx, opts)
	}
	return ctx
}
