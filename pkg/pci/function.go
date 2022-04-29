//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pci

import "fmt"

// Function describes an SR-IOV physical or virtual function. Physical functions
// will have no Parent Function struct pointer and will have one or more Function
// structs in the Virtual field.
type Function struct {
	Parent *Function `json:"parent,omitempty"`
	// MaxVirtual contains the maximum number of supported virtual
	// functions for this physical function
	MaxVirtual int `json:"max_virtual,omitempty"`
	// Virtual contains the physical function's virtual functions
	Virtual []*Function `json:"virtual_functions"`
}

// IsPhysical returns true if the PCIe function is a physical function, false
// if it is a virtual function. It is safe to assume that if a function is not
// physical, then is virtual (e.g. can't be anything else)
func (f *Function) IsPhysical() bool {
	return f.Parent == nil
}

func (f *Function) String() string {
	if f.IsPhysical() {
		return fmt.Sprintf("function: 'physical' virtual: '%d/%d'", len(f.Virtual), f.MaxVirtual)
	}
	return "function: 'virtual'"
}
