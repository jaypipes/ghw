// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func cachesForNode(nodeId uint32) ([]*MemoryCache, error) {
	// The /sys/devices/node/nodeX directory contains a subdirectory called
	// 'cpuX' for each logical processor assigned to the node. Each of those
	// subdirectories containers a 'cache' subdirectory which contains a number
	// of subdirectories beginning with 'index' and ending in the cache's
	// internal 0-based identifier. Those subdirectories contain a number of
	// files, including 'shared_cpu_list', 'size', and 'type' which we use to
	// determine cache characteristics.
	path := filepath.Join(
		PATH_DEVICES_SYSTEM_NODE,
		fmt.Sprintf("node%d", nodeId),
	)
	caches := make(map[string]*MemoryCache, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "cpu") {
			continue
		}
		if filename == "cpumap" || filename == "cpulist" {
			// There are two files in the node directory that start with 'cpu'
			// but are not subdirectories ('cpulist' and 'cpumap'). Ignore
			// these files.
			continue
		}
		// Grab the logical processor ID by cutting the integer from the
		// /sys/devices/system/node/nodeX/cpuX filename
		cpuPath := filepath.Join(path, filename)
		lpId, _ := strconv.Atoi(filename[3:])

		// Inspect the caches for each logical processor. There will be a
		// /sys/devices/system/node/nodeX/cpuX/cache directory containing a
		// number of directories beginning with the prefix "index" followed by
		// a number. The number indicates the level of the cache, which
		// indicates the "distance" from the processor. Each of these
		// directories contains information about the size of that level of
		// cache and the processors mapped to it.
		cachePath := filepath.Join(cpuPath, "cache")
		cacheDirFiles, err := ioutil.ReadDir(cachePath)
		if err != nil {
			return nil, err
		}
		for _, cacheDirFile := range cacheDirFiles {
			cacheDirFileName := cacheDirFile.Name()
			if !strings.HasPrefix(cacheDirFileName, "index") {
				continue
			}

			typePath := filepath.Join(cachePath, cacheDirFileName, "type")
			cacheTypeContents, err := ioutil.ReadFile(typePath)
			if err != nil {
				continue
			}
			cacheType := UNIFIED
			switch string(cacheTypeContents[:len(cacheTypeContents)-1]) {
			case "Data":
				cacheType = DATA
			case "Instruction":
				cacheType = INSTRUCTION
			default:
				cacheType = UNIFIED
			}

			levelPath := filepath.Join(cachePath, cacheDirFileName, "level")
			levelContents, err := ioutil.ReadFile(levelPath)
			if err != nil {
				continue
			}
			// levelContents is now a []byte with the last byte being a newline
			// character. Trim that off and convert the contents to an integer.
			level, _ := strconv.Atoi(string(levelContents[:len(levelContents)-1]))

			sizePath := filepath.Join(cachePath, cacheDirFileName, "size")
			sizeContents, err := ioutil.ReadFile(sizePath)
			if err != nil {
				continue
			}
			// size comes as XK\n, so we trim off the K and the newline.
			size, _ := strconv.Atoi(string(sizeContents[:len(sizeContents)-2]))

			scpuPath := filepath.Join(
				cachePath,
				cacheDirFileName,
				"shared_cpu_map",
			)
			sharedCpuMap, err := ioutil.ReadFile(scpuPath)
			if err != nil {
				continue
			}
			// The cache information is repeated for each node, so here, we
			// just ensure that we only have a one MemoryCache object for each
			// unique combination of level, type and processor map
			cacheKey := fmt.Sprintf("%d-%d-%s", level, cacheType, sharedCpuMap[:len(sharedCpuMap)-1])
			if cache, ok := caches[cacheKey]; !ok {
				cache = &MemoryCache{
					Level:             uint8(level),
					Type:              cacheType,
					SizeBytes:         uint64(size) * uint64(KB),
					LogicalProcessors: make([]uint32, 0),
				}
				caches[cacheKey] = cache
			}
			cache := caches[cacheKey]
			cache.LogicalProcessors = append(
				cache.LogicalProcessors,
				uint32(lpId),
			)
		}
	}

	cacheVals := make([]*MemoryCache, len(caches))
	x := 0
	for _, c := range caches {
		// ensure the cache's processor set is sorted by logical process ID
		sort.Sort(SortByLogicalProcessorId(c.LogicalProcessors))
		cacheVals[x] = c
		x++
	}

	return cacheVals, nil
}
