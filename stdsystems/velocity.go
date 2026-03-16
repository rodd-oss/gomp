/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"github.com/negrel/assert"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"math"
	"runtime"
	"time"
)

func NewVelocitySystem() VelocitySystem {
	return VelocitySystem{}
}

type VelocitySystem struct {
	Velocities  *stdcomponents.VelocityComponentManager
	Positions   *stdcomponents.PositionComponentManager
	RigidBodies *stdcomponents.RigidBodyComponentManager

	numWorkers int
	Engine     *core.Engine
}

func (s *VelocitySystem) Init() {
	s.numWorkers = runtime.NumCPU() - 2
}

func (s *VelocitySystem) Run(dt time.Duration) {
	dtSec := float32(dt.Seconds())

	s.Velocities.ProcessEntities(func(e ecs.Entity, workerId worker.WorkerId) {
		velocity := s.Velocities.GetUnsafe(e)
		assert.True(s.isVelocityValid(velocity))

		position := s.Positions.GetUnsafe(e)
		assert.True(s.isPositionValid(position))

		position.XY.X += velocity.X * dtSec
		position.XY.Y += velocity.Y * dtSec
	})
}

func (s *VelocitySystem) Destroy() {}

func (s *VelocitySystem) isVelocityValid(velocity *stdcomponents.Velocity) bool {
	if velocity == nil {
		return false
	}

	// Convert to float64
	x := float64(velocity.X)
	y := float64(velocity.Y)

	if math.IsInf(x, 1) || math.IsInf(x, -1) {
		return false
	}

	if math.IsInf(y, 1) || math.IsInf(y, -1) {
		return false
	}

	if math.IsNaN(x) || math.IsNaN(y) {
		return false
	}

	return true
}

func (s *VelocitySystem) isPositionValid(position *stdcomponents.Position) bool {
	if position == nil {
		return false
	}

	// Convert to float64
	x := float64(position.XY.X)
	y := float64(position.XY.Y)

	if math.IsInf(x, 1) || math.IsInf(x, -1) {
		return false
	}

	if math.IsInf(y, 1) || math.IsInf(y, -1) {
		return false
	}

	if math.IsNaN(x) || math.IsNaN(y) {
		return false
	}

	return true
}
