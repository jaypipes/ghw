//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// warning: don't use the context package here, this means not even the linuxpath package.
	// TODO(fromani) remove the path duplication
	sysClassNet = "/sys/class/net"
)

// ExpectedCloneNetContent returns a slice of strings pertaning to the network interfaces ghw
// cares about. We cannot use a static list because we want to filter away the virtual devices,
// which  ghw doesn't concern itself about. So we need to do some runtime discovery.
// Additionally, we want to make sure to clone the backing device data.
func ExpectedCloneNetContent() []string {
	var fileSpecs []string
	ifaceEntries := []string{
		"addr_assign_type",
		// intentionally avoid to clone "address" to avoid to leak any host-idenfifiable data.
	}
	entries, err := ioutil.ReadDir(sysClassNet)
	if err != nil {
		// we should not import context, hence we can't Warn()
		return fileSpecs
	}
	for _, entry := range entries {
		netName := entry.Name()
		netPath := filepath.Join(sysClassNet, netName)
		dest, err := os.Readlink(netPath)
		if err != nil {
			continue
		}
		if strings.Contains(dest, "devices/virtual/net") {
			// there is no point in cloning data for virtual devices,
			// because ghw concerns itself with HardWare.
			continue
		}

		// so, first copy the symlink itself
		fileSpecs = append(fileSpecs, netPath)

		// now we have to clone the content of the actual network interface
		// data related (and found into a subdir of) the backing hardware
		// device
		netIface := filepath.Clean(filepath.Join(sysClassNet, dest))
		for _, ifaceEntry := range ifaceEntries {
			fileSpecs = append(fileSpecs, filepath.Join(netIface, ifaceEntry))
		}

	}

	return fileSpecs
}
