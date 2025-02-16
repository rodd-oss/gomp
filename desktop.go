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

package gomp

import (
	"gomp/pkg/ecs"
)

type DesktopEngine = Engine[*DesktopComponents, *DesktopSystems]

func NewDesktopEngine(scenes []Scene) *DesktopEngine {
	world := ecs.CreateWorld("main")
	defer world.Destroy()

	components := NewDesktopComponents(world)
	systems := NewDesktopSystems(world, components, scenes)

	return NewEngine(world, components, systems)
}
