/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ClientWorld = ecs.GenericWorld[clientComponents, clientSystems]

type client struct {
	world *ClientWorld
}

type clientComponents struct {
	destroy   *ecs.ComponentManager[destroy]
	camera    *ecs.ComponentManager[camera]
	transform *ecs.ComponentManager[transform]
	health    *ecs.ComponentManager[health]
	color     *ecs.ComponentManager[color.RGBA]
}

type clientSystems struct {
	spawn   *systemSpawn
	calcHp  *systemCalcHp
	calcCol *systemCalcColor
	destroy *systemDestroyRemovedEntities
	draw    *systemDraw
}

func newGameClient() (c client) {
	// Create components
	colors := ecs.CreateComponentManager[color.RGBA](COLOR_COMPONENT_ID)
	transforms := ecs.CreateComponentManager[transform](TRANSFORM_COMPONENT_ID)
	health := ecs.CreateComponentManager[health](HEALTH_COMPONENT_ID)
	destroys := ecs.CreateComponentManager[destroy](DESTROY_COMPONENT_ID)
	cameras := ecs.CreateComponentManager[camera](CAMERA_COMPONENT_ID)

	// Create component managers

	components := clientComponents{
		color:     &colors,
		camera:    &cameras,
		health:    &health,
		destroy:   &destroys,
		transform: &transforms,
	}

	// Create systems
	systems := clientSystems{
		spawn:   new(systemSpawn),
		calcHp:  new(systemCalcHp),
		calcCol: new(systemCalcColor),
		destroy: new(systemDestroyRemovedEntities),
		draw:    new(systemDraw),
	}

	// Create world
	world := ecs.CreateGenericWorld(0, &components, &systems)

	// Register components
	world.RegisterComponents(
		components.color,
		components.camera,
		components.health,
		components.transform,
		components.destroy,
	)

	// Register update systems
	world.RegisterUpdateSystems().
		Sequential(
			systems.spawn,
			systems.calcHp,
			systems.calcCol,
			systems.destroy,
		)

	// Register draw systems
	world.RegisterDrawSystems().
		Sequential(
			systems.draw,
		)

	newClient := client{
		world: &world,
	}

	return newClient
}

func (c *client) Update() error {
	return c.world.RunUpdateSystems()
}

func (c *client) Draw(screen *ebiten.Image) {
	c.world.RunDrawSystems(screen)
}

func (c *client) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
