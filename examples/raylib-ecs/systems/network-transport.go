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
	"gomp/network"
	"gomp/pkg/ecs"
)

type networkTransportController struct {
	transport network.AnyNetworkTransport
	entity    ecs.Entity
}

func (s *networkTransportController) Init(world *ecs.World) {
	s.entity = world.CreateEntity("network-transport")
}
func (s *networkTransportController) Update(world *ecs.World) {
	//patchInComponent := components.PatchInComponentService.GetManager(world)
	//patchOutComponent := components.PatchOutComponentService.GetManager(world)
	//
	//if !s.transport.IsActive() {
	//	return
	//}
	//
	//numOfMsgs := len(s.transport.Receive())
	//for range numOfMsgs {
	//	data := <-s.transport.Receive()
	//	patch := components.PatchIn{Data: data}
	//	patchInComponent.Create(s.entity, patch)
	//}
	//
	//
	//patchOutComponent.All(func(entity ecs.Entity, patch *components.PatchOut) bool {
	//	s.transport.Send(patch.Data)
	//	return true
	//})
}
func (s *networkTransportController) FixedUpdate(world *ecs.World) {}
func (s *networkTransportController) Destroy(world *ecs.World)     {}
