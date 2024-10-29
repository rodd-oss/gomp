/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
)

type Scene struct {
	Engine *Engine

	World ecs.World
	Space *cp.Space

	Entities   []ecs.Entity
	Components []*ecs.ComponentType[NetworkSyncComponent[any, any]]

	currentTick uint
	syncPeriod  uint // in ticks
}

func (s *Scene) Update(dt float64) {
	needToSync := s.currentTick%s.syncPeriod == 0

	for i := range s.Components {
		s.Components[i].Each(s.World, func(e *ecs.Entry) {
			s.Components[i].GetValue(e).Update(dt)

			if needToSync {
				s.Components[i].GetValue(e).Sync()
			}
		})
	}

	s.Space.Step(dt)

	if s.currentTick%s.syncPeriod == 0 {
		//send s.Patch
	}

	// network sync
	s.currentTick++
}

type EntityState struct {
	Id   uint
}

type SceneState map[uint]type