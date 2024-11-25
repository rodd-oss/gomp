/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "testing"

type Transform struct {
	X, Y, Z float32
}

type Rotation struct {
	RX, RY, RZ int
}

type Scale struct {
	Value float32
}

func TestExample(t *testing.T) {
	var ecs = New()

	var transformComponent = CreateComponent[Transform](&ecs)
	var _ = CreateComponent[Rotation](&ecs)
	var scaleComponent = CreateComponent[Scale](&ecs)

	var playerEntity = ecs.CreateEntity()

	transformComponent.Set(playerEntity, Transform{0, 1, 2})
	scaleComponent.Set(playerEntity, Scale{0})
}
