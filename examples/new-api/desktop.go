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

package main

import (
	"gomp"
	"gomp/pkg/ecs"
)

func main() {
	world := ecs.CreateWorld("main")
	defer world.Destroy()

	components := NewDesktopComponents(world)
	systems := NewDesktopSystems(world, components)

	engine := gomp.NewEngine(world, components, systems)
	engine.Run(60)
}
