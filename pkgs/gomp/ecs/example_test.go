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

func BenchmarkEntityUpdate(b *testing.B) {
	b.ReportAllocs()
	count := 10_000_000

	var world = New("Main")

	world.RegisterComponents(
		&scaleComponent,
		&transformComponent,
	)

	tra := Transform{0, 1, 2}
	sc := Scale{1}

	var player *Entity
	for i := 0; i < count; i++ {
		player = world.CreateEntity("Player")
		if i%2 == 0 {
			transformComponent.Set(player, tra)
		}
		if i%10 == 0 {
			scaleComponent.Set(player, sc)
		}
	}

	b.ResetTimer()
	for range b.N {
		transformComponent.Each(&world, func(entity *Entity, data *Transform) {
			scale := scaleComponent.Get(entity)
			if scale == nil {
				return
			}

			data.X += 1
			data.Y -= 1
			data.Z += 2
		})

		// scaleComponent.Each(&world, func(entity *Entity, data *Scale) {
		// 	tr := transformComponent.Get(entity)
		// 	if tr == nil {
		// 		return
		// 	}

		// 	tr.X += 1
		// 	tr.Y -= 1
		// 	tr.Z += 2
		// })
	}
}

func BenchmarkCreateWorld(b *testing.B) {
	b.ReportAllocs()
	count := 1_000_000

	b.ResetTimer()
	for range b.N {
		b.StopTimer()
		var world = New("Main")

		world.RegisterComponents(
			&transformComponent,
		)

		tra := Transform{0, 1, 2}

		var player *Entity
		b.StartTimer()
		for i := 0; i < count; i++ {
			player = world.CreateEntity("Player")
			transformComponent.Set(player, tra)
		}
	}
}
