/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewTextureRenderSpriteSheetSystem(
	spriteSheets *stdcomponents.SpriteSheetComponentManager,
	textureRenders *stdcomponents.TextureRenderComponentManager,
) *TextureRenderSpriteSheetSystem {
	return &TextureRenderSpriteSheetSystem{
		spriteSheets:   spriteSheets,
		textureRenders: textureRenders,
	}
}

// TextureRenderSpriteSheetSystem is a system that prepares SpriteSheet to be rendered
type TextureRenderSpriteSheetSystem struct {
	spriteSheets   *stdcomponents.SpriteSheetComponentManager
	textureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderSpriteSheetSystem) Init() {}
func (s *TextureRenderSpriteSheetSystem) Run(dt time.Duration) {
	s.spriteSheets.AllParallel(func(entity ecs.Entity, spriteSheet *stdcomponents.SpriteSheet) bool {
		if spriteSheet == nil {
			return true
		}

		tr := s.textureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := stdcomponents.TextureRender{
				Texture: spriteSheet.Texture,
				Frame:   spriteSheet.Frame,
				Origin:  spriteSheet.Origin,
				Dest: rl.NewRectangle(
					0,
					0,
					spriteSheet.Frame.Width,
					spriteSheet.Frame.Height,
				),
			}

			s.textureRenders.Create(entity, newRender)
		} else {
			// Run spriteRender
			tr.Texture = spriteSheet.Texture
			tr.Frame = spriteSheet.Frame
			tr.Origin = spriteSheet.Origin
		}
		return true
	})
}
func (s *TextureRenderSpriteSheetSystem) Destroy() {}
