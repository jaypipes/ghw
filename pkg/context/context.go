//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/snapshot"
)

// Concrete merged set of configuration switches that act as an execution
// context when calling internal discovery methods
type Context struct {
	Chroot            string
	EnableTools       bool
	SnapshotPath      string
	SnapshotRoot      string
	SnapshotExclusive bool
	alert             option.Alerter
	err               error
	parent            *Context
}

// New returns a Context struct pointer that has had various options set on it
func New(opts ...*option.Option) *Context {
	merged := option.Merge(opts...)
	ctx := &Context{
		alert:  option.EnvOrDefaultAlerter(),
		Chroot: *merged.Chroot,
	}

	if merged.Snapshot != nil {
		ctx.SnapshotPath = merged.Snapshot.Path
		// root is optional, so a extra check is warranted
		if merged.Snapshot.Root != nil {
			ctx.SnapshotRoot = *merged.Snapshot.Root
		}
		ctx.SnapshotExclusive = merged.Snapshot.Exclusive
	}

	if merged.Alerter != nil {
		ctx.alert = merged.Alerter
	}

	if merged.EnableTools != nil {
		ctx.EnableTools = *merged.EnableTools
	}

	// New is not allowed to return error - it would break the established API.
	// so the only way out is to actually do the checks here and record the error,
	// and return it later, at the earliest possible occasion, in Setup()
	if ctx.SnapshotPath != "" && ctx.Chroot != option.DefaultChroot {
		// The env/client code supplied a value, but we are will overwrite it when unpacking shapshots!
		ctx.err = fmt.Errorf("Conflicting options: chroot %q and snapshot path %q", ctx.Chroot, ctx.SnapshotPath)
	}
	return ctx
}

// FromEnv returns a Context that has been populated from the environs or
// default options values
func FromEnv() *Context {
	chrootVal := option.EnvOrDefaultChroot()
	enableTools := option.EnvOrDefaultTools()
	snapPathVal := option.EnvOrDefaultSnapshotPath()
	snapRootVal := option.EnvOrDefaultSnapshotRoot()
	snapExclusiveVal := option.EnvOrDefaultSnapshotExclusive()
	return &Context{
		Chroot:            chrootVal,
		EnableTools:       enableTools,
		SnapshotPath:      snapPathVal,
		SnapshotRoot:      snapRootVal,
		SnapshotExclusive: snapExclusiveVal,
	}
}

// FromParent returns a Context which is subordinate to another one.
// This means only the parent Context is in charge to the environment (e.g. setting up the snapshot).
// Subordinate contexts will defer any environment-changing action to their parent.
func FromParent(ctx *Context) *Context {
	return &Context{
		Chroot:            ctx.Chroot,
		EnableTools:       ctx.EnableTools,
		SnapshotPath:      ctx.SnapshotPath,
		SnapshotRoot:      ctx.SnapshotRoot,
		SnapshotExclusive: ctx.SnapshotExclusive,
		alert:             ctx.alert,
		parent:            ctx,
	}
}

// Do wraps a Setup/Teardown pair around the given function
func (ctx *Context) Do(fn func() error) error {
	err := ctx.Setup()
	if err != nil {
		return err
	}
	defer ctx.Teardown()
	return fn()
}

// Setup prepares the extra optional data a Context may use.
// `Context`s are ready to use once returned by `New`. Optional features,
// like snapshot unpacking, may require extra steps. Run `Setup` to perform them.
// You should call `Setup` just once. It is safe to call `Setup` if you don't make
// use of optional extra features - `Setup` will do nothing.
func (ctx *Context) Setup() error {
	if ctx.err != nil {
		return ctx.err
	}
	if ctx.parent != nil {
		// nothing to do - the parent is in charge
		return nil
	}
	if ctx.SnapshotPath == "" {
		// nothing to do!
		return nil
	}

	var err error
	root := ctx.SnapshotRoot
	if root == "" {
		root, err = snapshot.Unpack(ctx.SnapshotPath)
	} else {
		var flags uint
		if ctx.SnapshotExclusive {
			flags |= snapshot.OwnTargetDirectory
		}
		_, err = snapshot.UnpackInto(ctx.SnapshotPath, root, flags)
	}
	if err != nil {
		return err
	}

	ctx.Chroot = root
	return nil
}

// Teardown releases any resource acquired by Setup.
// You should always call `Teardown` if you called `Setup` to free any resources
// acquired by `Setup`. Check `Do` for more automated management.
func (ctx *Context) Teardown() error {
	if ctx.parent != nil {
		// nothing to do - the parent is in charge
		return nil
	}
	if ctx.SnapshotRoot != "" {
		// if the client code provided the unpack directory,
		// then it is also in charge of the cleanup.
		return nil
	}
	return snapshot.Cleanup(ctx.SnapshotRoot)
}

func (ctx *Context) Warn(msg string, args ...interface{}) {
	ctx.alert.Printf("WARNING: "+msg, args...)
}
