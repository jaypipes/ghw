// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"runtime"

	"github.com/pkg/errors"
)

func topologyFillInfo(info *TopologyInfo) error {
	return errors.New("topologyFillInfo not implemented on " + runtime.GOOS)
}

func TopologyNodes() ([]*TopologyNode, error) {
	return nil, errors.New("Don't know how to get topology on " + runtime.GOOS)
}
