// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func (ctx *context) memFillInfo(info *MemoryInfo) error {
	tub := ctx.memTotalUsableBytes()
	if tub < 1 {
		return fmt.Errorf("Could not determine total usable bytes of memory")
	}
	info.TotalUsableBytes = tub
	tpb, err := ctx.memTotalPhysicalBytes()
	if err != nil {
		info.TotalPhysicalBytes = tub
		errMsg := fmt.Sprintf("fallback to total usable bytes after error\n"+
			"getting total physical bytes of RAM:\n"+
			"%v", err)
		warn(errMsg)
	} else {
		info.TotalPhysicalBytes = int64(tpb)
	}

	info.SupportedPageSizes = ctx.memSupportedPageSizes()
	return nil
}

func (ctx *context) memTotalPhysicalBytes() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	return info.Totalram, err
}

func (ctx *context) memTotalUsableBytes() int64 {
	// In Linux, /proc/meminfo contains a set of memory-related amounts, with
	// lines looking like the following:
	//
	// $ cat /proc/meminfo
	// MemTotal:       24677596 kB
	// MemFree:        21244356 kB
	// MemAvailable:   22085432 kB
	// ...
	// HugePages_Total:       0
	// HugePages_Free:        0
	// HugePages_Rsvd:        0
	// HugePages_Surp:        0
	// ...
	//
	// It's worth noting that /proc/meminfo returns exact information, not
	// "theoretical" information. For instance, on the above system, I have
	// 24GB of RAM but MemTotal is indicating only around 23GB. This is because
	// MemTotal contains the exact amount of *usable* memory after accounting
	// for the kernel's resident memory size and a few reserved bits. For more
	// information, see:
	//
	//  https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	filePath := ctx.pathProcMeminfo()
	r, err := os.Open(filePath)
	if err != nil {
		return -1
	}
	defer safeClose(r)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		key := strings.Trim(parts[0], ": \t")
		if key != "MemTotal" {
			continue
		}
		value, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return -1
		}
		inKb := (len(parts) == 3 && strings.TrimSpace(parts[2]) == "kB")
		if inKb {
			value = value * int(KB)
		}
		return int64(value)
	}
	return -1
}

func (ctx *context) memSupportedPageSizes() []uint64 {
	// In Linux, /sys/kernel/mm/hugepages contains a directory per page size
	// supported by the kernel. The directory name corresponds to the pattern
	// 'hugepages-{pagesize}kb'
	dir := ctx.pathSysKernelMMHugepages()
	out := make([]uint64, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, file := range files {
		parts := strings.Split(file.Name(), "-")
		sizeStr := parts[1]
		// Cut off the 'kb'
		sizeStr = sizeStr[0 : len(sizeStr)-2]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return out
		}
		out = append(out, uint64(size*int(KB)))
	}
	return out
}
