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
