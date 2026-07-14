//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

// tpmCmd represents the `tpm` command
var tpmCmd = &cobra.Command{
	Use:   "tpm",
	Short: "Show TPM information for the host system",
	RunE:  showTPM,
}

// showTPM shows TPM information for the host system.
func showTPM(cmd *cobra.Command, args []string) error {
	tpm, err := ghw.TPM(cmd.Context())
	if err != nil {
		return fmt.Errorf("error getting TPM info: %w", err)
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", tpm)
	case outputFormatJSON:
		fmt.Printf("%s\n", tpm.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", tpm.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(tpmCmd)
}
