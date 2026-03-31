//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

// baseboardCmd represents the install command
var baseboardCmd = &cobra.Command{
	Use:   "baseboard",
	Short: "Show baseboard information for the host system",
	RunE:  showBaseboard,
}

// showBaseboard shows baseboard information for the host system.
func showBaseboard(cmd *cobra.Command, args []string) error {
	baseboard, err := ghw.Baseboard(cmd.Context())
	if err != nil {
		return fmt.Errorf("error getting baseboard info: %w", err)
	}

	printInfo(baseboard)
	return nil
}

func init() {
	rootCmd.AddCommand(baseboardCmd)
}
