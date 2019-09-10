//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"testing"
)

func TestDMI(t *testing.T) {
	info, err := DMI()
	if err != nil {
		t.Fatalf("Expected no error creating DMIInfo, but got %v", err)
	}

	// NOTE(mnaser): This isn't the most ideal of tests but those two values
	//               seem to consistently always supply a value in my testing.

	if info.Product.Name == UNKNOWN {
		t.Fatalf("Got unknown product name.")
	}

	if info.Product.Version == UNKNOWN {
		t.Fatalf("Got unknown product version.")
	}
}
