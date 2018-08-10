//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cmd

import (
	"fmt"

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
	net := info.Network
	fmt.Printf("%v\n", net)

	for _, nic := range net.NICs {
		fmt.Printf(" %v\n", nic)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(netCmd)
}
