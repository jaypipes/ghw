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

// chassisCmd represents the install command
var chassisCmd = &cobra.Command{
	Use:   "chassis",
	Short: "Show chassis information for the host system",
	RunE:  showChassis,
}

// showChassis shows chassis information for the host system.
func showChassis(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	chassis, err := ghw.Chassis(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting chassis info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", chassis)
	case outputFormatJSON:
		fmt.Printf("%s\n", chassis.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", chassis.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(chassisCmd)
}
