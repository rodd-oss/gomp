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
	n         int
	transform WorldComponents[Transform]
}

func (s *TransformSystem) Init(world *ECS) {
	s.transform = transformComponent.Instances(world)
}
func (s *TransformSystem) Destroy(world *ECS) {}
func (s *TransformSystem) Run(world *ECS) {
	s.n++
	for _, t := range s.transform.All() {
		t.X += 1
		t.Y -= 1
		t.Z += 2
	}
}

type BulletSpawnSystem struct {
	n             int
	bulletSpawner WorldComponents[BulletSpawn]
	transform     WorldComponents[Transform]
	bullet        WorldComponents[Bullet]
}

func (s *BulletSpawnSystem) Init(world *ECS) {
	s.bulletSpawner = bulletSpawnerComponent.Instances(world)
	s.transform = transformComponent.Instances(world)
	s.bullet = bulletComponent.Instances(world)
}
func (s *BulletSpawnSystem) Destroy(world *ECS) {}
func (s *BulletSpawnSystem) Run(world *ECS) {
	s.n++

	var bulletData Bullet

	for id := range s.bulletSpawner.All() {
		tr, ok := s.transform.Get(id)
		if !ok {
			continue
		}

		newBullet := world.CreateEntity("bullet")
		s.transform.Set(newBullet, tr)
		bulletData.HP = 5
		s.bullet.Set(newBullet, bulletData)
	}
}

type BulletSystem struct {
	bullet WorldComponents[Bullet]
}

func (s *BulletSystem) Init(world *ECS) {
	s.bullet = bulletComponent.Instances(world)
}
func (s *BulletSystem) Destroy(world *ECS) {}
func (s *BulletSystem) Run(world *ECS) {
	for entId, b := range s.bullet.All() {
		b.HP -= 1
		if b.HP <= 0 {
			world.SoftDestroyEntity(entId)
		}
	}
}

type PlayerSpawnSystem struct {
	bulletSpawner WorldComponents[BulletSpawn]
	transform     WorldComponents[Transform]
}

func (s *PlayerSpawnSystem) Init(world *ECS) {
	s.bulletSpawner = bulletSpawnerComponent.Instances(world)
	s.transform = transformComponent.Instances(world)

	count := 100_000
	tra := Transform{0, 1, 2}
	bs := BulletSpawn{}

	var player EntityID
	for i := 0; i < count; i++ {
		player = world.CreateEntity("Player")
		s.transform.Set(player, tra)

		if i%2 == 0 {
			s.bulletSpawner.Set(player, bs)
		}
	}
}
func (s *PlayerSpawnSystem) Destroy(world *ECS) {}
func (s *PlayerSpawnSystem) Run(world *ECS)     {}
