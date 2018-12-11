//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"
	"os"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

var (
	version   string
	buildHash string
	buildDate string
	debug     bool
	info      *ghw.HostInfo
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghwc",
	Short: "ghwc - Discover hardware information.",
	Long: `
          __
 .-----. |  |--. .--.--.--.
 |  _  | |     | |  |  |  |
 |___  | |__|__| |________|
 |_____|

Discover hardware information.

https://github.com/jaypipes/ghw
`,
	RunE: showAll,
}

func showAll(cmd *cobra.Command, args []string) error {
	if err := showBlock(cmd, args); err != nil {
		return err
	}
	if err := showCPU(cmd, args); err != nil {
		return err
	}
	if err := showGPU(cmd, args); err != nil {
		return err
	}
	if err := showMemory(cmd, args); err != nil {
		return err
	}
	if err := showNetwork(cmd, args); err != nil {
		return err
	}
	if err := showTopology(cmd, args); err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, bh string, bd string) {
	version = v
	buildHash = bh
	buildDate = bd

	i, err := ghw.Host()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info = i

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable or disable debug mode")
}
