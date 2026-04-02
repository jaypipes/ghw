//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"
)

func printInfo(f interface{}) {
	switch outputFormat {
	case outputFormatJSON:
		f = JSONString(f, pretty)
	case outputFormatYAML:
		f = YAMLString(f)
	}
	fmt.Printf("%s", f)
}
