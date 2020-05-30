// +build !linux,!windows
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"

	"github.com/jaypipes/ghw/pkg/context"
)

func topologyFillInfo(ctx *context.Context, info *TopologyInfo) error {
	return errors.New("topologyFillInfo not implemented on " + runtime.GOOS)
}

// TopologyNodes has been deprecated in 0.2. Please use the TopologyInfo.Nodes
// attribute.
// TODO(jaypipes): Remove in 1.0.
func TopologyNodes() ([]*TopologyNode, error) {
	msg := `
The TopologyNodes() function has been DEPRECATED and will be removed in the 1.0
release of ghw. Please use the TopologyInfo.Nodes attribute.
`
	warn(msg)
	ctx := context.FromEnv()
	return topologyNodes(ctx)
}

func topologyNodes(ctx *context.Context) ([]*TopologyNode, error) {
	return nil, errors.New("Don't know how to get topology on " + runtime.GOOS)
}
