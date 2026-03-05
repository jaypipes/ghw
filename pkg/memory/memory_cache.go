//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/unitutil"
)

// CacheType indicates the type of memory stored in a memory cache.
type CacheType int

const (
	// CacheTypeUnified indicates the memory cache stores both instructions and
	// data.
	CacheTypeUnified CacheType = iota
	// CacheTypeInstruction indicates the memory cache stores only instructions
	// (executable bytecode).
	CacheTypeInstruction
	// CacheTypeData indicates the memory cache stores only data
	// (non-executable bytecode).
	CacheTypeData
)

const (
	// DEPRECATED: Please use CacheTypeUnified
	CACHE_TYPE_UNIFIED = CacheTypeUnified
	// DEPRECATED: Please use CacheTypeUnified
	CACHE_TYPE_INSTRUCTION = CacheTypeInstruction
	// DEPRECATED: Please use CacheTypeUnified
	CACHE_TYPE_DATA = CacheTypeData
)

var (
	memoryCacheTypeString = map[CacheType]string{
		CacheTypeUnified:     "Unified",
		CacheTypeInstruction: "Instruction",
		CacheTypeData:        "Data",
	}

	// NOTE(fromani): the keys are all lowercase and do not match
	// the keys in the opposite table `memoryCacheTypeString`.
	// This is done because of the choice we made in
	// CacheType:MarshalJSON.
	// We use this table only in UnmarshalJSON, so it should be OK.
	stringMemoryCacheType = map[string]CacheType{
		"unified":     CacheTypeUnified,
		"instruction": CacheTypeInstruction,
		"data":        CacheTypeData,
	}
)

func (a CacheType) String() string {
	return memoryCacheTypeString[a]
}

// NOTE(jaypipes): since serialized output is as "official" as we're going to
// get, let's lowercase the string output when serializing, in order to
// "normalize" the expected serialized output
func (a CacheType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.ToLower(a.String()))), nil
}

func (a *CacheType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	key := strings.ToLower(s)
	val, ok := stringMemoryCacheType[key]
	if !ok {
		return fmt.Errorf("unknown memory cache type: %q", key)
	}
	*a = val
	return nil
}

type SortByCacheLevelTypeFirstProcessor []*Cache

func (a SortByCacheLevelTypeFirstProcessor) Len() int      { return len(a) }
func (a SortByCacheLevelTypeFirstProcessor) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByCacheLevelTypeFirstProcessor) Less(i, j int) bool {
	if a[i].Level < a[j].Level {
		return true
	} else if a[i].Level == a[j].Level {
		if a[i].Type < a[j].Type {
			return true
		} else if a[i].Type == a[j].Type {
			// NOTE(jaypipes): len(LogicalProcessors) is always >0 and is always
			// sorted lowest LP ID to highest LP ID
			return a[i].LogicalProcessors[0] < a[j].LogicalProcessors[0]
		}
	}
	return false
}

type SortByLogicalProcessorId []uint32

func (a SortByLogicalProcessorId) Len() int           { return len(a) }
func (a SortByLogicalProcessorId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByLogicalProcessorId) Less(i, j int) bool { return a[i] < a[j] }

// Cache contains information about a single memory cache on a physical CPU
// package. Caches have a 1-based numeric level, with lower numbers indicating
// the cache is "closer" to the processing cores and reading memory from the
// cache will be faster relative to caches with higher levels. Note that this
// has nothing to do with RAM or memory modules like DIMMs.
type Cache struct {
	// Level is a 1-based numeric level that indicates the relative closeness
	// of this cache to processing cores on the physical package. Lower numbers
	// are "closer" to the processing cores and therefore have faster access
	// times.
	Level uint8 `json:"level"`
	// Type indicates what type of memory is stored in the cache. Can be
	// instruction (executable bytecodes), data or both.
	Type CacheType `json:"type"`
	// SizeBytes indicates the size of the cache in bytes.
	SizeBytes uint64 `json:"size_bytes"`
	// The set of logical processors (hardware threads) that have access to
	// this cache.
	LogicalProcessors []uint32 `json:"logical_processors"`
}

func (c *Cache) String() string {
	sizeKb := c.SizeBytes / uint64(unitutil.KB)
	typeStr := ""
	if c.Type == CacheTypeInstruction {
		typeStr = "i"
	} else if c.Type == CacheTypeData {
		typeStr = "d"
	}
	cacheIDStr := fmt.Sprintf("L%d%s", c.Level, typeStr)
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
		cacheIDStr,
		sizeKb,
		processorMapStr,
	)
}
