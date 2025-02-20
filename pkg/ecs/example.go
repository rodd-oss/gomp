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

const (
	invalidID ComponentID = iota
	rotationID
	transformID
	bulletID
	bulletSpawnID
)

var _ = CreateComponentService[Rotation](rotationID)
var transformComponent = CreateComponentService[Transform](transformID)
var bulletSpawnerComponent = CreateComponentService[BulletSpawn](bulletID)
var bulletComponent = CreateComponentService[Bullet](bulletSpawnID)

type TransformSystem struct {
	n         int
	transform *ComponentManager[Transform]
}

func (s *TransformSystem) Init(world *EntityManager) {
	s.transform = transformComponent.GetManager(world)
}
func (s *TransformSystem) Destroy(world *EntityManager) {}
func (s *TransformSystem) Run(world *EntityManager) {
	s.n++
	for _, t := range s.transform.All {
		t.X += 1
		t.Y -= 1
		t.Z += 2
	}
}

type BulletSpawnSystem struct {
	n             int
	bulletSpawner *ComponentManager[BulletSpawn]
	transform     *ComponentManager[Transform]
	bullet        *ComponentManager[Bullet]
}

func (s *BulletSpawnSystem) Init(world *EntityManager) {
	s.bulletSpawner = bulletSpawnerComponent.GetManager(world)
	s.transform = transformComponent.GetManager(world)
	s.bullet = bulletComponent.GetManager(world)
}
func (s *BulletSpawnSystem) Destroy(world *EntityManager) {}
func (s *BulletSpawnSystem) Run(world *EntityManager) {
	s.n++

	var bulletData Bullet
	bulletData.HP = 5

	for id := range s.bulletSpawner.All {
		tr := s.transform.Get(id)
		if tr == nil {
			continue
		}

		newBullet := world.Create()
		s.transform.Create(newBullet, *tr)
		s.bullet.Create(newBullet, bulletData)
	}
}

type BulletSystem struct {
	bullet *ComponentManager[Bullet]
}

func (s *BulletSystem) Init(world *EntityManager) {
	s.bullet = bulletComponent.GetManager(world)
}
func (s *BulletSystem) Destroy(world *EntityManager) {}
func (s *BulletSystem) Run(world *EntityManager) {
	for entId, b := range s.bullet.All {
		b.HP -= 1
		if b.HP <= 0 {
			world.Delete(entId)
		}
	}
}

type PlayerSpawnSystem struct {
	bulletSpawner *ComponentManager[BulletSpawn]
	transform     *ComponentManager[Transform]
}

func (s *PlayerSpawnSystem) Init(world *EntityManager) {
	s.bulletSpawner = bulletSpawnerComponent.GetManager(world)
	s.transform = transformComponent.GetManager(world)

	count := 100_000
	tra := Transform{0, 1, 2}
	bs := BulletSpawn{}

	var player Entity
	for i := 0; i < count; i++ {
		player = world.Create()
		s.transform.Create(player, tra)

		if i%2 == 0 {
			s.bulletSpawner.Create(player, bs)
		}
	}
}
func (s *PlayerSpawnSystem) Destroy(world *EntityManager) {}
func (s *PlayerSpawnSystem) Run(world *EntityManager)     {}
