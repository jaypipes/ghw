//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"strconv"
	"strings"
)

type MemoryCacheType int

const (
	UNIFIED MemoryCacheType = iota
	INSTRUCTION
	DATA
)

type SortByMemoryCacheLevel []*MemoryCache

func (a SortByMemoryCacheLevel) Len() int      { return len(a) }
func (a SortByMemoryCacheLevel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByMemoryCacheLevel) Less(i, j int) bool {
	if a[i].Level < a[j].Level {
		return true
	} else if a[i].Level == a[j].Level {
		return a[i].Type < a[j].Type
	}
	return false
}

type SortByLogicalProcessorId []uint32

func (a SortByLogicalProcessorId) Len() int           { return len(a) }
func (a SortByLogicalProcessorId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByLogicalProcessorId) Less(i, j int) bool { return a[i] < a[j] }

type MemoryCache struct {
	Level     uint8
	Type      MemoryCacheType
	SizeBytes uint64
	// The set of logical processors (hardware threads) that have access to the
	// cache
	LogicalProcessors []uint32
}

func (c *MemoryCache) String() string {
	sizeKb := c.SizeBytes / uint64(KB)
	typeStr := ""
	if c.Type == INSTRUCTION {
		typeStr = "i"
	} else if c.Type == DATA {
		typeStr = "d"
	}
	cacheIdStr := fmt.Sprintf("L%d%s", c.Level, typeStr)
	processorMapStr := ""
	if c.LogicalProcessors != nil {
		lpStrings := make([]string, len(c.LogicalProcessors))
		for x, lpid := range c.LogicalProcessors {
			lpStrings[x] = strconv.Itoa(int(lpid))
		}
		processorMapStr = " shared with logical processors: " + strings.Join(lpStrings, ",")
	}
	return fmt.Sprintf(
		"%s cache (%d KB)%s",
		cacheIdStr,
		sizeKb,
		processorMapStr,
	)
}
