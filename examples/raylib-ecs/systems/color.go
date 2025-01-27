/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkgs/ecs"
	"image/color"
)

type colorController struct {
	baseColor color.RGBA
}

func (s *colorController) Init(world *ecs2.World) {
	s.baseColor = color.RGBA{25, 220, 200, 255}
}
func (s *colorController) Update(world *ecs2.World) {}
func (s *colorController) FixedUpdate(world *ecs2.World) {
	sprites := components.SpriteService.GetManager(world)
	hps := components.HealthService.GetManager(world)

	sprites.AllParallel(func(entity ecs2.Entity, sprite *components.Sprite) bool {
		hp := hps.Get(entity)
		if hp == nil {
			return true
		}

		hpPercentage := float32(hp.Hp) / float32(hp.MaxHp)

		sprite.Tint = color.RGBA{
			uint8(hpPercentage * float32(s.baseColor.R)),
			uint8(hpPercentage * float32(s.baseColor.G)),
			uint8(hpPercentage * float32(s.baseColor.B)),
			s.baseColor.A,
		}
		return true
	})
}
func (s *colorController) Destroy(world *ecs2.World) {}
