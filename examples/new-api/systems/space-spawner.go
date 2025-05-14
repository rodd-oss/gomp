/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package systems

import (
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/entities"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"math/rand"
	"time"
)

func NewSpaceSpawnerSystem() SpaceSpawnerSystem {
	return SpaceSpawnerSystem{}
}

type SpaceSpawnerSystem struct {
	EntityManager   *ecs.EntityManager
	Positions       *stdcomponents.PositionComponentManager
	SpaceSpawners   *components.SpaceSpawnerComponentManager
	Asteroids       *components.AsteroidComponentManager
	Hp              *components.HpComponentManager
	Sprites         *stdcomponents.SpriteComponentManager
	CircleColliders *stdcomponents.CircleColliderComponentManager
	Velocities      *stdcomponents.VelocityComponentManager
	Rotations       *stdcomponents.RotationComponentManager
	Scales          *stdcomponents.ScaleComponentManager
	RigidBodies     *stdcomponents.RigidBodyComponentManager
	Renderables     *stdcomponents.RenderableComponentManager
	RenderOrders    *stdcomponents.RenderOrderComponentManager
	Textures        *stdcomponents.RLTextureProComponentManager
	YSorts          *stdcomponents.YSortComponentManager
}

func (s *SpaceSpawnerSystem) Init() {}
func (s *SpaceSpawnerSystem) Run(dt time.Duration) {
	s.SpaceSpawners.EachEntity(func(e ecs.Entity) bool {
		position := s.Positions.GetUnsafe(e)
		velocity := s.Velocities.GetUnsafe(e)

		if position.XY.X > 5000 || position.XY.X < 0 {
			velocity.X = -velocity.X
		}

		spawner := s.SpaceSpawners.GetUnsafe(e)
		if spawner.CooldownLeft > 0 {
			spawner.CooldownLeft -= dt
			return true
		}

		pos := s.Positions.GetUnsafe(e)
		entities.CreateAsteroid(entities.CreateAsteroidManagers{
			EntityManager:   s.EntityManager,
			Positions:       s.Positions,
			Rotations:       s.Rotations,
			Scales:          s.Scales,
			Velocities:      s.Velocities,
			CircleColliders: s.CircleColliders,
			Sprites:         s.Sprites,
			AsteroidTags:    s.Asteroids,
			Hp:              s.Hp,
			RigidBodies:     s.RigidBodies,
			Renderables:     s.Renderables,
			YSorts:          s.YSorts,
		}, pos.XY.X, pos.XY.Y, 0, 1+rand.Float32()*2, 0, 50+rand.Float32()*100)
		spawner.CooldownLeft = spawner.Cooldown
		return true
	})
}
func (s *SpaceSpawnerSystem) Destroy() {}
