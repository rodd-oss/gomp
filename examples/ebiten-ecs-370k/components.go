/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package main

import (
	ecs2 "gomp/pkgs/ecs"
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
	transformId ecs2.ComponentID = iota
	healthId
	colorId
	movementId
	cameraId
	destroyId
)

var transformComponentType = ecs2.CreateComponentService[transform](transformId)
var healthComponentType = ecs2.CreateComponentService[health](healthId)
var colorComponentType = ecs2.CreateComponentService[color.RGBA](colorId)
var movementComponentType = ecs2.CreateComponentService[movement](movementId)
var cameraComponentType = ecs2.CreateComponentService[camera](cameraId)
var destroyComponentType = ecs2.CreateComponentService[empty](destroyId)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// destroy system
