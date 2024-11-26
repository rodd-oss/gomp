/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"testing"
	"time"
)

type Transform struct {
	X, Y, Z float32
}

var transformComponent = CreateComponent[Transform]()

type Rotation struct {
	RX, RY, RZ int
}

var _ = CreateComponent[Rotation]()

type Scale struct {
	Value float32
}

var scaleComponent = CreateComponent[Scale]()

func TestExample(t *testing.T) {
	count := 10000000
	var world = New("Main")

	world.RegisterComponents(
		&scaleComponent,
		&transformComponent,
	)

	tra := Transform{0, 1, 2}

	var player *Entity
	start := time.Now()
	for i := 0; i < count; i++ {
		player = world.CreateEntity("Player")
		transformComponent.Set(player, tra)
	}
	t.Log("Creating", count, "entities in", time.Since(start))

	arr := transformComponent.Instances[&world].dense
	l := len(arr)
	if l != count {
		t.Fatal("Not equal", l, count)
	}

	start = time.Now()
	for i := 0; i < l; i++ {
		arr[i].X += 1
		arr[i].Y = 0
		arr[i].Z += 2
	}
	t.Log("Updating", l, "components in", time.Since(start))
}
