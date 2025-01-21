/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file deveopment:
-===-===-===-===-===-===-===-===-===-===

<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteController struct {
	colors     *ecs.ComponentManager[color.RGBA]
	transforms *ecs.ComponentManager[components.Transform]
	sprites    *ecs.ComponentManager[components.Sprite]
	t          rl.Texture2D
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
	s.colors.AllParallel(func(entity ecs.EntityID, color *color.RGBA) bool {
		if color == nil {
			return true
		}

		transform := s.transforms.Get(entity)
		if transform == nil {
			return true
		}

		sprite := s.sprites.Get(entity)
		if sprite == nil {
			sprite := components.Sprite{
				Position: rl.NewVector2(float32(transform.X), float32(transform.Y)),
				Tint:     *color,
			}

			s.sprites.Create(entity, sprite)
		} else {
			sprite.Tint = *color
		}

		return true
	})
}

func (s *spriteController) Destroy(world *ecs.World) {}

// ---------------
// System private methods
// ---------------
