//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

const (
	TargetRoot = "ghw-snapshot-*"
)

// Unpack expands the given snapshot in a temporary directory managed by `ghw`.
// Returns the path of that directory. Callers are responsible for cleaning up
// the temporary directory.
func Unpack(snapshotName string) (string, error) {
	targetRoot, err := os.MkdirTemp("", TargetRoot)
	if err != nil {
		return "", err
	}
	err = UnpackInto(snapshotName, targetRoot)
	return targetRoot, err
}

// UnpackInto expands the given snapshot in a client-supplied directory.
// Returns true if the snapshot was actually unpacked, false otherwise
func UnpackInto(snapshotName, targetRoot string) error {
	snap, err := os.Open(snapshotName)
	if err != nil {
		return err
	}
	defer snap.Close()
	return Untar(targetRoot, snap)
}

// Untar extracts data from the given reader (providing data in tar.gz format) and unpacks it in the given directory.
func Untar(root string, r io.Reader) error {
	var err error
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// we are done
			return nil
		}

		if err != nil {
			// bail out
			return err
		}

		if header == nil {
			// TODO: how come?
			continue
		}

		target := filepath.Join(root, header.Name)
		mode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, mode)
			if err != nil {
				return err
			}

		case tar.TypeReg:
			dst, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, mode)
			if err != nil {
				return err
			}

			_, err = io.Copy(dst, tr)
			if err != nil {
				return err
			}

			dst.Close()

		case tar.TypeSymlink:
			err = os.Symlink(header.Linkname, target)
			if err != nil {
				return err
			}
		}
	}
}

func isEmptyDir(name string) bool {
	entries, err := os.ReadDir(name)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
