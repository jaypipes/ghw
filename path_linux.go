// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"os"
	"path/filepath"
)

const (
	DEFAULT_ROOT_PATH = "/"
)

// To facilitate querying of sysfs filesystems that are bind-mounted to a
// non-default root mountpoint, we allow users to set the GHW_CHROOT environ
// vairable to an alternate mountpoint. For instance, assume that the user of
// ghw is a Golang binary being executed from an application container that has
// certain host filesystems bind-mounted into the container at /host. The user
// would ensure the GHW_CHROOT environ variable is set to "/host" and ghw will
// build its paths from that location instead of /
func pathRoot() string {
	path := DEFAULT_ROOT_PATH
	if override, exists := os.LookupEnv("GHW_CHROOT"); exists {
		path = override
	}
	return path
}

func pathProcCpuinfo() string {
	return filepath.Join(pathRoot(), "proc", "cpuinfo")
}

func pathEtcMtab() string {
	return filepath.Join(pathRoot(), "etc", "mtab")
}

func pathSysBlock() string {
	return filepath.Join(pathRoot(), "sys", "block")
}

func pathSysDevicesSystemNode() string {
	return filepath.Join(pathRoot(), "sys", "devices", "system", "node")
}

func pathSysBusPciDevices() string {
	return filepath.Join(pathRoot(), "sys", "bus", "pci", "devices")
}

func pathSysClassDrm() string {
	return filepath.Join(pathRoot(), "sys", "class", "drm")
}

func pathSysClassNet() string {
	return filepath.Join(pathRoot(), "sys", "class", "net")
}

func pathRunUdevData() string {
	return filepath.Join(pathRoot(), "run", "udev", "data")
}
