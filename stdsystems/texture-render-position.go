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

func NewTextureRenderPositionSystem() TextureRenderPositionSystem {
	return TextureRenderPositionSystem{}
}

// TextureRenderPositionSystem is a system that sets Position of textureRender
type TextureRenderPositionSystem struct {
	ViewPositions  *stdcomponents.ViewPositionComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderPositionSystem) Init() {}
func (s *TextureRenderPositionSystem) Run(interpolation float32) {
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		viewPosition := s.ViewPositions.Get(entity)
		if viewPosition == nil {
			return true
		}

		// Lerp(start, end, amount float32) start + amount*(end-start)
		interpolationX := interpolation * (viewPosition.CurrentPositionX - viewPosition.LastPositionX)
		interpolationY := interpolation * (viewPosition.CurrentPositionY - viewPosition.LastPositionY)
		tr.Dest.X = viewPosition.LastPositionX + interpolationX
		tr.Dest.Y = viewPosition.LastPositionY + interpolationY

		return true
	})
}
func (s *TextureRenderPositionSystem) Destroy() {}
