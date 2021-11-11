//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/adumandix/ghw/pkg/snapshot"
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

	verifyTestData(t, root)

	err = snapshot.Cleanup(root)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Expected %q to be gone, but still exists", root)
	}
}

// nolint: gocyclo
func TestUnpackInto(t *testing.T) {
	testRoot, err := ioutil.TempDir("", "ghw-test-snapshot-*")
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	_, err = snapshot.UnpackInto(testDataSnapshot, testRoot, 0)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	verifyTestData(t, testRoot)

	// note that in real production code the caller will likely manage its
	// snapshot root directory in a different way, here we call snapshot.Cleanup
	// to clean up after ourselves more than to test it
	err = snapshot.Cleanup(testRoot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	if _, err := os.Stat(testRoot); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Expected %q to be gone, but still exists", testRoot)
	}
}

// nolint: gocyclo
func TestUnpackIntoPresrving(t *testing.T) {
	testRoot, err := ioutil.TempDir("", "ghw-test-snapshot-*")
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(testRoot, "canary"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	_, err = snapshot.UnpackInto(testDataSnapshot, testRoot, snapshot.OwnTargetDirectory)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	entries, err := ioutil.ReadDir(testRoot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("Expected one entry in %q, but got %v", testRoot, entries)
	}

	canary := entries[0]
	if canary.Name() != "canary" {
		t.Fatalf("Expected entry %q, but got %q", "canary", canary.Name())
	}

	// note that in real production code the caller will likely manage its
	// snapshot root directory in a different way, here we call snapshot.Cleanup
	// to clean up after ourselves more than to test it
	err = snapshot.Cleanup(testRoot)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}

	if _, err := os.Stat(testRoot); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Expected %q to be gone, but still exists", testRoot)
	}
}

func verifyTestData(t *testing.T, root string) {
	verifyFileContent(t, filepath.Join(root, "ghw-test-0"), "ghw-test-0\n")
	verifyFileContent(t, filepath.Join(root, "ghw-test-1"), "ghw-test-1\n")
	verifyFileContent(t, filepath.Join(root, "nested", "ghw-test-2"), "ghw-test-2\n")

}

func verifyFileContent(t *testing.T, path, expected string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	content := string(data)
	if content != expected {
		t.Fatalf("Expected %q, but got %q", expected, content)
	}
}
