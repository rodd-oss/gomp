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

func NewTextureRenderTintSystem() TextureRenderTintSystem {
	return TextureRenderTintSystem{}
}

// TextureRenderTintSystem is a system that sets Scale of textureRender
type TextureRenderTintSystem struct {
	Tints          *stdcomponents.TintComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderTintSystem) Init() {}
func (s *TextureRenderTintSystem) Run(dt time.Duration) {
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		tint := s.Tints.Get(entity)
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
