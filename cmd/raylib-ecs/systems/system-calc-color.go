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

type systemCalcColor struct {
	baseColor color.RGBA
}

func (s *systemCalcColor) Init(world *ecs.World) {
	s.baseColor = color.RGBA{25, 220, 200, 255}
}

func (s *systemCalcColor) Run(world *ecs.World) {
	colors := components.ColorManager.Instances(world)
	healths := components.HealthManager.Instances(world)

	colors.AllParallel(func(entity ecs.EntityID, color *color.RGBA) bool {
		health := healths.GetPtr(entity)
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
func (s *systemCalcColor) Destroy(world *ecs.World) {}
