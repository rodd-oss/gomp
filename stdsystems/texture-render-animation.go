/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewTextureRenderAnimationSystem() TextureRenderAnimationSystem {
	return TextureRenderAnimationSystem{}
}

// TextureRenderAnimationSystem is a system that sets Position of textureRender
type TextureRenderAnimationSystem struct {
	Animations     *stdcomponents.AnimationPlayerComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderAnimationSystem) Init() {}
func (s *TextureRenderAnimationSystem) Run() {
	// Run sprites and spriteRenders
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		animation := s.Animations.Get(entity)
		if animation == nil {
			return true
		}

		frame := &tr.Frame
		if animation.Vertical {
			frame.Y += frame.Height * float32(animation.Current)
		} else {
			frame.X += frame.Width * float32(animation.Current)
		}

		return true
	})
}
func (s *TextureRenderAnimationSystem) Destroy() {}
