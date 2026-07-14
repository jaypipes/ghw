//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package config_test

import (
	"context"
	"testing"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/option"
)

// TestContextFromArgsHonorsEnvChrootWithModifiers ensures that GHW_CHROOT is
// honored even when callers also pass Modifiers. Regression test for the
// behavior change introduced in 0.23.0 where, when args were supplied, the
// base context was context.TODO() and the GHW_CHROOT environment variable was
// silently ignored.
func TestContextFromArgsHonorsEnvChrootWithModifiers(t *testing.T) {
	t.Setenv("GHW_CHROOT", "/from-env")

	ctx := config.ContextFromArgs(config.WithDisableTools())

	if got := config.Chroot(ctx); got != "/from-env" {
		t.Fatalf("Expected chroot to come from GHW_CHROOT (/from-env), got %q", got)
	}
	if config.ToolsEnabled(ctx) {
		t.Fatalf("Expected WithDisableTools modifier to be applied")
	}
}

// TestContextFromArgsExplicitChrootOverridesEnv ensures that a chroot passed
// via WithChroot wins over GHW_CHROOT.
func TestContextFromArgsExplicitChrootOverridesEnv(t *testing.T) {
	t.Setenv("GHW_CHROOT", "/from-env")

	ctx := config.ContextFromArgs(config.WithChroot("/from-arg"))

	if got := config.Chroot(ctx); got != "/from-arg" {
		t.Fatalf("Expected explicit WithChroot to override env, got %q", got)
	}
}

// TestContextFromArgsExplicitContextReplacesBase ensures that an explicit
// context.Context arg fully replaces the env-derived base context.
func TestContextFromArgsExplicitContextReplacesBase(t *testing.T) {
	t.Setenv("GHW_CHROOT", "/from-env")

	ctx := config.ContextFromArgs(context.Background())

	if got := config.Chroot(ctx); got != "/" {
		t.Fatalf("Expected explicit context.Background() to replace env-derived base (default chroot \"/\"), got %q", got)
	}
}

// TestContextFromArgsOldStyleOptionChrootOverridesEnv ensures that the
// old-style option.WithChroot also overrides GHW_CHROOT.
func TestContextFromArgsOldStyleOptionChrootOverridesEnv(t *testing.T) {
	t.Setenv("GHW_CHROOT", "/from-env")

	ctx := config.ContextFromArgs(option.WithChroot("/from-old-opt"))

	if got := config.Chroot(ctx); got != "/from-old-opt" {
		t.Fatalf("Expected old-style option.WithChroot to override env, got %q", got)
	}
}
