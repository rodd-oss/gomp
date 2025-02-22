/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewTextureRenderSpriteSystem() TextureRenderSpriteSystem {
	return TextureRenderSpriteSystem{}
}

// TextureRenderSpriteSystem is a system that prepares Sprite to be rendered
type TextureRenderSpriteSystem struct {
	Sprites        *stdcomponents.SpriteComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderSpriteSystem) Init() {}
func (s *TextureRenderSpriteSystem) Run() {
	// Run sprites and spriteRenders
	s.Sprites.AllParallel(func(entity ecs.Entity, sprite *stdcomponents.Sprite) bool {
		if sprite == nil {
			return true
		}

		spriteFrame := sprite.Frame
		spriteOrigin := sprite.Origin
		spriteTint := sprite.Tint

		tr := s.TextureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := stdcomponents.TextureRender{
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

			s.TextureRenders.Create(entity, newRender)
		} else {
			// Run spriteRender
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
func (s *TextureRenderSpriteSystem) Destroy() {}
