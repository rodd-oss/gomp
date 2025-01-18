/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Transform struct {
	X, Y int32
}

type Health struct {
	Hp, MaxHp int32
}

type Sprite struct {
	Pos  rl.Vector2
	Tint color.RGBA
}

var TransformService = ecs.CreateComponentService[Transform](TRANSFORM_ID)
var HealthService = ecs.CreateComponentService[Health](HEALTH_ID)
var ColorService = ecs.CreateComponentService[color.RGBA](COLOR_ID)
var SpriteService = ecs.CreateComponentService[Sprite](SPRITE_ID)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has Color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// Destroy system
