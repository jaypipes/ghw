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

// watchdogCmd represents the `watchdog` command
var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Show watchdog information for the host system",
	RunE:  showWatchdog,
}

// showWatchdog shows watchdog information for the host system.
func showWatchdog(cmd *cobra.Command, args []string) error {
	watchdog, err := ghw.Watchdog(cmd.Context())
	if err != nil {
		return fmt.Errorf("error getting watchdog info: %w", err)
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", watchdog)
	case outputFormatJSON:
		fmt.Printf("%s\n", watchdog.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", watchdog.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(watchdogCmd)
}
