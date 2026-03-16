/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"time"
)

func NewDampingSystem() DampingSystem {
	return DampingSystem{}
}

type DampingSystem struct {
	Velocities  *stdcomponents.VelocityComponentManager
	Positions   *stdcomponents.PositionComponentManager
	RigidBodies *stdcomponents.RigidBodyComponentManager
}

const (
	dampingFactor float32 = 0.98
)

func (s *DampingSystem) Init() {}

func (s *DampingSystem) Run(dt time.Duration) {
	s.Velocities.ProcessEntities(func(e ecs.Entity, _ worker.WorkerId) {
		velocity := s.Velocities.GetUnsafe(e)
		rigidbody := s.RigidBodies.GetUnsafe(e)

		if rigidbody != nil && !rigidbody.IsStatic {
			velocity.X *= dampingFactor / (config.TickRate * float32(dt.Seconds()))
			velocity.Y *= dampingFactor / (config.TickRate * float32(dt.Seconds()))
			if velocity.X < 0.1 && velocity.X > -0.1 {
				velocity.X = 0
			}
			if velocity.Y < 0.1 && velocity.Y > -0.1 {
				velocity.Y = 0
			}
		}
	})
}

func (s *DampingSystem) Destroy() {}
