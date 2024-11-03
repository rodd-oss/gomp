/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package entities

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/components"

	"github.com/jakecoffman/cp/v2"
)

var PlayerPhysics = components.PhysicsComponent(components.PhysicsData{
	Body: cp.NewKinematicBody(),
})

var Player = engine.CreateEntity(
	PlayerPhysics,
)
