/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewViewPositionSystem() ViewPositionSystem {
	return ViewPositionSystem{}
}

type ViewPositionSystem struct {
	EntityManager *ecs.EntityManager
	Positions     *stdcomponents.PositionComponentManager
	ViewPositions *stdcomponents.ViewPositionComponentManager
}

func (s *ViewPositionSystem) Init() {}
func (s *ViewPositionSystem) Run() {
	s.Positions.AllParallel(func(entity ecs.Entity, p *stdcomponents.Position) bool {
		viewPosition := s.ViewPositions.Get(entity)

		if viewPosition == nil {
			viewPosition = s.ViewPositions.Create(entity, stdcomponents.ViewPosition{})
		}

		viewPosition.LastPositionX = viewPosition.CurrentPositionX
		viewPosition.LastPositionY = viewPosition.CurrentPositionY
		viewPosition.CurrentPositionX = p.X
		viewPosition.CurrentPositionY = p.Y
		return true
	})
}
func (s *ViewPositionSystem) Destroy() {}
