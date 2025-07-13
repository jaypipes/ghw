//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/snapshot"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type optKey string

const (
	optsKey optKey = "ghwc.opts"
)

const (
	outputFormatHuman = "human"
	outputFormatJSON  = "json"
	outputFormatYAML  = "yaml"
	usageOutputFormat = `Output format.
Choices are 'json','yaml', and 'human'.`
	usageSnapshotPath = `Specify path to snapshot.`
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
	snapshotPath string
	pretty       bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghwc",
	Short: "ghwc - Discover hardware information.",
	Args:  validateRootCommand,
	Long: `          __
 .-----. |  |--. .--.--.--.
 |  _  | |     | |  |  |  |
 |___  | |__|__| |________|
 |_____|

Discover hardware information.

https://github.com/jaypipes/ghw
`,
	PersistentPreRunE: doPreRun,
	RunE:              showAll,
	SilenceUsage:      true,
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
	rootCmd.PersistentFlags().StringVarP(
		&snapshotPath,
		"snapshot-path", "s",
		"",
		usageSnapshotPath,
	)
}

func showAll(cmd *cobra.Command, args []string) error {

	switch outputFormat {
	case outputFormatHuman:
		for _, f := range []func(*cobra.Command, []string) error{
			showBlock,
			showCPU,
			showGPU,
			showMemory,
			showNetwork,
			showTopology,
			showChassis,
			showBIOS,
			showBaseboard,
			showProduct,
			showAccelerator,
			showUSB,
		} {
			err := f(cmd, args)
			if err != nil {
				return err
			}

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

func validateSnapshotPath() error {
	if _, err := os.Stat(snapshotPath); err != nil {
		return fmt.Errorf("invalid snapshot path: %w", err)
	}
	return nil
}

// validateRootCommand ensures any CLI options or arguments are valid,
// returning an error if not
func validateRootCommand(rootCmd *cobra.Command, args []string) error {
	if !haveValidOutputFormat() {
		return fmt.Errorf("invalid output format %q", outputFormat)
	}
	if snapshotPath != "" {
		if err := validateSnapshotPath(); err != nil {
			return err
		}
	}
	return nil
}

func doPreRun(cmd *cobra.Command, args []string) error {
	opts := []option.Option{}
	if snapshotPath != "" {
		// unpack the snapshot into a tempdir and clean up this tempdir after
		// the run...
		unpackDir, err := os.MkdirTemp("", "ghw-snap-*")
		if err != nil {
			return err
		}
		err = snapshot.UnpackInto(snapshotPath, unpackDir)
		if err != nil {
			return err
		}
		opts = append(opts, ghw.WithChroot(unpackDir))
		cmd.PersistentPostRunE = func(c *cobra.Command, args []string) error {
			_ = os.RemoveAll(unpackDir)
			return nil
		}
	}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, optsKey, opts)
	cmd.SetContext(ctx)
	return nil
}
