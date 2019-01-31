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
var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Show network information for the host system",
	RunE:  showNetwork,
}

// showNetwork show network information for the host system.
func showNetwork(cmd *cobra.Command, args []string) error {
	net, err := ghw.Network()
	if err != nil {
		return errors.Wrap(err, "error getting network info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", net)

		for _, nic := range net.NICs {
			fmt.Printf(" %v\n", nic)

			enabledCaps := make([]int, 0)
			for x, cap := range nic.Capabilities {
				if cap.IsEnabled {
					enabledCaps = append(enabledCaps, x)
				}
			}
			if len(enabledCaps) > 0 {
				fmt.Printf("  enabled capabilities:\n")
				for _, x := range enabledCaps {
					fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
				}
			}
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", net.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", net.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(netCmd)
}
