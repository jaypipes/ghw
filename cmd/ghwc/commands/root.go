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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	outputFormatHuman = "human"
	outputFormatJSON  = "json"
	outputFormatYAML  = "yaml"
	usageOutputFormat = `Output format.
Choices are 'json','yaml', and 'human'.`
)

var (
	version       string
	buildHash     string
	buildDate     string
	debug         bool
	outputFormat  string
	outputFormats = []string{
		outputFormatHuman,
		outputFormatJSON,
		outputFormatYAML,
	}
	pretty bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghwc",
	Short: "ghwc - Discover hardware information.",
	Args:  validateRootCommand,
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

	switch outputFormat {
	case outputFormatHuman:
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
	case outputFormatJSON:
		host, err := ghw.Host()
		if err != nil {
			return errors.Wrap(err, "error getting host info")
		}
		fmt.Printf("%s\n", host.JSONString(pretty))
	case outputFormatYAML:
		host, err := ghw.Host()
		if err != nil {
			return errors.Wrap(err, "error getting host info")
		}
		fmt.Printf("%s", host.YAMLString())
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute(v string, bh string, bd string) {
	version = v
	buildHash = bh
	buildDate = bd

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func haveValidOutputFormat() bool {
	for _, choice := range outputFormats {
		if choice == outputFormat {
			return true
		}
	}
	return false
}

// validateRootCommand ensures any CLI options or arguments are valid,
// returning an error if not
func validateRootCommand(rootCmd *cobra.Command, args []string) error {
	if !haveValidOutputFormat() {
		return fmt.Errorf("invalid output format %q", outputFormat)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "Enable or disable debug mode",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputFormat,
		"format", "f",
		outputFormatHuman,
		usageOutputFormat,
	)
	rootCmd.PersistentFlags().BoolVar(
		&pretty, "pretty", false, "When outputting JSON, use indentation",
	)
}
