//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"context"
	"os"
	"path/filepath"

	ghwcontext "github.com/jaypipes/ghw/pkg/context"
	pciaddr "github.com/jaypipes/ghw/pkg/pci/address"
)

const (
	// root directory: entry point to start scanning the PCI forest
	// warning: don't use the context package here, this means not even the linuxpath package.
	// TODO(fromani) remove the path duplication
	sysBusPCIDir = "/sys/bus/pci/devices"
)

// ExpectedClonePCIContent return a slice of glob patterns which represent the pseudofiles
// ghw cares about, pertaining to PCI devices only.
// Beware: the content is host-specific, because the PCI topology is host-dependent and unpredictable.
func ExpectedClonePCIContent(
	ctx context.Context,
) ([]string, error) {
	fileSpecs := []string{
		"/sys/bus/pci/drivers/*",
	}
	pciRoots := []string{
		sysBusPCIDir,
	}
	for {
		if len(pciRoots) == 0 {
			break
		}
		pciRoot := pciRoots[0]
		pciRoots = pciRoots[1:]
		specs, roots, err := scanPCIDeviceRoot(ctx, pciRoot)
		if err != nil {
			return nil, err
		}
		pciRoots = append(pciRoots, roots...)
		fileSpecs = append(fileSpecs, specs...)
	}
	return fileSpecs, nil
}

// scanPCIDeviceRoot reports a slice of glob patterns which represent the pseudofiles
// ghw cares about pertaining to all the PCI devices connected to the bus connected from the
// given root; usually (but not always) a CPU packages has 1+ PCI(e) roots, forming the first
// level; more PCI bridges are (usually) attached to this level, creating deep nested trees.
// hence we need to scan all possible roots, to make sure not to miss important devices.
//
// note about notifying errors. This function and its helper functions do use ghwcontext.Debug(ctx, ) everywhere
// to report recoverable errors, even though it would have been appropriate to use Warn().
// This is unfortunate, and again a byproduct of the fact we cannot use context.Context to avoid
// circular dependencies.
func scanPCIDeviceRoot(
	ctx context.Context,
	root string,
) (fileSpecs []string, pciRoots []string, err error) {
	ghwcontext.Debug(ctx, "scanning PCI device root %q", root)

	perDevEntries := []string{
		"class",
		"device",
		"driver",
		"iommu_group",
		"irq",
		"local_cpulist",
		"modalias",
		"numa_node",
		"revision",
		"vendor",
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, err
	}
	for _, entry := range entries {
		entryName := entry.Name()
		if addr := pciaddr.FromString(entryName); addr == nil {
			// doesn't look like a entry we care about
			// This is by far and large the most likely path
			// hence we should NOT trace/warn here.
			continue
		}

		entryPath := filepath.Join(root, entryName)
		pciEntry, err := findPCIEntryFromPath(ctx, root, entryName)
		if err != nil {
			ghwcontext.Debug(ctx, "error scanning %q: %v. skipping", entryName, err)
			continue
		}

		ghwcontext.Debug(ctx, "PCI entry is %q", pciEntry)
		fileSpecs = append(fileSpecs, entryPath)
		for _, perNetEntry := range perDevEntries {
			fileSpecs = append(fileSpecs, filepath.Join(pciEntry, perNetEntry))
		}

		bridge, err := isPCIBridge(entryPath)
		if err != nil {
			return nil, nil, err
		}
		if bridge {
			ghwcontext.Debug(ctx, "adding new PCI root %q", entryName)
			pciRoots = append(pciRoots, pciEntry)
		}
	}
	return fileSpecs, pciRoots, nil
}

func findPCIEntryFromPath(
	ctx context.Context,
	root string,
	entryName string,
) (string, error) {
	entryPath := filepath.Join(root, entryName)
	fi, err := os.Lstat(entryPath)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		// regular file, nothing to resolve
		return entryPath, nil
	}
	// resolve symlink
	target, err := os.Readlink(entryPath)
	if err != nil {
		return "", err
	}
	ghwcontext.Debug(ctx, "entry %q is symlink resolved to %q", entryPath, target)
	return filepath.Clean(filepath.Join(root, target)), nil
}

func isPCIBridge(entryPath string) (bool, error) {
	subNodes, err := os.ReadDir(entryPath)
	if err != nil {
		return false, err
	}
	for _, subNode := range subNodes {
		if !subNode.IsDir() {
			continue
		}
		if addr := pciaddr.FromString(subNode.Name()); addr != nil {
			// we got an entry in the directory pertaining to this device
			// which is a directory itself and it is named like a PCI address.
			// Hence we infer the device we are considering is a PCI bridge of sorts.
			// This is is indeed a bit brutal, but the only possible alternative
			// (besides blindly copying everything in /sys/bus/pci/devices) is
			// to detect the type of the device and pick only the bridges.
			// This approach duplicates the logic within the `pci` subkpg
			// - or forces us into awkward dep cycles, and has poorer forward
			// compatibility.
			return true, nil
		}
	}
	return false, nil
}
