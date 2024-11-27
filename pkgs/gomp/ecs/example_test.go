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

func BenchmarkEntityUpdate(b *testing.B) {
	b.ReportAllocs()
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
	b.Log("Creating", count, "entities in", time.Since(start))

	b.ResetTimer()
	for range b.N {
		transformComponent.Each(&world, func(data *Transform) {
			data.X += 1
			data.Y = 0
			data.Z += 2
		})
	}
}
