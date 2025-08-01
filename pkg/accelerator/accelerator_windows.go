// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

func (i *Info) load() error {
	i.Devices = []*AcceleratorDevice{}
	return nil
}
