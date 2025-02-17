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
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/scenes"
	"gomp/pkg/ecs"
)

func main() {
	// Initializing ECS world
	world := ecs.CreateWorld("main")
	defer world.Destroy()

	// Initializing components
	desktopComponents := gomp.NewDesktopComponents(world)
	gameComponents := components.NewGameComponents(world, &desktopComponents)

	// Initializing scenes
	sceneSet := scenes.NewSceneSet(world, &gameComponents)

	// Initializing systems
	systems := gomp.NewDesktopSystems(world, &desktopComponents)

	game := gomp.NewDesktopGame(&systems, sceneSet)
	game.CurrentSceneId = scenes.MainSceneId

	engine := gomp.NewEngine(&game)
	engine.Run(30, 165)
}
