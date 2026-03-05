// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package accelerator

import "context"

func (i *Info) load(ctx context.Context) error {
	i.Devices = []*AcceleratorDevice{}
	return nil
}
