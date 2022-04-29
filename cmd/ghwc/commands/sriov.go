//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"github.com/jaypipes/ghw"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// sriovCmd represents the install command
var sriovCmd = &cobra.Command{
	Use:   "sriov",
	Short: "Show Single Root I/O Virtualization device information for the host system",
	RunE:  showSRIOV,
}

// showSRIOV shows SRIOV information for the host system.
func showSRIOV(cmd *cobra.Command, args []string) error {
	info, err := ghw.PCI()
	if err != nil {
		return errors.Wrap(err, "error getting SRIOV info through PCI")
	}

	printInfo(info.DescribeDevices(info.GetSRIOVDevices()))
	return nil
}

func init() {
	rootCmd.AddCommand(sriovCmd)
}
