//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

const (
	testDataSnapshot = "testdata.tar.gz"
)

// nolint: gocyclo
func TestUnpack(t *testing.T) {
	root, err := snapshot.Unpack(testDataSnapshot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	defer func() {
		_ = os.RemoveAll(root)
	}()

	verifyTestData(t, root)
}

// nolint: gocyclo
func TestUnpackInto(t *testing.T) {
	testRoot, err := os.MkdirTemp("", "ghw-test-snapshot-*")
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	defer func() {
		_ = os.RemoveAll(testRoot)
	}()

	err = snapshot.UnpackInto(testDataSnapshot, testRoot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	verifyTestData(t, testRoot)
}

func verifyTestData(t *testing.T, root string) {
	verifyFileContent(t, filepath.Join(root, "ghw-test-0"), "ghw-test-0\n")
	verifyFileContent(t, filepath.Join(root, "ghw-test-1"), "ghw-test-1\n")
	verifyFileContent(t, filepath.Join(root, "nested", "ghw-test-2"), "ghw-test-2\n")

}

func verifyFileContent(t *testing.T, path, expected string) {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	content := string(data)
	if content != expected {
		t.Fatalf("Expected %q, but got %q", expected, content)
	}
}
