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

type Bullet struct {
	HP int
}

var _ = CreateComponent[Rotation]()
var transformComponent = CreateComponent[Transform]()
var bulletSpawnerComponent = CreateComponent[BulletSpawn]()
var bulletComponent = CreateComponent[Bullet]()

type TransformSystem struct {
	n int
}

func (s *TransformSystem) Init(world *ECS)    {}
func (s *TransformSystem) Destroy(world *ECS) {}
func (s *TransformSystem) Run(world *ECS) {
	s.n++
	transformComponent.Each(world, func(entity *Entity, data Transform) {
		data.X += 1
		data.Y -= 1
		data.Z += 2

		transformComponent.Set(entity, data)
	})
}

type BulletSpawnSystem struct {
	n int
}

func (s *BulletSpawnSystem) Init(world *ECS)    {}
func (s *BulletSpawnSystem) Destroy(world *ECS) {}
func (s *BulletSpawnSystem) Run(world *ECS) {
	s.n++
	bulletSpawnerComponent.Each(world, func(entity *Entity, data BulletSpawn) {
		tr, ok := transformComponent.Get(entity)
		if !ok {
			return
		}

		bullet := world.CreateEntity("bullet")
		transformComponent.Set(bullet, tr)
		bulletData := Bullet{5}
		bulletComponent.Set(bullet, bulletData)
	})
}

type BulletSystem struct {
}

func (s *BulletSystem) Init(world *ECS)    {}
func (s *BulletSystem) Destroy(world *ECS) {}
func (s *BulletSystem) Run(world *ECS) {
	bulletComponent.Each(world, func(entity *Entity, data Bullet) {
		data.HP -= 1
		if data.HP <= 0 {
			// world.SoftDestroyEntity(entity)
		}
		bulletComponent.Set(entity, data)
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
