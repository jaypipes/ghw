// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

type Win32_VideoController struct {
	Caption              string
	Description          string
	DeviceID             string
	Name                 string
	VideoArchitecture    uint16
	VideoMemoryType      uint16
	VideoModeDescription string
	VideoProcessor       string
}

func (ctx *context) gpuFillInfo(info *GPUInfo) error {
	// Getting disk drives from WMI
	var win32VideoControllerDescriptions []Win32_VideoController
	q1 := wmi.CreateQuery(&win32VideoControllerDescriptions, "")
	if err := wmi.Query(q1, &win32VideoControllerDescriptions); err != nil {
		return err
	}
	//fmt.Printf("TEST: %#v\n", win32VideoControllerDescriptions)
	// Converting into standard structures
	cards := make([]*GraphicsCard, 0)
	for _, description := range win32VideoControllerDescriptions {
		//topo := &TopologyInfo{}
		//topo.Architecture
		//topo.Nodes
		//node := &TopologyNode{}
		/*
			node.Cores
			c := &ProcessorCore{
				Id:                coreID,
				ID:                coreID,
				Index:             len(cores),
				LogicalProcessors: make([]int, 0),
			}
			node.Caches
			caches := make(map[string]*MemoryCache)
			cache = &MemoryCache{
					Level:             uint8(level),
					Type:              cacheType,
					SizeBytes:         uint64(size) * uint64(KB),
					LogicalProcessors: make([]uint32, 0),
				}
		*/
		cores := make([]*ProcessorCore, 0)
		core := &ProcessorCore{
			Id:                0,
			ID:                0,
			Index:             0,
			LogicalProcessors: make([]int, 0),
		}
		cores = append(cores, core)
		caches := make([]*MemoryCache, 0)
		cache := &MemoryCache{
			Level:             0,
			Type:              0,
			SizeBytes:         0,
			LogicalProcessors: make([]uint32, 0),
		}
		caches = append(caches, cache)
		card := &GraphicsCard{
			Address: description.DeviceID, // https://stackoverflow.com/questions/32073667/how-do-i-discover-the-pcie-bus-topology-and-slot-numbers-on-the-board
			Index:   0,
			Node: &TopologyNode{
				Cores:  cores,
				Caches: caches,
			},
		}
		cards = append(cards, card)
	}
	//car.DeviceInfo
	//card.Node
	info.GraphicsCards = cards
	return nil
}
