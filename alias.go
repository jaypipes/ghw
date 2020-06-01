//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/option"
)

type WithOption = option.Option

var (
	WithChroot = option.WithChroot
)

type CPUInfo = cpu.Info

var (
	CPU = cpu.New
)
