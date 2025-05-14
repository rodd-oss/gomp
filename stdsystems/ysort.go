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
	"gomp/pkg/worker"
	"gomp/stdcomponents"
)

const ySortOffsetScale float32 = 0.001

func NewYSortSystem() YSortSystem {
	return YSortSystem{}
}

type YSortSystem struct {
	EntityManager *ecs.EntityManager
	YSorts        *stdcomponents.YSortComponentManager
	Positions     *stdcomponents.PositionComponentManager
	RenderOrders  *stdcomponents.RenderOrderComponentManager
}

func (s *YSortSystem) Init() {}
func (s *YSortSystem) Run() {
	s.YSorts.ProcessEntities(func(entity ecs.Entity, _ worker.WorkerId) {
		pos := s.Positions.GetUnsafe(entity)
		renderOrder := s.RenderOrders.GetUnsafe(entity)

		// Calculate depth based on Y position
		yDepth := pos.XY.Y * ySortOffsetScale

		// Preserve original Z layer but add Y-based offset
		//renderOrder.CalculatedZ = float32(int(pos.Z)) + yDepth
		renderOrder.CalculatedZ = yDepth
	})
}
func (s *YSortSystem) Destroy() {}
