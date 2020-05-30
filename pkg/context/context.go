//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context

import (
	"github.com/jaypipes/ghw/pkg/option"
)

// Concrete merged set of configuration switches that act as an execution
// context when calling internal discovery methods
type Context struct {
	Chroot string
}

// New returns a Context struct pointer that has had various options set on it
func New(opts ...*option.Option) *Context {
	merged := option.Merge(opts...)
	return &Context{
		Chroot: *merged.Chroot,
	}
}

// FromEnv returns an Option that has been populated from the environs or
// default options values
func FromEnv() *Context {
	chrootVal := option.EnvOrDefaultChroot()
	return &Context{
		Chroot: chrootVal,
	}
}
