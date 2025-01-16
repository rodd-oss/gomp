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

type transform struct {
	x, y int32
}

type health struct {
	hp, maxHp int32
}

type movement struct {
	goToX, goToY int32
}

type cameraLayer struct {
	image *ebiten.Image
	zoom  float64
}
type camera struct {
	mainLayer  cameraLayer
	debugLayer cameraLayer
}

type empty struct{}

const (
	transformId ecs.ComponentID = iota
	healthId
	colorId
	movementId
	cameraId
	destroyId
)

var transformComponentType = ecs.CreateComponentService[transform](transformId)
var healthComponentType = ecs.CreateComponentService[health](healthId)
var colorComponentType = ecs.CreateComponentService[color.RGBA](colorId)
var movementComponentType = ecs.CreateComponentService[movement](movementId)
var cameraComponentType = ecs.CreateComponentService[camera](cameraId)
var destroyComponentType = ecs.CreateComponentService[empty](destroyId)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// destroy system
