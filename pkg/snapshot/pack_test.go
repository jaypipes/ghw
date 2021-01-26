//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

// nolint: gocyclo
func TestPackUnpack(t *testing.T) {
	root, err := snapshot.Unpack(testDataSnapshot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	defer os.RemoveAll(root)

	tmpfile, err := ioutil.TempFile("", "ght-test-snapshot-*.tgz")
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	err = snapshot.PackWithWriter(tmpfile, root)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	cloneRoot, err := snapshot.Unpack(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	verifyTestData(t, cloneRoot)
}
