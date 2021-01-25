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
)

// Attempting to tar up pseudofiles like /proc/cpuinfo is an exercise in
// futility. Notably, the pseudofiles, when read by syscalls, do not return the
// number of bytes read. This causes the tar writer to write zero-length files.
//
// Instead, it is necessary to build a directory structure in a tmpdir and
// create actual files with copies of the pseudofile contents

// CloneTreeInto copies all the pseudofiles that ghw will consume into the root
// `scratchDir`, preserving the hieratchy.
func CloneTreeInto(scratchDir string) error {
	var err error

	var createPaths = []string{
		"sys/block",
	}

	for _, path := range createPaths {
		if err = os.MkdirAll(filepath.Join(scratchDir, path), os.ModePerm); err != nil {
			return err
		}
	}

	if err = createBlockDevices(scratchDir); err != nil {
		return err
	}

	fileSpecs := ExpectedCloneContent()
	return CopyFilesInto(fileSpecs, scratchDir, nil)
}

// ExpectedCloneContent return a slice of glob patterns which represent the pseudofiles
// ghw cares about. The intended usage of this function is to validate a clone tree,
// checking that the content matches the expectations.
func ExpectedCloneContent() []string {
	return []string{
		"/etc/mtab",
		"/proc/cpuinfo",
		"/proc/meminfo",
		"/sys/bus/pci/devices/*",
		"/sys/devices/pci*/*/irq",
		"/sys/devices/pci*/*/local_cpulist",
		"/sys/devices/pci*/*/modalias",
		"/sys/devices/pci*/*/numa_node",
		"/sys/devices/pci*/pci_bus/*/cpulistaffinity",
		"/sys/devices/system/cpu/cpu*/cache/index*/*",
		"/sys/devices/system/cpu/cpu*/topology/*",
		"/sys/devices/system/memory/block_size_bytes",
		"/sys/devices/system/memory/memory*/online",
		"/sys/devices/system/memory/memory*/state",
		"/sys/devices/system/node/has_*",
		"/sys/devices/system/node/online",
		"/sys/devices/system/node/possible",
		"/sys/devices/system/node/node*/cpu*",
		"/sys/devices/system/node/node*/distance",
	}
}

// ValidateClonedTree checks the content of a cloned tree, whose root is `clonedDir`,
// against a slice of glob specs which must be included in the cloned tree.
// Is not wrong, and this functions doesn't enforce this, that the cloned tree includes
// more files than the necessary; ghw will just ignore the files it doesn't care about.
// Returns a slice of glob patters expected (given) but not found in the cloned tree,
// and the error during the validation (if any).
func ValidateClonedTree(fileSpecs []string, clonedDir string) ([]string, error) {
	missing := []string{}
	for _, fileSpec := range fileSpecs {
		matches, err := filepath.Glob(filepath.Join(clonedDir, fileSpec))
		if err != nil {
			return missing, err
		}
		if len(matches) == 0 {
			missing = append(missing, fileSpec)
		}
	}
	return missing, nil
}

func copyPseudoFile(path, targetPath string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	trace("creating %s\n", targetPath)
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	if _, err = f.Write(buf); err != nil {
		return err
	}
	f.Close()
	return nil
}

// CopyFileOptions allows to finetune the behaviour of the CopyFilesInto function
type CopyFileOptions struct {
	// IsSymlinkFn allows to control the behaviour when handling a symlink.
	// If this hook returns true, the source file is treated as symlink: the cloned
	// tree will thus contain a symlink, with its path adjusted to match the relative
	// path inside the cloned tree. If return false, the symlink will be deferred.
	// The easiest use case of this hook is if you want to avoid symlinks in your cloned
	// tree (having duplicated content). In this case you can just add a function
	// which always return false.
	IsSymlinkFn func(path string, info os.FileInfo) bool
}

// CopyFilesInto copies all the given glob files specs in the given `destDir` directory,
// preserving the directory structure. This means you can provide a deeply nested filespec
// like
// - /some/deeply/nested/file*
// and you DO NOT need to build the tree incrementally like
// - /some/
// - /some/deeply/
// ...
// all glob patterns supported in `filepath.Glob` are supported.
func CopyFilesInto(fileSpecs []string, destDir string, opts *CopyFileOptions) error {
	if opts == nil {
		opts = &CopyFileOptions{
			IsSymlinkFn: isSymlink,
		}
	}
	for _, fileSpec := range fileSpecs {
		trace("copying spec: %q\n", fileSpec)
		matches, err := filepath.Glob(fileSpec)
		if err != nil {
			return err
		}
		if err := copyFileTreeInto(matches, destDir, opts); err != nil {
			return err
		}
	}
	return nil
}

func copyFileTreeInto(paths []string, destDir string, opts *CopyFileOptions) error {
	for _, path := range paths {
		trace("  copying path: %q\n", path)
		baseDir := filepath.Dir(path)
		if err := os.MkdirAll(filepath.Join(destDir, baseDir), os.ModePerm); err != nil {
			return err
		}

		fi, err := os.Lstat(path)
		if err != nil {
			return err
		}
		// directories must be listed explicitely and created separately.
		// In the future we may want to expose this decision as hook point in
		// CopyFileOptions, when clear use cases emerge.
		if fi.IsDir() {
			trace("expanded glob path %q is a directory - skipped", path)
			continue
		}
		if opts.IsSymlinkFn(path, fi) {
			trace("    copying link: %q\n", path)
			if err := copyLink(path, filepath.Join(destDir, path)); err != nil {
				return err
			}
		} else {
			trace("    copying file: %q\n", path)
			if err := copyPseudoFile(path, filepath.Join(destDir, path)); err != nil {
				return err
			}
		}
	}
	return nil
}

func isSymlink(path string, fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}

func copyLink(path, targetPath string) error {
	target, err := os.Readlink(path)
	if err != nil {
		return err
	}
	if err := os.Symlink(target, targetPath); err != nil {
		return err
	}

	return nil
}
