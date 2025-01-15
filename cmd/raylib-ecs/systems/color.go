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

func (s *colorController) Run(world *ecs.World) {
	colorManager := components.ColorService.GetManager(world)
	healthManager := components.HealthService.GetManager(world)

	colorManager.AllParallel(func(entity ecs.EntityID, color *color.RGBA) bool {
		health := healthManager.GetPtr(entity)
		if health == nil {
			return true
		}

		hpPercentage := float32(health.Hp) / float32(health.MaxHp)

		color.R = uint8(hpPercentage * float32(s.baseColor.R))
		color.G = uint8(hpPercentage * float32(s.baseColor.G))
		color.B = uint8(hpPercentage * float32(s.baseColor.B))
		color.A = s.baseColor.A

		return true
	})
}
func (s *colorController) Destroy(world *ecs.World) {}
