//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"os"
)

type closer interface {
	Close() error
}

func safeClose(c closer) {
	err := c.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to close: %s", err)
	}
}

func warn(msg string, args ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, "WARNING: ")
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
}
