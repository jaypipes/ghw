//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

type PCIClassInfo struct {
	Id         string          // hex-encoded PCI_ID for the device class
	Name       string          // common string name for the class
	Subclasses []*PCIClassInfo // Any subclasses belonging to this class
}

// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type PCIProductInfo struct {
	VendorId   string            // vendor ID for the product
	Id         string            // hex-encoded PCI_ID for the product/model
	Name       string            // common string name of the vendor
	Subsystems []*PCIProductInfo // "subdevices" or "subsystems" for the product
}

type PCIVendorInfo struct {
	Id       string            // hex-encoded PCI_ID for the vendor
	Name     string            // common string name of the vendor
	Products []*PCIProductInfo // all top-level devices for the vendor
}

type PCIDeviceInfo struct {
	Vendor           PCIVendorInfo
	SubsystemVendor  PCIVendorInfo // optional subvendor information
	Product          PCIProductInfo
	SubsystemProduct PCIProductInfo // optional sub-device information
	Class            PCIClassInfo
	Subclass         PCIClassInfo // optional sub-class for the device
}

type PCIInfo struct {
	// hash of class ID -> class information
	Classes map[string]*PCIClassInfo
	// hash of vendor ID -> vendor information
	Vendors map[string]*PCIVendorInfo
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*PCIProductInfo
}

func (db *PCIInfo) GetDeviceInfo(address string) *PCIDeviceInfo {
	return nil
}

func PCI() (*PCIInfo, error) {
	info := &PCIInfo{}
	err := pciFillInfo(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
