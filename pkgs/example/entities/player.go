/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package entities

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/components"
)

var PlayerPhysics = components.PhysicsComponent(components.PhysicsData{})
var PlayerRender = components.RenderComponent(components.RenderData{})

var Player = engine.CreateEntity(
	PlayerPhysics,
	PlayerRender,
)
