/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/cmd/raylib-ecs/systems"
	"gomp_game/pkgs/gomp/ecs"
)

func main() {
	world := ecs.CreateWorld("main-raylib")
	defer world.Destroy()

	world.RegisterComponents(
		&components.TransformService,
		&components.HealthService,
		&components.ColorService,
	)

	world.RegisterSystems(
		&systems.SpawnService,
		&systems.HpService,
		&systems.ColorService,
		&systems.RenderService,
	)

	world.Run(144)
}
