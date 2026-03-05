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
	"github.com/stretchr/testify/require"
)

// NOTE: we intentionally use `os.RemoveAll` - not `snapshot.Cleanup` because we
// want to make sure we never leak directories. `snapshot.Cleanup` is used and
// tested explicitly in `unpack_test.go`.

// nolint: gocyclo
func TestPackUnpack(t *testing.T) {
	require := require.New(t)
	root, err := snapshot.Unpack(testDataSnapshot)
	require.Nil(err)
	defer os.RemoveAll(root)

	tmpfile, err := os.CreateTemp("", "ght-test-snapshot-*.tgz")
	require.Nil(err)
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	ctx := context.TODO()
	err = snapshot.PackWithWriter(ctx, tmpfile, root)
	require.Nil(err)
	err = tmpfile.Close()
	require.Nil(err)

	cloneRoot, err := snapshot.Unpack(tmpfile.Name())
	require.Nil(err)
	defer os.RemoveAll(cloneRoot)

	verifyTestData(t, cloneRoot)
}
