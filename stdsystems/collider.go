/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- rpecb Donated 500 RUB

Thank you for your support!
*/

package stdsystems

import (
	"github.com/negrel/assert"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"gomp/vectors"
	"time"
)

func NewColliderSystem() ColliderSystem {
	return ColliderSystem{}
}

type ColliderSystem struct {
	EntityManager                      *ecs.EntityManager
	Positions                          *stdcomponents.PositionComponentManager
	Scales                             *stdcomponents.ScaleComponentManager
	Rotations                          *stdcomponents.RotationComponentManager
	Velocities                         *stdcomponents.VelocityComponentManager
	GenericColliders                   *stdcomponents.GenericColliderComponentManager
	BoxColliders                       *stdcomponents.BoxColliderComponentManager
	CircleColliders                    *stdcomponents.CircleColliderComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	AABB                               *stdcomponents.AABBComponentManager

	numWorkers             int
	accAABB                [][]ecs.Entity
	accGenericColliders    [][]ecs.Entity
	accColliderSleepCreate [][]ecs.Entity
	accColliderSleepDelete [][]ecs.Entity
	Engine                 *core.Engine
}

func (s *ColliderSystem) Init() {
	s.numWorkers = s.Engine.Pool().NumWorkers()
	s.accAABB = make([][]ecs.Entity, s.numWorkers)
	s.accGenericColliders = make([][]ecs.Entity, s.numWorkers)

	s.accColliderSleepCreate = make([][]ecs.Entity, s.numWorkers)
	s.accColliderSleepDelete = make([][]ecs.Entity, s.numWorkers)
}
func (s *ColliderSystem) Run(dt time.Duration) {
	for i := range s.accAABB {
		s.accAABB[i] = s.accAABB[i][:0]
	}
	for i := range s.accGenericColliders {
		s.accGenericColliders[i] = s.accGenericColliders[i][:0]
	}
	for i := range s.accColliderSleepCreate {
		s.accColliderSleepCreate[i] = s.accColliderSleepCreate[i][:0]
	}
	for i := range s.accColliderSleepDelete {
		s.accColliderSleepDelete[i] = s.accColliderSleepDelete[i][:0]
	}
	s.BoxColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		if !s.GenericColliders.Has(entity) {
			s.accGenericColliders[workerId] = append(s.accGenericColliders[workerId], entity)
		}
		if !s.AABB.Has(entity) {
			s.accAABB[workerId] = append(s.accAABB[workerId], entity)
		}
	})
	s.CircleColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		if !s.GenericColliders.Has(entity) {
			s.accGenericColliders[workerId] = append(s.accGenericColliders[workerId], entity)
		}
		if !s.AABB.Has(entity) {
			s.accAABB[workerId] = append(s.accAABB[workerId], entity)
		}
	})
	for i := range s.accAABB {
		a := s.accAABB[i]
		for _, entity := range a {
			s.AABB.Create(entity, stdcomponents.AABB{})
		}
	}
	for i := range s.accGenericColliders {
		a := s.accGenericColliders[i]
		for _, entity := range a {
			s.GenericColliders.Create(entity, stdcomponents.GenericCollider{})
		}
	}

	s.BoxColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		boxCollider := s.BoxColliders.GetUnsafe(entity)

		genCollider := s.GenericColliders.GetUnsafe(entity)

		genCollider.Layer = boxCollider.Layer
		genCollider.Mask = boxCollider.Mask
		genCollider.Offset.X = boxCollider.Offset.X
		genCollider.Offset.Y = boxCollider.Offset.Y
		genCollider.Shape = stdcomponents.BoxColliderShape
		genCollider.AllowSleep = boxCollider.AllowSleep
	})

	s.BoxColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		boxCollider := s.BoxColliders.GetUnsafe(entity)
		assert.NotNil(boxCollider)

		position := s.Positions.GetUnsafe(entity)
		assert.NotNil(position)

		scale := s.Scales.GetUnsafe(entity)
		assert.NotNil(scale)

		rotation := s.Rotations.GetUnsafe(entity)
		assert.NotNil(rotation)

		aabb := s.AABB.GetUnsafe(entity)
		assert.NotNil(aabb)

		a := boxCollider.WH
		b := vectors.Vec2{X: 0, Y: boxCollider.WH.Y}
		c := vectors.Vec2{X: 0, Y: 0}
		d := vectors.Vec2{X: boxCollider.WH.X, Y: 0}

		c = c.Sub(boxCollider.Offset).Rotate(rotation.Angle)
		a = a.Sub(boxCollider.Offset).Rotate(rotation.Angle)
		b = b.Sub(boxCollider.Offset).Rotate(rotation.Angle)
		d = d.Sub(boxCollider.Offset).Rotate(rotation.Angle)

		aabb.Min = vectors.Vec2{X: min(b.X, c.X, a.X, d.X), Y: min(b.Y, c.Y, a.Y, d.Y)}.Mul(scale.XY)
		aabb.Max = vectors.Vec2{X: max(b.X, c.X, a.X, d.X), Y: max(b.Y, c.Y, a.Y, d.Y)}.Mul(scale.XY)

		aabb.Min = position.XY.Add(aabb.Min)
		aabb.Max = position.XY.Add(aabb.Max)
	})

	s.CircleColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		circleCollider := s.CircleColliders.GetUnsafe(entity)
		assert.NotNil(circleCollider)

		genCollider := s.GenericColliders.GetUnsafe(entity)
		assert.NotNil(genCollider)

		genCollider.Layer = circleCollider.Layer
		genCollider.Mask = circleCollider.Mask
		genCollider.Offset.X = circleCollider.Offset.X
		genCollider.Offset.Y = circleCollider.Offset.Y
		genCollider.Shape = stdcomponents.CircleColliderShape
		genCollider.AllowSleep = circleCollider.AllowSleep
	})

	s.CircleColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		circleCollider := s.CircleColliders.GetUnsafe(entity)
		assert.NotNil(circleCollider)

		position := s.Positions.GetUnsafe(entity)
		assert.NotNil(position)

		scale := s.Scales.GetUnsafe(entity)
		assert.NotNil(scale)

		aabb := s.AABB.GetUnsafe(entity)
		assert.NotNil(aabb)

		offset := circleCollider.Offset.Mul(scale.XY)
		scaledRadius := scale.XY.Scale(circleCollider.Radius)
		aabb.Min = position.XY.Add(offset).Sub(scaledRadius)
		aabb.Max = position.XY.Add(offset).Add(scaledRadius)
	})

	s.GenericColliders.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		genCollider := s.GenericColliders.GetUnsafe(entity)
		if genCollider.AllowSleep {
			shouldSleep := true
			velocity := s.Velocities.GetUnsafe(entity)
			if velocity != nil {
				if velocity.Vec2().LengthSquared() != 0 {
					shouldSleep = false
				}
			}
			isSleeping := s.ColliderSleepStateComponentManager.GetUnsafe(entity)
			if shouldSleep {
				if isSleeping == nil {
					s.accColliderSleepCreate[workerId] = append(s.accColliderSleepCreate[workerId], entity)
				}
			} else {
				if isSleeping != nil {
					s.accColliderSleepDelete[workerId] = append(s.accColliderSleepDelete[workerId], entity)
				}
			}
		}
	})
	for i := range s.accColliderSleepCreate {
		a := s.accColliderSleepCreate[i]
		for _, entity := range a {
			s.ColliderSleepStateComponentManager.Create(entity, stdcomponents.ColliderSleepState{})
		}
	}
	for i := range s.accColliderSleepDelete {
		a := s.accColliderSleepDelete[i]
		for _, entity := range a {
			s.ColliderSleepStateComponentManager.Delete(entity)
		}
	}

}
func (s *ColliderSystem) Destroy() {}
