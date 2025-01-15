/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
)

type colorSystem struct {
	transformComponent ecs.WorldComponents[transform]
	healthComponent    ecs.WorldComponents[health]
	colorComponent     ecs.WorldComponents[color.RGBA]
	movementComponent  ecs.WorldComponents[movement]

	baseColor color.RGBA
}

func (s *colorSystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)

	s.baseColor = color.RGBA{25, 220, 200, 255}
}
func (s *colorSystem) Run(world *ecs.World) {
	s.colorComponent.AllParallel(func(ei ecs.EntityID, c *color.RGBA) bool {
		health := s.healthComponent.GetPtr(ei)
		if health == nil {
			return true
		}

		hpPercentage := float32(health.hp) / float32(health.maxHp)

		c.R = uint8(hpPercentage * float32(s.baseColor.R))
		c.G = uint8(hpPercentage * float32(s.baseColor.G))
		c.B = uint8(hpPercentage * float32(s.baseColor.B))
		c.A = s.baseColor.A

		return true
	})
}
func (s *colorSystem) Destroy(world *ecs.World) {}
