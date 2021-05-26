//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package context_test

import (
	"os"
	"testing"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/option"
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

// nolint: gocyclo
func TestContextReadiness(t *testing.T) {
	ctx := context.New()
	if ctx.IsReady() {
		t.Fatalf("context ready before Setup()")
	}

	ctx.Do(func() error {
		if !ctx.IsReady() {
			t.Fatalf("context NOT ready inside Do()")
		}
		return nil
	})

	if ctx.IsReady() {
		t.Fatalf("context ready after Teardown()")
	}
}

// nolint: gocyclo
func TestContextReadinessNested(t *testing.T) {
	ctx := context.New()
	if ctx.IsReady() {
		t.Fatalf("context ready before Setup()")
	}

	ctx.Do(func() error {
		if !ctx.IsReady() {
			t.Fatalf("context NOT ready inside outer Do()")
		}
		ctx.Do(func() error {
			if !ctx.IsReady() {
				t.Fatalf("context NOT ready inside inner Do()")
			}
			return nil
		})
		if !ctx.IsReady() {
			t.Fatalf("context NOT ready after inner Do()")
		}
		return nil
	})

	if ctx.IsReady() {
		t.Fatalf("context ready after Teardown() - refcount = %d", ctx.RefCount())
	}
}

// nolint: gocyclo
func TestContextReadinessDeeplyNested(t *testing.T) {
	// we don't expect more nesting than this atm
	ctx := context.New()
	if ctx.IsReady() {
		t.Fatalf("context ready before Setup()")
	}
	if ctx.RefCount() != 0 {
		t.Fatalf("context refcount unexpected value %d", ctx.RefCount())
	}

	ctx.Do(func() error {
		if !ctx.IsReady() {
			t.Fatalf("context NOT ready inside outer Do()")
		}
		ctx.Do(func() error {
			if !ctx.IsReady() {
				t.Fatalf("context NOT ready inside middle Do()")
			}
			ctx.Do(func() error {
				if !ctx.IsReady() {
					t.Fatalf("context NOT ready inside inner Do()")
				}
				return nil
			})
			if !ctx.IsReady() {
				t.Fatalf("context NOT ready after inner Do()")
			}
			return nil
		})
		if !ctx.IsReady() {
			t.Fatalf("context NOT ready after middle Do()")
		}
		return nil
	})

	if ctx.IsReady() {
		t.Fatalf("context ready after Teardown() - refcount = %d", ctx.RefCount())
	}
	if ctx.RefCount() != 0 {
		t.Fatalf("context refcount unexpected value %d", ctx.RefCount())
	}
}
