//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// topologyCmd represents the install command
var topologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Show topology information for the host system",
	RunE:  showTopology,
}

// showTopology show topology information for the host system.
func showTopology(cmd *cobra.Command, args []string) error {
	topology := info.Topology
	fmt.Printf("%v\n", topology)

	for _, node := range topology.Nodes {
		fmt.Printf(" %v\n", node)
		for _, cache := range node.Caches {
			fmt.Printf("  %v\n", cache)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(topologyCmd)
}
