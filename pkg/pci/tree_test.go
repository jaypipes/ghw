//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci_test

import (
	"reflect"
	"testing"

	"github.com/jaypipes/ghw/pkg/pci"
)

// buildTree assembles a small device tree shaped like:
//
//	root
//	 ├─ a
//	 │   ├─ aa
//	 │   └─ ab
//	 └─ b
//
// Each device is identified by its Address so the tests can assert
// over a deterministic visit order.
func buildTree() (root, a, aa, ab, b *pci.Device) {
	root = &pci.Device{Address: "root"}
	a = &pci.Device{Address: "a", Parent: root}
	aa = &pci.Device{Address: "aa", Parent: a}
	ab = &pci.Device{Address: "ab", Parent: a}
	b = &pci.Device{Address: "b", Parent: root}
	a.Children = []*pci.Device{aa, ab}
	root.Children = []*pci.Device{a, b}
	return
}

func TestDeviceWalkPreOrder(t *testing.T) {
	root, _, _, _, _ := buildTree()
	var got []string
	root.Walk(func(d *pci.Device) bool {
		got = append(got, d.Address)
		return true
	})
	want := []string{"root", "a", "aa", "ab", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Walk visit order = %v, want %v", got, want)
	}
}

func TestDeviceWalkPruneSubtree(t *testing.T) {
	root, _, _, _, _ := buildTree()
	var got []string
	root.Walk(func(d *pci.Device) bool {
		got = append(got, d.Address)
		// Skip everything below "a", but still visit a itself.
		return d.Address != "a"
	})
	want := []string{"root", "a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Walk visit order = %v, want %v", got, want)
	}
}

func TestDeviceWalkNilReceiver(t *testing.T) {
	var d *pci.Device
	called := false
	d.Walk(func(*pci.Device) bool {
		called = true
		return true
	})
	if called {
		t.Errorf("Walk on nil receiver invoked fn")
	}
}

func TestDeviceAncestors(t *testing.T) {
	root, _, aa, _, _ := buildTree()

	if got := root.Ancestors(); len(got) != 0 {
		t.Errorf("Ancestors of root = %v, want empty", got)
	}

	got := []string{}
	for _, a := range aa.Ancestors() {
		got = append(got, a.Address)
	}
	want := []string{"a", "root"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Ancestors of aa = %v, want %v", got, want)
	}
}

func TestDeviceRoot(t *testing.T) {
	root, _, aa, _, b := buildTree()

	if r := root.Root(); r != root {
		t.Errorf("root.Root() = %v, want self", r)
	}
	if r := aa.Root(); r != root {
		t.Errorf("aa.Root() = %v, want root", r)
	}
	if r := b.Root(); r != root {
		t.Errorf("b.Root() = %v, want root", r)
	}

	var nilDev *pci.Device
	if r := nilDev.Root(); r != nil {
		t.Errorf("(nil).Root() = %v, want nil", r)
	}
}
