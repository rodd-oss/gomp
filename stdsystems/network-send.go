/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package stdsystems

import (
	"fmt"
	"gomp/network"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewNetworkSendSystem() NetworkSendSystem {
	return NetworkSendSystem{}
}

type NetworkSendSystem struct {
	World     *ecs.EntityManager
	Positions *stdcomponents.PositionComponentManager
	Rotations *stdcomponents.RotationComponentManager
	Mirroreds *stdcomponents.FlipComponentManager
}

func (s *NetworkSendSystem) Init() {
	s.Positions.TrackChanges = true
	s.Rotations.TrackChanges = true
	s.Mirroreds.TrackChanges = true

	s.Positions.SetEncoder(func(components []stdcomponents.Position) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.XY.X)
			data = append(data, []byte(binary)...)
		}

		return data
	})
	s.Rotations.SetEncoder(func(components []stdcomponents.Rotation) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.Angle)
			data = append(data, []byte(binary)...)
		}

		return data
	})
	s.Mirroreds.SetEncoder(func(components []stdcomponents.Flip) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.X)
			data = append(data, []byte(binary)...)
		}

		return data
	})
}
func (s *NetworkSendSystem) Run(dt time.Duration) {
	//patch := world.PatchGet()
	//world.PatchReset()
	//log.Printf("%v", patch)
	if network.Quic.Mode() != network.ModeNone {
		network.Quic.Send([]byte("patch"), 0)
	}
}
func (s *NetworkSendSystem) Destroy() {}
