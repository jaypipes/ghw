//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package testdata

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func SnapshotsDirectory() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Cannot retrieve testdata directory")
	}
	basedir := filepath.Dir(file)
	return filepath.Join(basedir, "snapshots"), nil
}

func SamplesDirectory() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Cannot retrieve testdata directory")
	}
	basedir := filepath.Dir(file)
	return filepath.Join(basedir, "samples"), nil
}

func PCIDBChroot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve testdata directory")
	}
	basedir := filepath.Dir(file)
	return filepath.Join(basedir, "usr", "share", "hwdata", "pci.ids")
}
