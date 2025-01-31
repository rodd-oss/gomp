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

package systems

import (
	"gomp/examples/raylib-ecs/components"
	"gomp/pkg/ecs"
	"sync/atomic"
)

type NetworkMode int

const (
	Server NetworkMode = iota
	Client
)

type networkController struct {
	mode   NetworkMode
	lookup *ecs.PagedMap[ecs.Entity, components.NetworkId]
	nextId atomic.Int32
}

func (s *networkController) Init(world *ecs.World) {
	s.lookup = ecs.NewPagedMap[ecs.Entity, components.NetworkId]()
	s.nextId.Store(0)
}
func (s *networkController) Update(world *ecs.World) {
	networks := components.NetworkComponentService.GetManager(world)

	switch s.mode {
	case Server:
		networks.All(func(entity ecs.Entity, network *components.Network) bool {
			if network.Id == 0 {
				network.Id = s.registerEntity(entity)
			}

			return true
		})

	case Client:
		networks.All(func(entity ecs.Entity, network *components.Network) bool {
			// apply patch to entity
			return true
		})
	}
}
func (s *networkController) FixedUpdate(world *ecs.World) {}
func (s *networkController) Destroy(world *ecs.World)     {}

func (s *networkController) registerEntity(entity ecs.Entity) components.NetworkId {
	id := components.NetworkId(s.nextId.Add(1))
	s.lookup.Set(entity, id)
	return id
}
