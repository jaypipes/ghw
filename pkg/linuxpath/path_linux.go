// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package linuxpath

import (
	"fmt"
	"path/filepath"

	"github.com/jaypipes/ghw/pkg/option"
)

// PathRoots holds the roots of all the filesystem subtrees
// ghw wants to access.
type PathRoots struct {
	Etc  string
	Proc string
	Run  string
	Sys  string
	Var  string
}

// DefaultPathRoots return the canonical default value for PathRoots
func DefaultPathRoots() PathRoots {
	return PathRoots{
		Etc:  "/etc",
		Proc: "/proc",
		Run:  "/run",
		Sys:  "/sys",
		Var:  "/var",
	}
}

// PathRootsFromContext initialize PathRoots from the given Context,
// allowing overrides of the canonical default paths.
func PathRootsFromContext(opts *option.Options) PathRoots {
	roots := DefaultPathRoots()
	if pathEtc, ok := opts.PathOverrides["/etc"]; ok {
		roots.Etc = pathEtc
	}
	if pathProc, ok := opts.PathOverrides["/proc"]; ok {
		roots.Proc = pathProc
	}
	if pathRun, ok := opts.PathOverrides["/run"]; ok {
		roots.Run = pathRun
	}
	if pathSys, ok := opts.PathOverrides["/sys"]; ok {
		roots.Sys = pathSys
	}
	if pathVar, ok := opts.PathOverrides["/var"]; ok {
		roots.Var = pathVar
	}
	return roots
}

type Paths struct {
	SysRoot                string
	VarLog                 string
	ProcMeminfo            string
	ProcCpuinfo            string
	ProcMounts             string
	SysKernelMMHugepages   string
	SysBlock               string
	SysDevicesSystemNode   string
	SysDevicesSystemMemory string
	SysDevicesSystemCPU    string
	SysBusPciDevices       string
	SysBusUsbDevices       string
	SysClassDRM            string
	SysClassDMI            string
	SysClassNet            string
	RunUdevData            string
}

// New returns a new Paths struct containing filepath fields relative to the
// supplied Context
func New(opts *option.Options) *Paths {
	roots := PathRootsFromContext(opts)
	return &Paths{
		SysRoot:                filepath.Join(opts.Chroot, roots.Sys),
		VarLog:                 filepath.Join(opts.Chroot, roots.Var, "log"),
		ProcMeminfo:            filepath.Join(opts.Chroot, roots.Proc, "meminfo"),
		ProcCpuinfo:            filepath.Join(opts.Chroot, roots.Proc, "cpuinfo"),
		ProcMounts:             filepath.Join(opts.Chroot, roots.Proc, "self", "mounts"),
		SysKernelMMHugepages:   filepath.Join(opts.Chroot, roots.Sys, "kernel", "mm", "hugepages"),
		SysBlock:               filepath.Join(opts.Chroot, roots.Sys, "block"),
		SysDevicesSystemNode:   filepath.Join(opts.Chroot, roots.Sys, "devices", "system", "node"),
		SysDevicesSystemMemory: filepath.Join(opts.Chroot, roots.Sys, "devices", "system", "memory"),
		SysDevicesSystemCPU:    filepath.Join(opts.Chroot, roots.Sys, "devices", "system", "cpu"),
		SysBusPciDevices:       filepath.Join(opts.Chroot, roots.Sys, "bus", "pci", "devices"),
		SysBusUsbDevices:       filepath.Join(opts.Chroot, roots.Sys, "bus", "usb", "devices"),
		SysClassDRM:            filepath.Join(opts.Chroot, roots.Sys, "class", "drm"),
		SysClassDMI:            filepath.Join(opts.Chroot, roots.Sys, "class", "dmi"),
		SysClassNet:            filepath.Join(opts.Chroot, roots.Sys, "class", "net"),
		RunUdevData:            filepath.Join(opts.Chroot, roots.Run, "udev", "data"),
	}
}

func (p *Paths) NodeCPU(nodeID int, lpID int) string {
	return filepath.Join(
		p.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
		fmt.Sprintf("cpu%d", lpID),
	)
}

func (p *Paths) NodeCPUCache(nodeID int, lpID int) string {
	return filepath.Join(
		p.NodeCPU(nodeID, lpID),
		"cache",
	)
}

func (p *Paths) NodeCPUCacheIndex(nodeID int, lpID int, cacheIndex int) string {
	return filepath.Join(
		p.NodeCPUCache(nodeID, lpID),
		fmt.Sprintf("index%d", cacheIndex),
	)
}
