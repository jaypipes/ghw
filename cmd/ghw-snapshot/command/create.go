//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

var (
	// output filepath to save snapshot to
	outPath string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new ghw snapshot",
	RunE:  doCreate,
}

// doCreate creates a ghw snapshot
func doCreate(cmd *cobra.Command, args []string) error {
	scratchDir, err := os.MkdirTemp("", "ghw-snapshot")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratchDir)

	snapshot.SetTraceFunction(trace)
	if err = snapshot.CloneTreeInto(scratchDir); err != nil {
		return err
	}

	if outPath == "" {
		outPath, err = defaultOutPath()
		if err != nil {
			return err
		}
		trace("using default output filepath %s\n", outPath)
	}

	return snapshot.PackFrom(outPath, scratchDir)
}

func systemFingerprint() (string, error) {
	hn, err := os.Hostname()
	if err != nil {
		return "unknown", err
	}
	m := md5.New()
	_, err = io.WriteString(m, hn)
	if err != nil {
		return "unknown", err
	}
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

func defaultOutPath() (string, error) {
	fp, err := systemFingerprint()
	if err != nil {
		return "unknown", err
	}
	return fmt.Sprintf("%s-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH, fp), nil
}

func init() {
	createCmd.PersistentFlags().StringVarP(
		&outPath,
		"out", "o",
		outPath,
		"Path to place snapshot. Defaults to file in current directory with name $OS-$ARCH-$HASHSYSTEMNAME.tar.gz",
	)
	rootCmd.AddCommand(createCmd)
}
