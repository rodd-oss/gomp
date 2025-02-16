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

func NewTextureRenderTintSystem(tint *ecs.ComponentManager[stdcomponents.Tint], render *ecs.ComponentManager[stdcomponents.TextureRender]) *TextureRenderTintSystem {
	return &TextureRenderTintSystem{
		tints:          tint,
		textureRenders: render,
	}
}

// TextureRenderTintSystem is a system that sets Scale of textureRender
type TextureRenderTintSystem struct {
	tints          *stdcomponents.TintComponentManager
	textureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderTintSystem) Init() {}
func (s *TextureRenderTintSystem) Run(dt time.Duration) {
	s.textureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		tint := s.tints.Get(entity)
		if tint == nil {
			return true
		}

		trTint := &tr.Tint
		trTint.A = tint.A
		trTint.R = tint.R
		trTint.G = tint.G
		trTint.B = tint.B

		return true
	})
}
func (s *TextureRenderTintSystem) Destroy() {}
