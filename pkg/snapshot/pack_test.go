//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot_test

import (
	"context"
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/snapshot"
)

// NOTE: we intentionally use `os.RemoveAll` - not `snapshot.Cleanup` because we
// want to make sure we never leak directories. `snapshot.Cleanup` is used and
// tested explicitly in `unpack_test.go`.

// nolint: gocyclo
func TestPackUnpack(t *testing.T) {
	root, err := snapshot.Unpack(testDataSnapshot)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	tmpfile, err := os.CreateTemp("", "ght-test-snapshot-*.tgz")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	ctx := context.TODO()
	err = snapshot.PackWithWriter(ctx, tmpfile, root)
	if err != nil {
		t.Fatal(err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	cloneRoot, err := snapshot.Unpack(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cloneRoot)

	verifyTestData(t, cloneRoot)
}
