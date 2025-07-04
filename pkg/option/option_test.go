//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package option_test

import (
	"runtime"
	"testing"

	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/pcidb"
)

type optTestCase struct {
	name   string
	opts   []*option.Option
	merged *option.Option
}

// nolint: gocyclo
func TestOption(t *testing.T) {
	var pciTest *optTestCase

	if runtime.GOOS == "linux" {
		pcidb, err := pcidb.New()
		if err != nil {
			t.Fatalf("error creating new pcidb: %v", err)
		}
		pciTest = &optTestCase{
			name: "pcidb",
			opts: []*option.Option{
				option.WithPCIDB(pcidb),
				option.WithChroot("/my/chroot/dir"),
			},
			merged: &option.Option{
				Chroot: stringPtr("/my/chroot/dir"),
				PCIDB:  pcidb,
			},
		}
	}

	optTCases := []optTestCase{
		{
			name: "multiple chroots",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithChroot("/my/chroot/dir/2"),
			},
			merged: &option.Option{
				Chroot:      stringPtr("/my/chroot/dir/2"),
				EnableTools: boolPtr(true),
			},
		},
		{
			name: "multiple chroots interleaved",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithSnapshot(option.SnapshotOptions{
					Path: "/my/snapshot/dir",
				}),
				option.WithChroot("/my/chroot/dir/2"),
			},
			merged: &option.Option{
				// latest seen takes precedence
				Chroot: stringPtr("/my/chroot/dir/2"),
				Snapshot: &option.SnapshotOptions{
					Path: "/my/snapshot/dir",
				},
			},
		},
		{
			name: "multiple snapshots overriding path",
			opts: []*option.Option{
				option.WithSnapshot(option.SnapshotOptions{
					Path: "/my/snapshot/dir",
				}),
				option.WithSnapshot(option.SnapshotOptions{
					Exclusive: true,
				}),
			},
			merged: &option.Option{
				Chroot: stringPtr("/"),
				Snapshot: &option.SnapshotOptions{
					// note Path is gone because the second instance
					// has default (empty) path, and the latest always
					// takes precedence.
					Path:      "",
					Exclusive: true,
				},
			},
		},
		{
			name: "chroot and snapshot",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithSnapshot(option.SnapshotOptions{
					Path:      "/my/snapshot/dir",
					Exclusive: true,
				}),
			},
			merged: &option.Option{
				Chroot: stringPtr("/my/chroot/dir"),
				Snapshot: &option.SnapshotOptions{
					Path:      "/my/snapshot/dir",
					Exclusive: true,
				},
			},
		},
		{
			name: "chroot and snapshot with root",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithSnapshot(option.SnapshotOptions{
					Path:      "/my/snapshot/dir",
					Root:      stringPtr("/my/overridden/chroot/dir"),
					Exclusive: true,
				}),
			},
			merged: &option.Option{
				// caveat! the option merge logic DOES NOT DO the override
				Chroot: stringPtr("/my/chroot/dir"),
				Snapshot: &option.SnapshotOptions{
					Path:      "/my/snapshot/dir",
					Root:      stringPtr("/my/overridden/chroot/dir"),
					Exclusive: true,
				},
			},
		},
		{
			name: "chroot and disabling tools",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithDisableTools(),
			},
			merged: &option.Option{
				Chroot:      stringPtr("/my/chroot/dir"),
				EnableTools: boolPtr(false),
			},
		},
		{
			name: "paths",
			opts: []*option.Option{
				option.WithPathOverrides(option.PathOverrides{
					"/run": "/host-run",
					"/var": "/host-var",
				}),
			},
			merged: &option.Option{
				PathOverrides: option.PathOverrides{
					"/run": "/host-run",
					"/var": "/host-var",
				},
			},
		},
		{
			name: "chroot paths",
			opts: []*option.Option{
				option.WithChroot("/my/chroot/dir"),
				option.WithPathOverrides(option.PathOverrides{
					"/run": "/host-run",
					"/var": "/host-var",
				}),
			},
			merged: &option.Option{
				Chroot: stringPtr("/my/chroot/dir"),
				PathOverrides: option.PathOverrides{
					"/run": "/host-run",
					"/var": "/host-var",
				},
			},
		},
	}
	if pciTest != nil {
		optTCases = append(optTCases, *pciTest)
	}
	for _, optTCase := range optTCases {
		t.Run(optTCase.name, func(t *testing.T) {
			opt := option.Merge(optTCase.opts...)
			if what, ok := optionEqual(optTCase.merged, opt); !ok {
				t.Errorf("expected %#v got %#v - difference on %s", optTCase.merged, opt, what)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func optionEqual(a, b *option.Option) (string, bool) {
	if a == nil || b == nil {
		return "top-level", false
	}
	if a.Chroot != nil {
		if b.Chroot == nil {
			return "chroot ptr", false
		}
		if *a.Chroot != *b.Chroot {
			return "chroot value", false
		}
	}
	if a.Snapshot != nil {
		if b.Snapshot == nil {
			return "snapshot ptr", false
		}
		return optionSnapshotEqual(a.Snapshot, b.Snapshot)
	}
	if a.EnableTools != nil {
		if b.EnableTools == nil {
			return "enabletools ptr", false
		}
		if *a.EnableTools != *b.EnableTools {
			return "enabletools value", false
		}
	}
	return "", true
}

func optionSnapshotEqual(a, b *option.SnapshotOptions) (string, bool) {
	if a.Path != b.Path {
		return "snapshot path", false
	}
	if a.Exclusive != b.Exclusive {
		return "snapshot exclusive flag", false
	}
	if a.Root != nil {
		if b.Root == nil {
			return "snapshot root ptr", false
		}
		if *a.Root != *b.Root {
			return "snapshot root value", false
		}
	}
	return "", true
}
