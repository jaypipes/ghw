//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"github.com/jaypipes/ghw"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// packagesCmd represents the install command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Show installed packages/software information for the host system",
	RunE:  showPackages,
}

// showMemory show memory information for the host system.
func showPackages(cmd *cobra.Command, args []string) error {
	packages, err := ghw.Packages()
	if err != nil {
		return errors.Wrap(err, "error getting packages info")
	}

	printInfo(packages)
	return nil
}

func init() {
	rootCmd.AddCommand(packagesCmd)
}
