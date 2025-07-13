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

// productCmd represents the install command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Show product information for the host system",
	RunE:  showProduct,
}

// showProduct shows product information for the host system.
func showProduct(cmd *cobra.Command, args []string) error {
	opts := cmd.Context().Value(optsKey).([]ghw.Option)
	product, err := ghw.Product(opts...)
	if err != nil {
		return errors.Wrap(err, "error getting product info")
	}

	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%v\n", product)
	case outputFormatJSON:
		fmt.Printf("%s\n", product.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", product.YAMLString())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(productCmd)
}
