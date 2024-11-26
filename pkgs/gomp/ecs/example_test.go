/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"testing"
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

func BenchmarkExample8(b *testing.B) {
	count := b.N
	if count > 10000000 {
		count = 10000000
	}
	var world = New("Main")

	world.RegisterComponents(
		&scaleComponent,
		&transformComponent,
	)

	tra := Transform{0, 1, 2}

	var player *Entity
	b.ResetTimer()
	for i := 0; i < count; i++ {
		player = world.CreateEntity("Player")
		transformComponent.Set(player, tra)
	}
	b.StopTimer()

	arr := transformComponent.Instances[&world].dense
	l := len(arr)
	if l != count {
		b.Fatal("Not equal")
	}

	for i := 0; i < l; i++ {
		arr[i].X += 1
		arr[i].Y = 0
		arr[i].Z += 2
	}

}
