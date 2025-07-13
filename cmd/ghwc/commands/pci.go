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

// pciCmd represents the install command
var pciCmd = &cobra.Command{
	Use:   "pci",
	Short: "Show information about PCI devices on the host system",
	RunE:  showPCI,
}

// showPCI shows information for PCI devices on the host system.
func showPCI(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	pci, err := ghw.PCI(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting PCI info")
	}

	printInfo(pci)
	return nil
}

func init() {
	rootCmd.AddCommand(pciCmd)
}
