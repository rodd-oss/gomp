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

package systems

import (
	"fmt"
	"gomp/examples/raylib-ecs/components"
	"gomp/network"
	"gomp/pkg/ecs"
)

type networkSendController struct{}

func (s *networkSendController) Init(world *ecs.EntityManager) {
	positions := components.PositionService.GetManager(world)
	rotations := components.RotationService.GetManager(world)
	mirroreds := components.MirroredService.GetManager(world)

	positions.TrackChanges = true
	rotations.TrackChanges = true
	mirroreds.TrackChanges = true

	positions.SetEncoder(func(components []components.Position) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.X)
			data = append(data, []byte(binary)...)
		}

		return data
	})
	rotations.SetEncoder(func(components []components.Rotation) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.Angle)
			data = append(data, []byte(binary)...)
		}

		return data
	})
	mirroreds.SetEncoder(func(components []components.Mirrored) []byte {
		data := make([]byte, 0)
		for _, component := range components {
			binary := fmt.Sprintf("%b", component.X)
			data = append(data, []byte(binary)...)
		}

		return data
	})
}
func (s *networkSendController) Update(world *ecs.EntityManager) {}
func (s *networkSendController) FixedUpdate(world *ecs.EntityManager) {
	//patch := world.PatchGet()
	//world.PatchReset()
	//log.Printf("%v", patch)
	if network.Quic.Mode() != network.ModeNone {
		network.Quic.Send([]byte("patch"), 0)
	}
}
func (s *networkSendController) Destroy(world *ecs.EntityManager) {}
