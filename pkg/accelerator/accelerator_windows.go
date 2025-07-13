// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

import "github.com/jaypipes/ghw/pkg/option"

func (i *Info) load(opt ...option.Option) error {
	i.Devices = []*AcceleratorDevice{}
	return nil
}
