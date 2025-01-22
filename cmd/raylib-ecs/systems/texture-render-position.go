/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
)

// TextureRenderPosition is a system that sets Position of textureRender
type trPositionController struct{}

func (s *trPositionController) Init(world *ecs.World)        {}
func (s *trPositionController) FixedUpdate(world *ecs.World) {}
func (s *trPositionController) Update(world *ecs.World) {
	// Get component managers
	positions := components.PositionService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	textureRenders.AllParallel(func(entity ecs.EntityID, tr *components.TextureRender) bool {
		if tr == nil {
			return true
		}

		position := positions.Get(entity)
		if position == nil {
			return true
		}

		tr.Dest.X = position.X
		tr.Dest.Y = position.Y

		return true
	})
}
func (s *trPositionController) Destroy(world *ecs.World) {}
