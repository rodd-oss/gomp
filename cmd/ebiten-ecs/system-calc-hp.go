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

type systemCalcHp struct {
	baseColor color.RGBA
}

func (s *systemCalcHp) Init(world *ClientWorld) {
	s.baseColor = color.RGBA{25, 220, 200, 255}
}

func (s *systemCalcHp) Run(world *ClientWorld) {
	world.Components.Health.AllParallel(func(entity ecs.EntityID, h *health) bool {
		h.hp--

		if h.hp <= 0 {
			world.Components.Destroy.Create(entity, struct{}{})
		}
		hpPercentage := float32(h.hp) / float32(h.maxHp)
		h.color.R = uint8(hpPercentage * float32(s.baseColor.R))
		h.color.G = uint8(hpPercentage * float32(s.baseColor.G))
		h.color.B = uint8(hpPercentage * float32(s.baseColor.B))
		h.color.A = s.baseColor.A

		return true
	})
}
func (s *systemCalcHp) Destroy(world *ClientWorld) {}
