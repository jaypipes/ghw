//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"math"
)

type MemoryInfo struct {
	TotalPhysicalBytes int64
	TotalUsableBytes   int64
	// An array of sizes, in bytes, of memory pages supported by the host
	SupportedPageSizes []uint64
}

func Memory() (*MemoryInfo, error) {
	info := &MemoryInfo{}
	err := memFillInfo(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (i *MemoryInfo) String() string {
	tpbs := "unknown"
	if i.TotalPhysicalBytes > 0 {
		tpb := i.TotalPhysicalBytes
		unit, unitStr := unitWithString(tpb)
		tpb = int64(math.Ceil(float64(i.TotalPhysicalBytes) / float64(unit)))
		tpbs = fmt.Sprintf("%d%s", tpb, unitStr)
	}
	tubs := "unknown"
	if i.TotalUsableBytes > 0 {
		tub := i.TotalUsableBytes
		unit, unitStr := unitWithString(tub)
		tub = int64(math.Ceil(float64(i.TotalUsableBytes) / float64(unit)))
		tubs = fmt.Sprintf("%d%s", tub, unitStr)
	}
	return fmt.Sprintf("memory (%s physical, %s usable)", tpbs, tubs)
}
