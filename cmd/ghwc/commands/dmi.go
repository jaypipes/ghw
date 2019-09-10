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

// dmiCmd represents the install command
var dmiCmd = &cobra.Command{
	Use:   "dmi",
	Short: "Show DMI information for the host system",
	RunE:  showDMI,
}

// showDMI show DMI information for the host system.
func showDMI(cmd *cobra.Command, args []string) error {
	dmi, err := ghw.DMI()
	if err != nil {
		return errors.Wrap(err, "error getting DMI info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", dmi)
	case outputFormatJSON:
		fmt.Printf("%s\n", dmi.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", dmi.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(dmiCmd)
}
