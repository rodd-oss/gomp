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

package components

import "gomp/pkg/ecs"

type Health struct {
	Hp, MaxHp int32
}

type HealthComponentManager = ecs.ComponentManager[Health]

func NewHealthComponentManager(world *ecs.World) *ecs.ComponentManager[Health] {
	return ecs.NewComponentManager[Health](world, HEALTH_ID)
}
