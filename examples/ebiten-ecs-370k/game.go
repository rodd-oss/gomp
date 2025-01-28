/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	ecs2 "gomp/pkg/ecs"
)

type game struct {
	world            *ecs2.World
	cameraComponents *ecs2.ComponentManager[camera]
	op               *ebiten.DrawImageOptions
}

func newGame() game {
	world := ecs2.CreateWorld("1 mil pixel")

	world.RegisterComponentServices(
		&destroyComponentType,
		&cameraComponentType,
		&transformComponentType,
		&healthComponentType,
		&colorComponentType,
		&movementComponentType,
	)

	// world.RegisterSystems().
	// 	Sequential(
	// 		new(spawnSystem),
	// 		new(hpSystem),
	// 		new(colorSystem),
	// 		new(destroySystem),
	// 		new(cameraSystem),
	// 	)

	newGame := game{
		world:            &world,
		cameraComponents: cameraComponentType.GetManager(&world),
		op:               new(ebiten.DrawImageOptions),
	}

	return newGame
}

func (g *game) Update() error {
	err := g.world.FixedUpdate()
	if err != nil {
		return err
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	var mainCamera *camera

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
