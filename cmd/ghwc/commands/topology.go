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

// topologyCmd represents the install command
var topologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Show topology information for the host system",
	RunE:  showTopology,
}

// showTopology show topology information for the host system.
func showTopology(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	topology, err := ghw.Topology(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting topology info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", topology)

		for _, node := range topology.Nodes {
			fmt.Printf(" %v\n", node)
			for _, cache := range node.Caches {
				fmt.Printf("  %v\n", cache)
			}
			fmt.Printf("  %v\n", node.Memory)
			fmt.Printf("  distances\n")
			for nodeID, dist := range node.Distances {
				fmt.Printf("    to node #%d %v\n", nodeID, dist)
			}
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", topology.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", topology.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(topologyCmd)
}
