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

func NewTextureRenderRotationSystem() TextureRenderRotationSystem {
	return TextureRenderRotationSystem{}
}

// TextureRenderRotationSystem is a system that sets Rotation of textureRender
type TextureRenderRotationSystem struct {
	Rotations      *stdcomponents.RotationComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderRotationSystem) Init() {}
func (s *TextureRenderRotationSystem) Run(dt time.Duration) {
	// Run sprites and spriteRenders
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		rotation := s.Rotations.Get(entity)
		if rotation == nil {
			return true
		}

		tr.Rotation = rotation.Angle

		return true
	})
}
func (s *TextureRenderRotationSystem) Destroy() {}
