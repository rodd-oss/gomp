/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewVelocitySystem(
	velocities *stdcomponents.VelocityComponentManager,
	positions *stdcomponents.PositionComponentManager,
) *VelocitySystem {
	return &VelocitySystem{
		velocities: velocities,
		positions:  positions,
	}
}

type VelocitySystem struct {
	velocities *stdcomponents.VelocityComponentManager
	positions  *stdcomponents.PositionComponentManager
}

func (s *VelocitySystem) Init() {}
func (s *VelocitySystem) Run(dt time.Duration) {
	s.velocities.AllParallel(func(entity ecs.Entity, v *stdcomponents.Velocity) bool {
		position := s.positions.Get(entity)
		position.X += v.X * float32(dt.Seconds())
		position.Y += v.Y * float32(dt.Seconds())
		return true
	})
}
func (s *VelocitySystem) Destroy() {}
