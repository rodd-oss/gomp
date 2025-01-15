/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
)

const (
	INVALID_ID ecs.ComponentID = iota
	COLOR_ID
	TRANSFORM_ID
	HEALTH_ID
	DESTROY_ID
)

type Transform struct {
	X, Y int32
}

type Health struct {
	Hp, MaxHp int32
}

var TransformService = ecs.CreateComponent[Transform](TRANSFORM_ID)
var HealthService = ecs.CreateComponent[Health](HEALTH_ID)
var ColorService = ecs.CreateComponent[color.RGBA](COLOR_ID)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has Color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// Destroy system
