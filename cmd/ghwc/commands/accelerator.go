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

// acceleratorCmd represents the install command
var acceleratorCmd = &cobra.Command{
	Use:   "accelerator",
	Short: "Show processing accelerators information for the host system",
	RunE:  showGPU,
}

// showAccelerator show processing accelerators information for the host system.
func showAccelerator(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	accel, err := ghw.Accelerator(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting Accelerator info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", accel)

		for _, card := range accel.Devices {
			fmt.Printf(" %v\n", card)
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", accel.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", accel.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(acceleratorCmd)
}
