//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"os"
	"testing"
)

func TestPCI(t *testing.T) {
	if _, ok := os.LookupEnv("GHW_TESTING_SKIP_PCI"); ok {
		t.Skip("Skipping PCI tests.")
	}
	info, err := PCI()
	if err != nil {
		t.Fatalf("Expected no error creating PciInfo, but got %v", err)
	}

	if len(info.Classes) == 0 {
		t.Fatalf("Expected >0 PCI classes, but found 0.")
	}

	wirelessController, exists := info.Classes["0d"]
	if !exists {
		t.Fatalf("Expected to find wireless controller class in hash for identifier '0d'")
	}
	if wirelessController.Name != "Wireless controller" {
		t.Fatalf("Expected wireless controller class name to be 'Wireless controller' but got '%v'", wirelessController.Name)
	}
}
