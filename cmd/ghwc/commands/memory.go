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

// memoryCmd represents the install command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Show memory information for the host system",
	RunE:  showMemory,
}

// showMemory show memory information for the host system.
func showMemory(cmd *cobra.Command, args []string) error {
	mem := info.Memory
	fmt.Printf("%v\n", mem)
	return nil
}

func init() {
	rootCmd.AddCommand(memoryCmd)
}
