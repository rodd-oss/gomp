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
	Destroy   *ecs.ComponentManager[destroy]
	Camera    *ecs.ComponentManager[camera]
	Transform *ecs.ComponentManager[transform]
	Health    *ecs.ComponentManager[health]
	Color     *ecs.ComponentManager[color.RGBA]
}

type clientSystems struct {
	Spawn   *systemSpawn
	CalcHp  *systemCalcHp
	CalcCol *systemCalcColor
	Destroy *systemDestroyRemovedEntities
	Draw    *systemDraw
}

func newGameClient() (c client) {
	// TODO: move initializing components with reflect inside CreateGenericWorld() function?
	// Create component managers
	components := clientComponents{
		Color:     ecs.CreateComponentManager[color.RGBA](COLOR_COMPONENT_ID),
		Camera:    ecs.CreateComponentManager[camera](CAMERA_COMPONENT_ID),
		Health:    ecs.CreateComponentManager[health](HEALTH_COMPONENT_ID),
		Destroy:   ecs.CreateComponentManager[destroy](DESTROY_COMPONENT_ID),
		Transform: ecs.CreateComponentManager[transform](TRANSFORM_COMPONENT_ID),
	}

	// Create systems
	systems := new(clientSystems)

	// Create world and register components and systems
	world := ecs.CreateGenericWorld(0, &components, systems)

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
