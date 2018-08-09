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

// blockCmd represents the install command
var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Show block storage information for the host system",
	RunE:  showBlock,
}

// showBlock show block storage information for the host system.
func showBlock(cmd *cobra.Command, args []string) error {
	block := info.Block
	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		fmt.Printf(" %v\n", disk)
		for _, part := range disk.Partitions {
			fmt.Printf("  %v\n", part)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(blockCmd)
}
