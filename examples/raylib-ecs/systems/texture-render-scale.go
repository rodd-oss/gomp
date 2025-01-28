/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkg/ecs"
)

// TextureRenderScale is a system that sets Scale of textureRender
type trScaleController struct{}

func (s *trScaleController) Init(world *ecs2.World)        {}
func (s *trScaleController) FixedUpdate(world *ecs2.World) {}
func (s *trScaleController) Update(world *ecs2.World) {
	// Get component managers
	scales := components.ScaleService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	textureRenders.AllParallel(func(entity ecs2.Entity, tr *components.TextureRender) bool {
		if tr == nil {
			return true
		}

		scale := scales.Get(entity)
		if scale == nil {
			return true
		}

		tr.Dest.Width *= scale.X
		tr.Dest.Height *= scale.Y

		return true
	})
}
func (s *trScaleController) Destroy(world *ecs2.World) {}
