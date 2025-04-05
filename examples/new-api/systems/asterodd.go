/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)]

Thank you for your support!
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/entities"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math/rand"
	"time"
)

func NewAssteroddSystem() AssteroddSystem {
	return AssteroddSystem{}
}

type AssteroddSystem struct {
	EntityManager   *ecs.EntityManager
	Positions       *stdcomponents.PositionComponentManager
	Rotations       *stdcomponents.RotationComponentManager
	Scales          *stdcomponents.ScaleComponentManager
	Velocities      *stdcomponents.VelocityComponentManager
	Sprites         *stdcomponents.SpriteComponentManager
	BoxColliders    *stdcomponents.BoxColliderComponentManager
	CircleColliders *stdcomponents.CircleColliderComponentManager
	RigidBodies     *stdcomponents.RigidBodyComponentManager

	PlayerTags       *components.PlayerTagComponentManager
	AsteroidTags     *components.AsteroidComponentManager
	BulletTags       *components.BulletTagComponentManager
	Hps              *components.HpComponentManager
	Weapons          *components.WeaponComponentManager
	SpaceshipIntents *components.SpaceshipIntentComponentManager
	SpaceSpawnerTags *components.SpaceSpawnerComponentManager
	Collisions       *stdcomponents.CollisionComponentManager
	SceneManager     *components.AsteroidSceneManagerComponentManager
	WallTags         *components.WallTagComponentManager
	SoundEffects     *components.SoundEffectsComponentManager
}

func (s *AssteroddSystem) Init() {
	entities.CreateSpaceShip(entities.CreateSpaceShipManagers{
		EntityManager:    s.EntityManager,
		Positions:        s.Positions,
		Rotations:        s.Rotations,
		Scales:           s.Scales,
		Velocities:       s.Velocities,
		Sprites:          s.Sprites,
		BoxColliders:     s.BoxColliders,
		RigidBodies:      s.RigidBodies,
		PlayerTags:       s.PlayerTags,
		Hps:              s.Hps,
		Weapons:          s.Weapons,
		SpaceshipIntents: s.SpaceshipIntents,
		SoundEffects:     s.SoundEffects,
	}, 300, 300, -44.9)
	entities.CreateSatellite(entities.CreateSatelliteManagers{
		EntityManager: s.EntityManager,
		Positions:     s.Positions,
		Rotations:     s.Rotations,
		Scales:        s.Scales,
		Velocities:    s.Velocities,
		Sprites:       s.Sprites,
		BoxColliders:  s.BoxColliders,
		RigidBodies:   s.RigidBodies,
	}, 500, 500, 0)
	entities.CreateSpaceSpawner(entities.CreateSpaceSpawnerManagers{
		EntityManager: s.EntityManager,
		Positions:     s.Positions,
		Velocities:    s.Velocities,
		SpaceSpawners: s.SpaceSpawnerTags,
	}, 16, 100, 1000, time.Millisecond*200)

	wallManager := entities.CreateWallManagers{
		EntityManager: s.EntityManager,
		Positions:     s.Positions,
		Rotations:     s.Rotations,
		Scales:        s.Scales,
		BoxColliders:  s.BoxColliders,
		Sprites:       s.Sprites,
		WallTags:      s.WallTags,
		RigidBodies:   s.RigidBodies,
	}
	entities.CreateWall(&wallManager, 0, -1000, 0, 5000, 1000)
	entities.CreateWall(&wallManager, 0, 5000, 0, 5000, 1000)
	entities.CreateWall(&wallManager, -1000, -1000, 0, 1000, 7000)
	entities.CreateWall(&wallManager, 5000, -1000, 0, 1000, 7000)

	for range 1000 {
		randPos := vectors.Vec2{
			X: float32(rand.Intn(5000)) + float32(rand.Intn(1000))/10000,
			Y: float32(rand.Intn(5000)) + float32(rand.Intn(1000))/10000,
		}

		entities.CreateBullet(entities.CreateBulletManagers{
			EntityManager:   s.EntityManager,
			Positions:       s.Positions,
			Rotations:       s.Rotations,
			Scales:          s.Scales,
			Velocities:      s.Velocities,
			CircleColliders: s.CircleColliders,
			RigidBodies:     s.RigidBodies,
			Sprites:         s.Sprites,
			BulletTags:      s.BulletTags,
			Hps:             s.Hps,
		}, randPos.X, randPos.Y, 0, 0, 0)

		// bug case
		//entities.CreateBullet(entities.CreateBulletManagers{
		//	EntityManager:   s.EntityManager,
		//	Positions:       s.Positions,
		//	Rotations:       s.Rotations,
		//	Scales:          s.Scales,
		//	Velocities:      s.Velocities,
		//	CircleColliders: s.CircleColliders,
		//	RigidBodies:     s.RigidBodies,
		//	Sprites:         s.Sprites,
		//	BulletTags:      s.BulletTags,
		//	Hps:             s.Hps,
		//}, 4.1, 4.1, 0, 100, 100)
		//entities.CreateBullet(entities.CreateBulletManagers{
		//	EntityManager:   s.EntityManager,
		//	Positions:       s.Positions,
		//	Rotations:       s.Rotations,
		//	Scales:          s.Scales,
		//	Velocities:      s.Velocities,
		//	CircleColliders: s.CircleColliders,
		//	RigidBodies:     s.RigidBodies,
		//	Sprites:         s.Sprites,
		//	BulletTags:      s.BulletTags,
		//	Hps:             s.Hps,
		//}, -4.1, -4.1, 0, -100, -100)
	}

	manager := s.EntityManager.Create()
	s.SceneManager.Create(manager, components.AsteroidSceneManager{})
}
func (s *AssteroddSystem) Run(dt time.Duration) {
	s.PlayerTags.EachEntity(func(e ecs.Entity) bool {
		intents := s.SpaceshipIntents.Get(e)

		intents.MoveUp = false
		intents.MoveDown = false
		intents.RotateLeft = false
		intents.RotateRight = false
		intents.Fire = false

		if rl.IsKeyDown(rl.KeyW) {
			intents.MoveUp = true
		}
		if rl.IsKeyDown(rl.KeyS) {
			intents.MoveDown = true
		}

		if rl.IsKeyDown(rl.KeyA) {
			intents.RotateLeft = true
		}
		if rl.IsKeyDown(rl.KeyD) {
			intents.RotateRight = true
		}

		if rl.IsKeyDown(rl.KeySpace) {
			intents.Fire = true
		}

		return true
	})

	s.SceneManager.EachEntity(func(e ecs.Entity) bool {
		sceneManager := s.SceneManager.Get(e)
		s.PlayerTags.EachEntity(func(e ecs.Entity) bool {
			playerHp := s.Hps.Get(e)
			if playerHp == nil {
				return true
			}
			sceneManager.PlayerHp = playerHp.Hp
			return false
		})
		return true
	})
}
func (s *AssteroddSystem) Destroy() {}
