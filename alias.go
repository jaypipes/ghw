//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/memory"
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

type MemoryInfo = memory.Info
type MemoryCacheType = memory.CacheType
type MemoryModule = memory.Module

const (
	MEMORY_CACHE_TYPE_UNIFIED     = memory.CACHE_TYPE_UNIFIED
	MEMORY_CACHE_TYPE_INSTRUCTION = memory.CACHE_TYPE_INSTRUCTION
	MEMORY_CACHE_TYPE_DATA        = memory.CACHE_TYPE_DATA
)

var (
	Memory = memory.New
)
