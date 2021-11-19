// +build linux
//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

var (
	// show debug output
	debug = false
	// output filepath to save snapshot to
	outPath string
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

	snapshot.SetTraceFunction(trace)
	if err = snapshot.CloneTreeInto(scratchDir); err != nil {
		return err
	}

	if outPath == "" {
		outPath = defaultOutPath()
		trace("using default output filepath %s\n", outPath)
	}

	return snapshot.PackFrom(outPath, scratchDir)
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
