/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/stdcomponents"
)

func NewTextureRenderSpriteSheetSystem() TextureRenderSpriteSheetSystem {
	return TextureRenderSpriteSheetSystem{}
}

// TextureRenderSpriteSheetSystem is a system that prepares SpriteSheet to be rendered
type TextureRenderSpriteSheetSystem struct {
	SpriteSheets   *stdcomponents.SpriteSheetComponentManager
	TextureRenders *stdcomponents.RLTextureProComponentManager
}

func (s *TextureRenderSpriteSheetSystem) Init() {}
func (s *TextureRenderSpriteSheetSystem) Run() {
	//s.SpriteSheets.EachParallel(func(entity ecs.Entity, spriteSheet *stdcomponents.SpriteSheet) bool {
	//	if spriteSheet == nil {
	//		return true
	//	}
	//
	//	tr := s.RlTexturePros.GetUnsafe(entity)
	//	if tr == nil {
	//		// Create new spriteRender
	//		newRender := stdcomponents.RLTexturePro{
	//			TextureId: spriteSheet.TextureId,
	//			Frame:   spriteSheet.Frame,
	//			Origin:  spriteSheet.Origin,
	//			Dest: rl.NewRectangle(
	//				0,
	//				0,
	//				spriteSheet.Frame.Width,
	//				spriteSheet.Frame.Height,
	//			),
	//		}
	//
	//		s.RlTexturePros.Create(entity, newRender)
	//	} else {
	//		// Run spriteRender
	//		tr.TextureId = spriteSheet.TextureId
	//		tr.Frame = spriteSheet.Frame
	//		tr.Origin = spriteSheet.Origin
	//	}
	//	return true
	//})
}
func (s *TextureRenderSpriteSheetSystem) Destroy() {}
