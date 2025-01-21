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

	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteController struct {
	t rl.Texture2D
}

// ---------------
// Basic System Methods
// ---------------

func (s *spriteController) Init(world *ecs.World) {}

func (s *spriteController) Update(world *ecs.World) {}

func (s *spriteController) FixedUpdate(world *ecs.World) {
	// Get component managers
	positions := components.PositionService.GetManager(world)
	rotations := components.RotationService.GetManager(world)
	scales := components.ScaleService.GetManager(world)
	sprites := components.SpriteService.GetManager(world)
	spriteRenders := components.SpriteRenderService.GetManager(world)

	// Update sprites and spriteRenders
	sprites.AllParallel(func(entity ecs.EntityID, sprite *components.Sprite) bool {
		if sprite == nil {
			return true
		}

		position := positions.Get(entity)
		if position == nil {
			return true
		}

		rotation := rotations.Get(entity)
		if rotation == nil {
			return true
		}

		scale := scales.Get(entity)
		if scale == nil {
			return true
		}

		spriteTextureRegion := &sprite.TextureRegion
		spriteOrigin := &sprite.Origin
		spriteTint := &sprite.Tint

		spriteRender := spriteRenders.Get(entity)
		if spriteRender == nil {
			// Create new spriteRender
			newRender := components.SpriteRender{
				Sprite: *sprite,
				Dest: rl.NewRectangle(
					position.X,
					position.Y,
					spriteTextureRegion.Width*scale.X,
					spriteTextureRegion.Height*scale.Y,
				),
				Rotation: rotation.Angle,
			}

			spriteRenders.Create(entity, newRender)
		} else {
			renderSprite := &spriteRender.Sprite
			renderDest := &spriteRender.Dest

			// Update TextureRegion
			renderTextureRegion := &renderSprite.TextureRegion
			renderTextureRegion.X = spriteTextureRegion.X
			renderTextureRegion.Y = spriteTextureRegion.Y
			renderTextureRegion.Width = spriteTextureRegion.Width
			renderTextureRegion.Height = spriteTextureRegion.Height

			// Update Origin
			renderOrigin := &renderSprite.Origin
			renderOrigin.X = spriteTextureRegion.Width * spriteOrigin.X * scale.X
			renderOrigin.Y = spriteTextureRegion.Height * spriteOrigin.Y * scale.Y

			// Update Tint
			renderTint := &renderSprite.Tint
			renderTint.A = spriteTint.A
			renderTint.R = spriteTint.R
			renderTint.G = spriteTint.G
			renderTint.B = spriteTint.B

			// Update destination
			renderDest.X = position.X
			renderDest.Y = position.Y
			renderDest.Width = spriteTextureRegion.Width * scale.X
			renderDest.Height = spriteTextureRegion.Height * scale.Y

			// Update rotation
			spriteRender.Rotation = rotation.Angle
		}

		return true
	})
}

func (s *spriteController) Destroy(world *ecs.World) {}

// ---------------
// System private methods
// ---------------
