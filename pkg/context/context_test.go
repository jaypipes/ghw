//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context_test

import (
	"os"
	"testing"

	"github.com/adumandix/ghw/pkg/context"
	"github.com/adumandix/ghw/pkg/option"
)

const (
	testDataSnapshot = "../snapshot/testdata.tar.gz"
)

// nolint: gocyclo
func TestSnapshotContext(t *testing.T) {
	ctx := context.New(option.WithSnapshot(option.SnapshotOptions{
		Path: testDataSnapshot,
	}))

	var uncompressedDir string
	err := ctx.Do(func() error {
		uncompressedDir = ctx.Chroot
		return nil
	})

	if uncompressedDir == "" {
		t.Fatalf("Expected the uncompressed dir path to not be empty")
	}
	if err != nil {
		t.Fatalf("Expected nil err, but got %v", err)
	}
	if _, err = os.Stat(uncompressedDir); !os.IsNotExist(err) {
		t.Fatalf("Expected the uncompressed dir to be deleted: %s", uncompressedDir)
	}
}
