//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jaypipes/ghw"
	ghwcontext "github.com/jaypipes/ghw/pkg/context"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads a new ghw snapshot",
	RunE:  doRead,
}

// doRead reads a ghw snapshot from the input snapshot path argument
func doRead(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("supply a single argument with the filepath to the snapshot you wish to read")
	}
	inPath := args[0]
	if _, err := os.Stat(inPath); err != nil {
		return err
	}
	os.Setenv("GHW_SNAPSHOT_PATH", inPath)
	ctx := ghwcontext.New()

	return ctx.Do(func() error {
		info, err := ghw.Host()
		fmt.Println(info.String())
		return err
	})
}

func init() {
	rootCmd.AddCommand(readCmd)
}
