//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

//go:build darwin
// +build darwin

package block

import "testing"

func TestStorageControllerFromPlist(t *testing.T) {
	tests := []struct {
		name string
		info *diskUtilInfoPlist
		want StorageController
	}{
		{
			name: "standard IONVMeController (Intel Mac / Thunderbolt NVMe)",
			info: &diskUtilInfoPlist{
				DeviceTreePath: "IODeviceTree:/PCI0@0/RP17@1B/IONVMeController",
			},
			want: StorageControllerNVMe,
		},
		{
			name: "Apple Silicon internal NVMe via Apple Fabric bus protocol",
			info: &diskUtilInfoPlist{
				DeviceTreePath: "IODeviceTree:/arm-io@10F00000/ans@77400000/iop-ans-nub/AppleANS3NVMeController",
				BusProtocol:    "Apple Fabric",
			},
			want: StorageControllerNVMe,
		},
		{
			name: "T2 Mac internal NVMe via Apple Fabric bus protocol, no NVMe substring in path",
			info: &diskUtilInfoPlist{
				DeviceTreePath: "IODeviceTree:/PCI0@0/RP17@1B/ANS2@0/AppleANS2Controller",
				BusProtocol:    "Apple Fabric",
			},
			want: StorageControllerNVMe,
		},
		{
			name: "SATA/USB disk",
			info: &diskUtilInfoPlist{
				DeviceTreePath: "IODeviceTree:/PCI0@0/SATA@0",
				BusProtocol:    "SATA",
			},
			want: StorageControllerSCSI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := storageControllerFromPlist(tt.info)
			if got != tt.want {
				t.Errorf("storageControllerFromPlist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindEntryID(t *testing.T) {
	// Two storage devices sharing the same node name, as seen under macOS
	// virtualization (AppleVirtIOStorageDevice), to make sure they resolve
	// to different, correct IDs rather than ambiguously matching either one.
	tree := &ioregTreeNode{
		Name: "Root",
		Children: []ioregTreeNode{
			{
				Name: "device-tree",
				Children: []ioregTreeNode{
					{
						Name:     "arm-io",
						Location: "10F00000",
						Children: []ioregTreeNode{
							{
								Name:     "pcie",
								Location: "30000000",
								Children: []ioregTreeNode{
									{
										Name:     "pci106b,1a00",
										Location: "4",
										Children: []ioregTreeNode{
											{EntryID: 100, Name: "AppleVirtIOStorageDevice"},
										},
									},
									{
										Name:     "pci106b,1a00",
										Location: "5",
										Children: []ioregTreeNode{
											{EntryID: 200, Name: "AppleVirtIOStorageDevice"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		path    string
		want    int64
		wantErr bool
	}{
		{
			name: "first of two same-named devices",
			path: "IODeviceTree:/arm-io@10F00000/pcie@30000000/pci106b,1a00@4/AppleVirtIOStorageDevice",
			want: 100,
		},
		{
			name: "second of two same-named devices",
			path: "IODeviceTree:/arm-io@10F00000/pcie@30000000/pci106b,1a00@5/AppleVirtIOStorageDevice",
			want: 200,
		},
		{
			name:    "nonexistent segment",
			path:    "IODeviceTree:/arm-io@10F00000/pcie@30000000/does-not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findEntryID(tree, tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("findEntryID() expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("findEntryID() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("findEntryID() = %d, want %d", got, tt.want)
			}
		})
	}
}
