//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const debugHeader = `
Date: %s
Build: %s
Version: %s
Git Hash: %s
`

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of gofile",
	Run: func(cmd *cobra.Command, args []string) {
		goVersion := fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Printf(debugHeader, buildDate, goVersion, version, buildHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
