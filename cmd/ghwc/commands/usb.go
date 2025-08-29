//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// netCmd represents the install command
var usbCmd = &cobra.Command{
	Use:   "usb",
	Short: "Show USB information for the host system",
	RunE:  showUSB,
}

// showNetwork show network information for the host system.
func showUSB(cmd *cobra.Command, args []string) error {
	usb, err := ghw.USB()
	if err != nil {
		return errors.Wrap(err, "error getting network info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", usb)
		for _, usb := range usb.USBs {
			fmt.Printf(" %+v\n", usb)
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", usb.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", usb.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(usbCmd)
}
