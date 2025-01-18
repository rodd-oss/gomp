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

	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteController struct {
	colors     ecs.WorldComponents[color.RGBA]
	transforms ecs.WorldComponents[components.Transform]
	sprites    ecs.WorldComponents[components.Sprite]
}

// ---------------
// Basic System Methods
// ---------------

func (s *spriteController) Init(world *ecs.World) {
	s.colors = components.ColorService.GetManager(world)
	s.transforms = components.TransformService.GetManager(world)
	s.sprites = components.SpriteService.GetManager(world)
}

func (s *spriteController) Update(world *ecs.World) {}

func (s *spriteController) FixedUpdate(world *ecs.World) {
	s.colors.AllParallel(s.mapColorsToSprites)
}

func (s *spriteController) Destroy(world *ecs.World) {}

// ---------------
// System private methods
// ---------------

func (s *spriteController) mapColorsToSprites(entity ecs.EntityID, color *color.RGBA) bool {
	if color == nil {
		return true
	}

	transform := s.transforms.GetPtr(entity)
	if transform == nil {
		return true
	}

	sprite := s.sprites.GetPtr(entity)
	if sprite == nil {
		sprite := components.Sprite{
			Pos:  rl.NewVector2(float32(transform.X), float32(transform.Y)),
			Tint: *color,
		}

		s.sprites.Set(entity, sprite)
	} else {
		sprite.Tint = *color
	}

	return true
}
