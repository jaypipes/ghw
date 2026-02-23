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

// pciCmd represents the install command
var pciCmd = &cobra.Command{
	Use:   "pci",
	Short: "Show information about PCI devices on the host system",
	RunE:  showPCI,
}

// showPCI shows information for PCI devices on the host system.
func showPCI(cmd *cobra.Command, args []string) error {
	pci, err := ghw.PCI(cmd.Context())
	if err != nil {
		return fmt.Errorf("error getting PCI info: %w", err)
	}

	printInfo(pci)
	return nil
}

func init() {
	rootCmd.AddCommand(pciCmd)
}
