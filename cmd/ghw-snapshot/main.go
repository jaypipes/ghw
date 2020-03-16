// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// version of application at compile time (-X 'main.version=$(VERSION)').
	version = "(Unknown Version)"
	// buildHash GIT hash of application at compile time (-X 'main.buildHash=$(GITCOMMIT)').
	buildHash = "No Git-hash Provided."
	// buildDate of application at compile time (-X 'main.buildDate=$(BUILDDATE)').
	buildDate = "No Build Date Provided."
	// show debug output
	debug = false
	// output filepath to save snapshot to
	outPath string
)

var (
	createPseudofilePaths = []string{
		"/proc/cpuinfo",
		"/proc/meminfo",
		"/etc/mtab",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghw-snapshot",
	Short: "ghw-snapshot - Snapshot filesystem containing system information.",
	RunE:  execute,
}

func trace(msg string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(msg, args...)
}

func systemFingerprint() string {
	hn, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	m := md5.New()
	io.WriteString(m, hn)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func defaultOutPath() string {
	fp := systemFingerprint()
	return fmt.Sprintf("%s-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH, fp)
}

func execute(cmd *cobra.Command, args []string) error {
	scratchDir, err := ioutil.TempDir("", "ghw-snapshot")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratchDir)

	var createPaths = []string{
		"proc",
		"etc",
		"sys/block",
	}

	for _, path := range createPaths {
		if err = os.MkdirAll(filepath.Join(scratchDir, path), os.ModePerm); err != nil {
			return err
		}
	}

	if err = createPseudofiles(scratchDir); err != nil {
		return err
	}
	if err = createBlockDevices(scratchDir); err != nil {
		return err
	}
	return doSnapshot(scratchDir)
}

// Attempting to tar up pseudofiles like /proc/cpuinfo is an exercise in
// futility. Notably, the pseudofiles, when read by syscalls, do not return the
// number of bytes read. This causes the tar writer to write zero-length files.
//
// Instead, it is necessary to build a directory structure in a tmpdir and
// create actual files with copies of the pseudofile contents
func createPseudofiles(buildDir string) error {
	for _, path := range createPseudofilePaths {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(buildDir, path)
		trace("creating %s\n", targetPath)
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		if _, err = f.Write(buf); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func createBlockDevices(buildDir string) error {
	// Grab all the block device pseudo-directories from /sys/block symlinks
	// (excluding loopback devices) and inject them into our build filesystem
	// with all but the circular symlink'd subsystem directories
	devLinks, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return err
	}
	for _, devLink := range devLinks {
		dname := devLink.Name()
		if strings.HasPrefix(dname, "loop") {
			continue
		}
		devPath := filepath.Join("/sys/block", dname)
		fi, err := os.Lstat(devPath)
		if err != nil {
			return err
		}
		var link string
		if fi.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(devPath)
			if err != nil {
				return err
			}
		}
		// Create a symlink in our build filesystem that is a directory
		// pointing to the actual device bus path where the block device's
		// information directory resides
		linkPath := filepath.Join(buildDir, "sys/block", dname)
		linkTargetPath := filepath.Join(
			buildDir,
			"sys/block",
			strings.TrimPrefix(link, string(os.PathSeparator)),
		)
		trace("creating device directory %s\n", linkTargetPath)
		if err = os.MkdirAll(linkTargetPath, os.ModePerm); err != nil {
			return err
		}
		trace("linking device directory %s to %s\n", linkPath, linkTargetPath)
		if err = os.Symlink(linkTargetPath, linkPath); err != nil {
			return err
		}
		// Now read the source block device directory and populate the
		// newly-created target link in the build directory with the
		// appropriate block device pseudofiles
		srcDeviceDir := filepath.Join(
			"/sys/block",
			strings.TrimPrefix(link, string(os.PathSeparator)),
		)
		if err = createBlockDeviceDir(linkTargetPath, srcDeviceDir); err != nil {
			return err
		}
	}
	return nil
}

func createBlockDeviceDir(buildDeviceDir string, srcDeviceDir string) error {
	// Populate the supplied directory (in our build filesystem) with all the
	// appropriate information pseudofile contents for the block device.
	devName := filepath.Base(srcDeviceDir)
	devFiles, err := ioutil.ReadDir(srcDeviceDir)
	if err != nil {
		return err
	}
	for _, f := range devFiles {
		fname := f.Name()
		fp := filepath.Join(srcDeviceDir, fname)
		fi, err := os.Lstat(fp)
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			// Ignore any symlinks in the deviceDir since they simply point to
			// either self-referential links or information we aren't
			// interested in like "subsystem"
			continue
		} else if fi.IsDir() {
			if strings.HasPrefix(fname, devName) {
				// We're interested in are the directories that begin with the
				// block device name. These are directories with information
				// about the partitions on the device
				buildPartitionDir := filepath.Join(
					buildDeviceDir, fname,
				)
				srcPartitionDir := filepath.Join(
					srcDeviceDir, fname,
				)
				trace("creating partition directory %s\n", buildPartitionDir)
				err = os.MkdirAll(buildPartitionDir, os.ModePerm)
				if err != nil {
					return err
				}
				err = createPartitionDir(buildPartitionDir, srcPartitionDir)
				if err != nil {
					return err
				}
			}
		} else if fi.Mode().IsRegular() {
			// Regular files in the block device directory are both regular and
			// pseudofiles containing information such as the size (in sectors)
			// and whether the device is read-only
			buf, err := ioutil.ReadFile(fp)
			if err != nil {
				return err
			}
			targetPath := filepath.Join(buildDeviceDir, fname)
			trace("creating %s\n", targetPath)
			f, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err = f.Write(buf); err != nil {
				return err
			}
			f.Close()
		}
	}
	// There is a special file $DEVICE_DIR/queue/rotational that, for some hard
	// drives, contains a 1 or 0 indicating whether the device is a spinning
	// disk or not
	srcQueueDir := filepath.Join(
		srcDeviceDir,
		"queue",
	)
	buildQueueDir := filepath.Join(
		buildDeviceDir,
		"queue",
	)
	err = os.MkdirAll(buildQueueDir, os.ModePerm)
	if err != nil {
		return err
	}
	fp := filepath.Join(srcQueueDir, "rotational")
	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(buildQueueDir, "rotational")
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

func createPartitionDir(buildPartitionDir string, srcPartitionDir string) error {
	// Populate the supplied directory (in our build filesystem) with all the
	// appropriate information pseudofile contents for the partition.
	partFiles, err := ioutil.ReadDir(srcPartitionDir)
	if err != nil {
		return err
	}
	for _, f := range partFiles {
		fname := f.Name()
		fp := filepath.Join(srcPartitionDir, fname)
		fi, err := os.Lstat(fp)
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			// Ignore any symlinks in the partition directory since they simply
			// point to information we aren't interested in like "subsystem"
			continue
		} else if fi.IsDir() {
			// The subdirectories in the partition directory are not
			// interesting for us. They have information about power events and
			// traces
			continue
		} else if fi.Mode().IsRegular() {
			// Regular files in the block device directory are both regular and
			// pseudofiles containing information such as the size (in sectors)
			// and whether the device is read-only
			buf, err := ioutil.ReadFile(fp)
			if err != nil {
				return err
			}
			targetPath := filepath.Join(buildPartitionDir, fname)
			trace("creating %s\n", targetPath)
			f, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err = f.Write(buf); err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}

func doSnapshot(buildDir string) error {
	if outPath == "" {
		outPath = defaultOutPath()
		trace("using default output filepath %s\n", outPath)
	}

	var f *os.File
	var err error

	if _, err = os.Stat(outPath); os.IsNotExist(err) {
		if f, err = os.Create(outPath); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		f, err := os.OpenFile(outPath, os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		fs, err := f.Stat()
		if err != nil {
			return err
		}
		if fs.Size() > 0 {
			return fmt.Errorf("File %s already exists and is of size >0", outPath)
		}
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return createSnapshot(tw, buildDir)
}

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&outPath,
		"out", "o",
		outPath,
		"Path to place snapshot. Defaults to file in current directory with name $OS-$ARCH-$HASHSYSTEMNAME.tar.gz",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&debug, "debug", "d", false, "Enable or disable debug mode",
	)
}

func createSnapshot(tw *tar.Writer, buildDir string) error {
	return filepath.Walk(buildDir, func(path string, fi os.FileInfo, err error) error {
		if path == buildDir {
			return nil
		}
		var link string

		if fi.Mode()&os.ModeSymlink != 0 {
			trace("processing symlink %s\n", path)
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		hdr.Name = strings.TrimPrefix(strings.TrimPrefix(path, buildDir), string(os.PathSeparator))

		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err = io.Copy(tw, f); err != nil {
				return err
			}
			f.Close()
		}
		return nil
	})
}
