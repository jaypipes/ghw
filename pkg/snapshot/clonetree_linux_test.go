//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package snapshot_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/jaypipes/ghw/internal/config"
	"github.com/jaypipes/ghw/pkg/linuxdt"
	"github.com/jaypipes/ghw/pkg/snapshot"
	"github.com/jaypipes/ghw/pkg/util"
)

// NOTE: we intentionally use `os.RemoveAll` - not `snapshot.Cleanup` because we
// want to make sure we never leak directories. `snapshot.Cleanup` is used and
// tested explicitly in `unpack_test.go`.

// nolint: gocyclo
func TestCloneTree(t *testing.T) {
	root, err := snapshot.Unpack(testDataSnapshot)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	cloneRoot, err := os.MkdirTemp("", "ghw-test-clonetree-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cloneRoot)

	fileSpecs := []string{
		filepath.Join(root, "ghw-test-*"),
		filepath.Join(root, "different/subtree/ghw*"),
		filepath.Join(root, "nested/ghw-test*"),
		filepath.Join(root, "nested/tree/of/subdirectories/forming/deep/unbalanced/tree/ghw-test-3"),
	}
	ctx := context.TODO()
	err = snapshot.CopyFilesInto(ctx, fileSpecs, cloneRoot, nil)
	if err != nil {
		t.Fatal(err)
	}

	origContent, err := scanTree(root, "", []string{""})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(origContent)

	cloneContent, err := scanTree(cloneRoot, cloneRoot, []string{"", "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(cloneContent)

	if len(origContent) != len(cloneContent) {
		t.Fatalf("Expected tree size %d got %d", len(origContent), len(cloneContent))
	}
	if !reflect.DeepEqual(origContent, cloneContent) {
		t.Fatalf("subtree content different expected %#v got %#v", origContent, cloneContent)
	}
}

// nolint: gocyclo
func TestCloneSystemTree(t *testing.T) {
	// ok, this is tricky. Validating a cloned tree is a complex business.
	// We do the bare minimum here to check that both the CloneTree and the ValidateClonedTree did something
	// sensible. To really do a meaningful test we need a more advanced functional test, starting with from
	// a ghw snapshot.
	cloneRoot, err := os.MkdirTemp("", "ghw-test-clonetree-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cloneRoot)

	ctx := context.TODO()
	err = snapshot.CloneTreeInto(ctx, cloneRoot)
	if err != nil {
		t.Fatal(err)
	}

	cloneContent, err := snapshot.ExpectedCloneContent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	missing, err := snapshot.ValidateClonedTree(cloneContent, cloneRoot)
	if err != nil {
		t.Fatal(err)
	}

	if len(missing) > 0 && areEntriesOnSysfs(missing) {
		t.Fatalf("Expected content %#v missing into the cloned tree %q", missing, cloneRoot)
	}
}

func areEntriesOnSysfs(sysfsEntries []string) bool {
	// turns out some ISA bridges do not actually expose the driver entry. The reason is not clear.
	// So let's check if we actually have the entry we were looking for on sysfs. If so, we
	// actually failed to clone an entry, and we must fail the test. Otherwise we carry on.
	for _, sysfsEntry := range sysfsEntries {
		if _, err := os.Lstat(sysfsEntry); err == nil {
			return true
		}
	}
	return false
}

func scanTree(root, prefix string, excludeList []string) ([]string, error) {
	var contents []string
	return contents, filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fp := strings.TrimPrefix(path, prefix); !includedInto(fp, excludeList) {
			contents = append(contents, fp)
		}
		return nil
	})
}

func includedInto(s string, items []string) bool {
	if items == nil {
		return false
	}
	for _, item := range items {
		if s == item {
			return true
		}
	}
	return false
}

// A snapshot of a DeviceTree board carries the properties ghw reads from the
// DeviceTree, all but the board's serial number.
func TestCloneDeviceTree(t *testing.T) {
	host := t.TempDir()
	base := filepath.Join(host, "sys", "firmware", "devicetree", "base")
	for name, value := range map[string]string{
		"compatible":            "nvidia,p3737-0000+p3701-0000\x00nvidia,tegra234\x00",
		"model":                 "Jetson AGX Orin Developer Kit\x00",
		"chassis-type":          "embedded\x00",
		"serial-number":         "0123456789\x00",
		"chosen/u-boot,version": "2023.04\x00",
	} {
		path := filepath.Join(base, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var specs []string
	for _, spec := range snapshot.ExpectedCloneStaticContent() {
		specs = append(specs, filepath.Join(host, spec))
	}
	clone := t.TempDir()
	if err := snapshot.CopyFilesInto(context.Background(), specs, clone, nil); err != nil {
		t.Fatal(err)
	}
	read := func(root string) []string {
		ctx := config.ContextFromArgs(config.WithChroot(root))
		return []string{
			linuxdt.Model(ctx), linuxdt.Vendor(ctx), linuxdt.SoC(ctx),
			linuxdt.ChassisType(ctx), linuxdt.UBootVersion(ctx),
		}
	}
	if got, want := read(filepath.Join(clone, host)), read(host); !reflect.DeepEqual(got, want) {
		t.Errorf("from the snapshot = %q, want %q", got, want)
	}
	ctx := config.ContextFromArgs(config.WithChroot(filepath.Join(clone, host)))
	if serial := linuxdt.SerialNumber(ctx); serial != util.UNKNOWN {
		t.Errorf("snapshot carries the serial number %q", serial)
	}
}
