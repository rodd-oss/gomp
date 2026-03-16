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

package entities

import (
	"gomp/examples/new-api/components"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"time"
)

type CreateSpaceSpawnerManagers struct {
	EntityManager *ecs.EntityManager
	Positions     *stdcomponents.PositionComponentManager
	Velocities    *stdcomponents.VelocityComponentManager
	SpaceSpawners *components.SpaceSpawnerComponentManager
}

func CreateSpaceSpawner(
	props CreateSpaceSpawnerManagers,
	posX, posY float32,
	velX float32,
	spawnRate time.Duration,
) ecs.Entity {
	e := props.EntityManager.Create()
	props.Positions.Create(e, stdcomponents.Position{
		XY: vectors.Vec2{
			X: posX,
			Y: posY,
		},
	})
	props.Velocities.Create(e, stdcomponents.Velocity{
		X: velX,
		Y: 0,
	})
	props.SpaceSpawners.Create(e, components.SpaceSpawnerTag{
		Cooldown:     spawnRate,
		CooldownLeft: 0,
	})

	return e
}
