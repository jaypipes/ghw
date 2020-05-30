// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"path/filepath"

	"github.com/jaypipes/ghw/pkg/context"
)

func pathVarLog(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "var", "log")
}

func pathProcMeminfo(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "proc", "meminfo")
}

func pathSysKernelMMHugepages(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "kernel", "mm", "hugepages")
}

func pathProcCpuinfo(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "proc", "cpuinfo")
}

func pathEtcMtab(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "etc", "mtab")
}

func pathSysBlock(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "block")
}

func pathSysDevicesSystemNode(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "devices", "system", "node")
}

func pathSysBusPciDevices(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "bus", "pci", "devices")
}

func pathSysClassDrm(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "class", "drm")
}

func pathSysClassDMI(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "class", "dmi")
}

func pathSysClassNet(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "sys", "class", "net")
}

func pathRunUdevData(ctx *context.Context) string {
	return filepath.Join(ctx.Chroot, "run", "udev", "data")
}
