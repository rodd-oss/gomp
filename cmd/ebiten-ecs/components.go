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

const (
	TRANSFORM_COMPONENT_ID ecs.ComponentID = iota
	HEALTH_COMPONENT_ID
	DESTROY_COMPONENT_ID
	SPAWN_COMPONENT_ID
	CAMERA_COMPONENT_ID
)

type transform struct {
	x, y int32
}

type health struct {
	hp    int32
	maxHp int32
	color color.RGBA
}

type cameraLayer struct {
	image *ebiten.Image
	zoom  float64
}
type camera struct {
	mainLayer  cameraLayer
	debugLayer cameraLayer
}

type destroy struct{}

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has Color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// Destroy system
