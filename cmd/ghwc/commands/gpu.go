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

// gpuCmd represents the install command
var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Show graphics/GPU information for the host system",
	RunE:  showGPU,
}

// showGPU show graphics/GPU information for the host system.
func showGPU(cmd *cobra.Command, args []string) error {
	gpu, err := ghw.GPU()
	if err != nil {
		return errors.Wrap(err, "error getting GPU info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", gpu)

		for _, card := range gpu.GraphicsCards {
			fmt.Printf(" %v\n", card)
		}
	case outputFormatJSON:
		fmt.Printf("%s\n", gpu.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", gpu.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(gpuCmd)
}
