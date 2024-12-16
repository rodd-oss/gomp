/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"testing"
)

func BenchmarkSystems(b *testing.B) {
	b.ReportAllocs()
	var world = New("Main")

	world.RegisterComponents(
		&bulletSpawnerComponent,
		&transformComponent,
		&bulletComponent,
	)

	world.RegisterSystems().
		Parallel(
			new(PlayerSpawnSystem),
			new(BulletSpawnSystem),
		).
		Sequential(
			new(BulletSystem),
			new(TransformSystem),
		)

	b.ResetTimer()
	for range b.N {
		world.RunSystems()
	}
}

// func BenchmarkEntityUpdate(b *testing.B) {
// 	b.ReportAllocs()
// 	count := 1_000_000

// 	var world = New("Main")

// 	world.RegisterComponents(
// 		&bulletSpawnerComponent,
// 		&transformComponent,
// 	)

// 	world.RegisterSystems().
// 		Parallel(
// 			new(TransformSystem),
// 			new(BulletSpawnSystem),
// 		).
// 		Sequential(
// 			new(BulletSpawnSystem),
// 			new(TransformSystem),
// 		)

// 	tra := Transform{0, 1, 2}
// 	sc := BulletSpawn{}

// 	var player *Entity
// 	transform := transformComponent.Instances(&world)
// 	bullet := bulletSpawnerComponent.Instances(&world)
// 	for i := 0; i < count; i++ {
// 		player = world.CreateEntity("Player")
// 		if i%2 == 0 {
// 			transform.Set(player.ID, tra)
// 		}
// 		if i%10 == 0 {
// 			bullet.Set(player.ID, sc)
// 		}
// 	}

// 	b.ResetTimer()
// 	for range b.N {
// 		transform := transformComponent.Instances(&world)

// 		for _, v := range transform.All() {
// 			v.X += 1
// 			v.Y -= 1
// 			v.Z += 2
// 		}
// 	}
// }

// func BenchmarkCreateWorld(b *testing.B) {
// 	b.ReportAllocs()
// 	count := 1_000_000

// 	b.ResetTimer()
// 	for range b.N {
// 		b.StopTimer()
// 		var world = New("Main")

// 		world.RegisterComponents(
// 			&transformComponent,
// 		)

// 		tra := Transform{0, 1, 2}

// 		var player *Entity
// 		b.StartTimer()
// 		transform := transformComponent.Instances(&world)
// 		for i := 0; i < count; i++ {
// 			player = world.CreateEntity("Player")
// 			transform.Set(player.ID, tra)
// 		}
// 	}
// }

// func BenchmarkEntityCreate(b *testing.B) {
// 	var world = New("Main")
// 	world.RegisterComponents(
// 		&bulletSpawnerComponent,
// 		&transformComponent,
// 	)

// 	b.ResetTimer()
// 	b.ReportAllocs()
// 	var tra Transform

// 	transform := transformComponent.Instances(&world)

// 	for i := 0; i < b.N; i++ {
// 		tra.X = float32(i)
// 		tra.Y = float32(i + 1)
// 		tra.Z = float32(i + 4)
// 		player := world.CreateEntity("Player")
// 		transform.Set(player.ID, tra)
// 	}
// }

// func TestEntityUpdate(t *testing.T) {
// 	var world = New("Main")
// 	world.RegisterComponents(
// 		&bulletSpawnerComponent,
// 		&transformComponent,
// 	)

// 	transform := transformComponent.Instances(&world)

// 	var cases []EntityID
// 	for i := 0; i < 10_000_000; i++ {
// 		tra := Transform{float32(i), float32(-i), 2}
// 		player := world.CreateEntity("Player")
// 		transform.Set(player.ID, tra)
// 		cases = append(cases, player.ID)
// 	}
// 	// check
// 	for i, id := range cases {
// 		tra := Transform{float32(i), float32(-i), 2}

// 		if id == 1000000 {
// 			t.Log("test", id)
// 		}
// 		e, ok := world.Entities.Get(id)
// 		if !ok {
// 			t.Fatalf("not found entity with id: %v", id)
// 		}

// 		if e.ID != id {
// 			t.Fatalf("want: %v, got: %v", id, e.ID)
// 		}

// 		entity, ok := transform.Get(e.ID)
// 		if !ok {
// 			t.Fatalf("not found component %v", id)
// 		}

// 		if entity != tra {
// 			t.Fatalf("want: %v, got: %v", tra, entity)
// 		}
// 	}
// }
