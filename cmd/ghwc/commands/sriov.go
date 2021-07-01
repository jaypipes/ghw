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

// sriovCmd represents the listing command
var sriovCmd = &cobra.Command{
	Use:   "sriov",
	Short: "Show SRIOV devices information for the host system",
	RunE:  showSRIOV,
}

// showSRIOV show SRIOV physical device information for the host system.
func showSRIOV(cmd *cobra.Command, args []string) error {
	sriov, err := ghw.SRIOV()
	if err != nil {
		return errors.Wrap(err, "error getting SRIOV info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", sriov)

		for _, dev := range sriov.PhysicalFunctions {
			fmt.Printf(" %v\n", dev)
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", sriov.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", sriov.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(sriovCmd)
}
