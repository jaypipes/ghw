//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghw-snapshot",
	Short: "ghw-snapshot - create and read ghw snapshots.",
	Long: `
        __                                                   __           __   
.-----.|  |--.--.--.--.______.-----.-----.---.-.-----.-----.|  |--.-----.|  |_ 
|  _  ||     |  |  |  |______|__ --|     |  _  |  _  |__ --||     |  _  ||   _|
|___  ||__|__|________|      |_____|__|__|___._|   __|_____||__|__|_____||____|
|_____|                                        |__|                            

Create and read ghw snapshots.

https://github.com/jaypipes/ghw
`,
	RunE: doCreate,
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func trace(msg string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(msg, args...)
}

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "Enable or disable debug mode",
	)
}
