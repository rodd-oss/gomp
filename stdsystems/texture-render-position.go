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

func NewTextureRenderPositionSystem() TextureRenderPositionSystem {
	return TextureRenderPositionSystem{}
}

// TextureRenderPositionSystem is a system that sets Position of textureRender
type TextureRenderPositionSystem struct {
	Positions      *stdcomponents.PositionComponentManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
}

func (s *TextureRenderPositionSystem) Init() {}
func (s *TextureRenderPositionSystem) Run(dt time.Duration) {
	s.TextureRenders.AllParallel(func(entity ecs.Entity, tr *stdcomponents.TextureRender) bool {
		if tr == nil {
			return true
		}

		position := s.Positions.Get(entity)
		if position == nil {
			return true
		}

		tr.Dest.X = position.X
		tr.Dest.Y = position.Y

		return true
	})
}
func (s *TextureRenderPositionSystem) Destroy() {}
