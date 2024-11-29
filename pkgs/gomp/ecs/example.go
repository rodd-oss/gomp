/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Transform struct {
	X, Y, Z float32
}

type Rotation struct {
	RX, RY, RZ int
}

type BulletSpawn struct{}

var _ = CreateComponent[Rotation]()
var transformComponent = CreateComponent[Transform]()
var bulletSpawnerComponent = CreateComponent[BulletSpawn]()

type TransformSystem struct {
	n int
}

func (s *TransformSystem) Init(world *ECS)    {}
func (s *TransformSystem) Destroy(world *ECS) {}
func (s *TransformSystem) Run(world *ECS) {
	s.n++
	transformComponent.Each(world, func(entity *Entity, data *Transform) {
		data.X += 1
		data.Y -= 1
		data.Z += 2
	})
}

type BulletSpawnSystem struct {
	n int
}

func (s *BulletSpawnSystem) Init(world *ECS)    {}
func (s *BulletSpawnSystem) Destroy(world *ECS) {}
func (s *BulletSpawnSystem) Run(world *ECS) {
	s.n++
	bulletSpawnerComponent.Each(world, func(entity *Entity, data *BulletSpawn) {
		tr := transformComponent.Get(entity)
		if tr == nil {
			return
		}

		bullet := world.CreateEntity("bullet")
		transformComponent.Set(bullet, *tr)
	})
}

type PlayerSpawnSystem struct{}

func (s *PlayerSpawnSystem) Init(world *ECS) {
	count := 50000
	tra := Transform{0, 1, 2}
	bs := BulletSpawn{}

	var player *Entity
	for i := 0; i < count; i++ {
		player = world.CreateEntity("Player")
		transformComponent.Set(player, tra)

		if i%2 == 0 {
			bulletSpawnerComponent.Set(player, bs)
		}
	}
}
func (s *PlayerSpawnSystem) Destroy(world *ECS) {}
func (s *PlayerSpawnSystem) Run(world *ECS)     {}
