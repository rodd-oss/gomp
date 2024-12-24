/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/pkgs/gomp/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
	world            *ecs.World
	cameraComponents ecs.WorldComponents[camera]
	op               *ebiten.DrawImageOptions
}

func newGame() (newGame game) {
	world := ecs.New("1 mil pixel")

	world.RegisterComponentTypes(
		&destroyComponentType,
		&cameraComponentType,
		&transformComponentType,
		&healthComponentType,
		&colorComponentType,
		&movementComponentType,
	)

	world.RegisterUpdateSystems().
		Sequential(
			new(systemSpawn),
			new(systemCalcHp),
			new(systemCalcColor),
			new(systemDestroyRemovedEntities),
		)

	world.RegisterDrawSystems().
		Sequential(
			new(systemDraw),
		)

	newGame.world = &world
	newGame.cameraComponents = cameraComponentType.Instances(&world)
	newGame.op = new(ebiten.DrawImageOptions)

	return newGame
}

func (g *game) Update() error {
	return g.world.RunUpdateSystems()
}

func (g *game) Draw(screen *ebiten.Image) {
	var mainCamera *camera

	g.world.RunDrawSystems()

	g.cameraComponents.AllData(func(c *camera) bool {
		mainCamera = c
		return false
	})

	if mainCamera == nil {
		return
	}

	g.op.GeoM.Reset()
	g.op.GeoM.Scale(mainCamera.mainLayer.zoom, mainCamera.mainLayer.zoom)
	screen.DrawImage(mainCamera.mainLayer.image, g.op)

	g.op.GeoM.Reset()
	g.op.GeoM.Scale(mainCamera.debugLayer.zoom, mainCamera.debugLayer.zoom)
	screen.DrawImage(mainCamera.debugLayer.image, g.op)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
