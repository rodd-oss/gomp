/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
)

type colorController struct {
	baseColor color.RGBA
}

func (s *colorController) Init(world *ecs.World) {
	s.baseColor = color.RGBA{25, 220, 200, 255}
}
func (s *colorController) Update(world *ecs.World) {}
func (s *colorController) FixedUpdate(world *ecs.World) {
	colorManager := components.SpriteService.GetManager(world)
	healthManager := components.HealthService.GetManager(world)

	colorManager.AllParallel(func(entity ecs.EntityID, sprite *components.Sprite) bool {
		health := healthManager.Get(entity)
		if health == nil {
			return true
		}

		hpPercentage := float32(health.Hp) / float32(health.MaxHp)

		sprite.Tint = color.RGBA{
			uint8(hpPercentage * float32(s.baseColor.R)),
			uint8(hpPercentage * float32(s.baseColor.G)),
			uint8(hpPercentage * float32(s.baseColor.B)),
			s.baseColor.A,
		}
		return true
	})
}
func (s *colorController) Destroy(world *ecs.World) {}
