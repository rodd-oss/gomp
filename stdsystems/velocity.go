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

func NewVelocitySystem() VelocitySystem {
	return VelocitySystem{}
}

type VelocitySystem struct {
	Velocities *stdcomponents.VelocityComponentManager
	Positions  *stdcomponents.PositionComponentManager
}

func (s *VelocitySystem) Init() {}
func (s *VelocitySystem) Run(dt time.Duration) {
	dtSec := float32(dt.Seconds())
	s.Velocities.AllParallel(func(entity ecs.Entity, v *stdcomponents.Velocity) bool {
		position := s.Positions.Get(entity)
		position.X += v.X * dtSec
		position.Y += v.Y * dtSec
		return true
	})
}
func (s *VelocitySystem) Destroy() {}
