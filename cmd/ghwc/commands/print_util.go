//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import "fmt"

type formattable interface {
	String() string
	JSONString(bool) string
	YAMLString() string
}

func printInfo(f formattable) {
	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%s\n", f)
	case outputFormatJSON:
		fmt.Printf("%s\n", f.JSONString(pretty))
	case outputFormatYAML:
		fmt.Printf("%s", f.YAMLString())
	}
}
