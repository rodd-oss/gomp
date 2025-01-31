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
	"gomp/pkg/ecs"
)

type networkPatchInController struct{}

func (s *networkPatchInController) Init(world *ecs.World)   {}
func (s *networkPatchInController) Update(world *ecs.World) {}
func (s *networkPatchInController) FixedUpdate(world *ecs.World) {
	//patchInComponents := components.PatchInComponentService.GetManager(world)
	//
	//patchInComponents.All(func(entity ecs.Entity, p *components.PatchIn) bool {
	//	data := p.Data
	//	patch := parse(data)
	//	component := world.getcomponentbyid(patch.id)
	//	component.applyPatch(patch.data)
	//	return true
	//})
}
func (s *networkPatchInController) Destroy(world *ecs.World) {}
