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

// usbCmd represents the `usb` command
var usbCmd = &cobra.Command{
	Use:   "usb",
	Short: "Show USB information for the host system",
	RunE:  showUSB,
}

// showUSB show usb information for the host system.
func showUSB(cmd *cobra.Command, args []string) error {
	usb, err := ghw.USB(cmd.Context())
	if err != nil {
		return fmt.Errorf("error getting USB info: %w", err)
	}

	printInfo(usb)
	return nil
}

func init() {
	rootCmd.AddCommand(usbCmd)
}
