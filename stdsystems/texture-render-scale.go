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

func NewTextureRenderScaleSystem(scale *ecs.ComponentManager[stdcomponents.Scale], render *ecs.ComponentManager[stdcomponents.TextureRender]) *TextureRenderScaleSystem {
	return &TextureRenderScaleSystem{
		scales:         scale,
		textureRenders: render,
	}
}

// TextureRenderScaleSystem is a system that sets Scale of textureRender
type TextureRenderScaleSystem struct {
	scales         *stdcomponents.ScaleComponentManager
	textureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderScaleSystem) Init() {}
func (s *TextureRenderScaleSystem) Run(dt time.Duration) {
	s.textureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		scale := s.scales.Get(entity)
		if scale == nil {
			return true
		}

		tr.Dest.Width *= scale.X
		tr.Dest.Height *= scale.Y

		return true
	})
}
func (s *TextureRenderScaleSystem) Destroy() {}
