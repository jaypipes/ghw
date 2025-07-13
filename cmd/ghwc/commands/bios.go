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

// biosCmd represents the install command
var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "Show BIOS information for the host system",
	RunE:  showBIOS,
}

// showBIOS shows BIOS host system.
func showBIOS(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	bios, err := ghw.BIOS(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting BIOS info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", bios)
	case outputFormatJSON:
		fmt.Printf("%s\n", bios.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", bios.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(biosCmd)
}
