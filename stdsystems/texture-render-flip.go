/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewTextureRenderFlipSystem(
	flips *stdcomponents.FlipComponentManager,
	textureRenders *stdcomponents.TextureRenderComponentManager,
) *TextureRenderFlipSystem {
	return &TextureRenderFlipSystem{
		flips:          flips,
		textureRenders: textureRenders,
	}
}

// TextureRenderFlipSystem is a system that sets Scale of textureRender
type TextureRenderFlipSystem struct {
	flips          *stdcomponents.FlipComponentManager
	textureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderFlipSystem) Init() {}
func (s *TextureRenderFlipSystem) Run(dt time.Duration) {
	// Run sprites and spriteRenders
	s.textureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		mirrored := s.flips.Get(entity)
		if mirrored == nil {
			return true
		}

		if mirrored.X {
			tr.Frame.Width *= -1
		}
		if mirrored.Y {
			tr.Frame.Height *= -1
		}

		return true
	})
}
func (s *TextureRenderFlipSystem) Destroy() {}
