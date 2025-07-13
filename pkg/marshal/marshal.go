//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package marshal

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v3"
)

// SafeYAML returns a string after marshalling the supplied parameter into YAML.
func SafeYAML(p interface{}) string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	var jsonObj interface{}
	if err := yaml.Unmarshal(b, &jsonObj); err != nil {
		return ""
	}

	yb, err := yaml.Marshal(jsonObj)
	if err != nil {
		return ""
	}

	return string(yb)
}

// SafeJSON returns a string after marshalling the supplied parameter into
// JSON. Accepts an optional argument to trigger pretty/indented formatting of
// the JSON string.
func SafeJSON(p interface{}, indent bool) string {
	var b []byte
	var err error
	if !indent {
		b, err = json.Marshal(p)
	} else {
		b, err = json.MarshalIndent(&p, "", "  ")
	}
	if err != nil {
		return ""
	}
	return string(b)
}
