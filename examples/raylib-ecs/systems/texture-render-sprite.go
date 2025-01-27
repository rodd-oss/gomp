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
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkgs/ecs"
)

// TextureRenderSprite is a system that prepares Sprite to be rendered
type trSpriteController struct{}

func (s *trSpriteController) Init(world *ecs2.World)        {}
func (s *trSpriteController) FixedUpdate(world *ecs2.World) {}
func (s *trSpriteController) Update(world *ecs2.World) {
	// Get component managers
	sprites := components.SpriteService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	sprites.AllParallel(func(entity ecs2.Entity, sprite *components.Sprite) bool {
		if sprite == nil {
			return true
		}

		spriteFrame := sprite.Frame
		spriteOrigin := sprite.Origin
		spriteTint := sprite.Tint

		tr := textureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := components.TextureRender{
				Texture: sprite.Texture,
				Frame:   sprite.Frame,
				Origin:  sprite.Origin,
				Tint:    sprite.Tint,
				Dest: rl.NewRectangle(
					0,
					0,
					sprite.Frame.Width,
					sprite.Frame.Height,
				),
			}

			textureRenders.Create(entity, newRender)
		} else {
			// Update spriteRender
			// tr.Texture = sprite.Texture
			trFrame := &tr.Frame
			trFrame.X = spriteFrame.X
			trFrame.Y = spriteFrame.Y
			trFrame.Width = spriteFrame.Width
			trFrame.Height = spriteFrame.Height

			trOrigin := &tr.Origin
			trOrigin.X = spriteOrigin.X
			trOrigin.Y = spriteOrigin.Y

			trTint := &tr.Tint
			trTint.A = spriteTint.A
			trTint.R = spriteTint.R
			trTint.G = spriteTint.G
			trTint.B = spriteTint.B

			trDest := &tr.Dest
			trDest.Width = spriteFrame.Width
			trDest.Height = spriteFrame.Height
		}
		return true
	})
}
func (s *trSpriteController) Destroy(world *ecs2.World) {}
