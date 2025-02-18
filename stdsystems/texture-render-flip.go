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

func NewTextureRenderFlipSystem() TextureRenderFlipSystem {
	return TextureRenderFlipSystem{}
}

// TextureRenderFlipSystem is a system that sets Scale of textureRender
type TextureRenderFlipSystem struct {
	Flips          *stdcomponents.FlipComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderFlipSystem) Init() {}
func (s *TextureRenderFlipSystem) Run(dt time.Duration) {
	// Run sprites and spriteRenders
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		mirrored := s.Flips.Get(entity)
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
