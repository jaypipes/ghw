//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import "fmt"

func printInfo(f interface{}) {
	switch outputFormat {
	case outputFormatHuman:
		fmt.Printf("%s\n", f)
	case outputFormatJSON:
		fmt.Printf("%s\n", JSONString(f, pretty))
	case outputFormatYAML:
		fmt.Printf("%s", YAMLString(f))
	}
}
