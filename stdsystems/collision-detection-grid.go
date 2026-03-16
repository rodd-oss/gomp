/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- SoundOfTheWind Donated 100 RUB
<- Anonymous Donated 500 RUB

Thank you for your support!
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"time"
)

func NewCollisionDetectionGridSystem() CollisionDetectionGridSystem {
	return CollisionDetectionGridSystem{
		cellSizeX:      192,
		cellSizeY:      192,
		spatialBuckets: make(map[stdcomponents.SpatialHash][]ecs.Entity, 32),
	}
}

type CollisionDetectionGridSystem struct {
	EntityManager   *ecs.EntityManager
	Positions       *stdcomponents.PositionComponentManager
	Scales          *stdcomponents.ScaleComponentManager
	GenericCollider *stdcomponents.GenericColliderComponentManager
	BoxColliders    *stdcomponents.BoxColliderComponentManager
	Collisions      *stdcomponents.CollisionComponentManager
	SpatialIndex    *stdcomponents.SpatialHashComponentManager

	cellSizeX int
	cellSizeY int

	// Cache
	activeCollisions  map[CollisionPair]ecs.Entity // Maps collision pairs to proxy entities
	currentCollisions map[CollisionPair]struct{}
	spatialBuckets    map[stdcomponents.SpatialHash][]ecs.Entity
	entityToCell      map[ecs.Entity]stdcomponents.SpatialHash
	aabbs             map[ecs.Entity]aabb
}

type aabb struct {
	Left, Right, Top, Bottom float32
}

type CollisionPair struct {
	E1, E2 ecs.Entity
}

func (c CollisionPair) Normalize() CollisionPair {
	if c.E1 > c.E2 {
		return CollisionPair{c.E2, c.E1}
	}
	return c
}

func (s *CollisionDetectionGridSystem) Init() {
	s.activeCollisions = make(map[CollisionPair]ecs.Entity)
}

type CollisionEvent struct {
	entityA, entityB ecs.Entity
	position         vectors.Vec2
	normal           vectors.Vec2
	depth            float32
}

func (s *CollisionDetectionGridSystem) Run(dt time.Duration) {
	//if len(s.entityToCell) < s.GenericCollider.Len() {
	//	s.entityToCell = make(map[ecs.Entity]stdcomponents.SpatialHash, s.GenericCollider.Len())
	//}
	//// Reuse spatialBuckets to reduce allocations
	//for k := range s.spatialBuckets {
	//	delete(s.spatialBuckets, k)
	//}
	//s.currentCollisions = make(map[CollisionPair]struct{})
	//
	//// Build spatial buckets and entity-to-cell map
	//s.GenericCollider.EachEntity()(func(entity ecs.Entity) bool {
	//	position := s.Positions.GetUnsafe(entity)
	//	scale := s.Scales.GetUnsafe(entity)
	//
	//	collider := s.GenericCollider.GetUnsafe(entity)
	//	cellX := int(position.XY.X-(collider.Offset.X*scale.XY.X)) / s.cellSizeX
	//	cellY := int(position.XY.Y-(collider.Offset.Y*scale.XY.Y)) / s.cellSizeY
	//	cell := stdcomponents.SpatialHash{X: cellX, Y: cellY}
	//	s.entityToCell[entity] = cell
	//	s.spatialBuckets[cell] = append(s.spatialBuckets[cell], entity)
	//	return true
	//})
	//
	//// Precompute AABBs for box colliders
	//s.aabbs = make(map[ecs.Entity]aabb, s.BoxColliders.Len())
	//s.BoxColliders.EachEntity()(func(entity ecs.Entity) bool {
	//	position := s.Positions.GetUnsafe(entity)
	//	collider := s.BoxColliders.GetUnsafe(entity)
	//	scale := s.Scales.GetUnsafe(entity)
	//	newAABB := aabb{
	//		Left:   position.XY.X - (collider.Offset.X * scale.XY.X),
	//		Right:  position.XY.X + (collider.WH.X-collider.Offset.X)*scale.XY.X,
	//		Top:    position.XY.Y - (collider.Offset.Y * scale.XY.Y),
	//		Bottom: position.XY.Y + (collider.WH.Y-collider.Offset.Y)*scale.XY.Y,
	//	}
	//	s.aabbs[entity] = newAABB
	//	return true
	//})
	//
	//// Create collision channel
	//collisionChan := make(chan CollisionEvent, 4096)
	//doneChan := make(chan struct{})
	//
	//// Start result collector
	//go func() {
	//	for event := range collisionChan {
	//		pair := CollisionPair{event.entityA, event.entityB}.Normalize()
	//		s.currentCollisions[pair] = struct{}{}
	//
	//		if _, exists := s.activeCollisions[pair]; !exists {
	//			proxy := s.EntityManager.Create()
	//			s.Collisions.Create(proxy, stdcomponents.Collision{E1: pair.E1, E2: pair.E2, State: stdcomponents.CollisionStateEnter})
	//			s.Positions.Create(proxy, stdcomponents.Position{
	//				XY: vectors.Vec2{
	//					X: event.position.X,
	//					Y: event.position.Y,
	//				},
	//			})
	//			s.activeCollisions[pair] = proxy
	//		} else {
	//			proxy := s.activeCollisions[pair]
	//			s.Collisions.GetUnsafe(proxy).State = stdcomponents.CollisionStateStay
	//			s.Positions.GetUnsafe(proxy).XY.X = event.position.X
	//			s.Positions.GetUnsafe(proxy).XY.Y = event.position.Y
	//		}
	//	}
	//	close(doneChan)
	//}()
	//
	//entities := s.GenericCollider.RawEntities(make([]ecs.Entity, 0, s.GenericCollider.Len()))
	//
	//// Worker pool setup
	//var wg sync.WaitGroup
	//maxNumWorkers := runtime.NumCPU() - 2
	//entitiesLength := len(entities)
	//// get minimum 1 worker for small amount of entities, and maximum maxNumWorkers for a lot entities
	//numWorkers := max(min(entitiesLength/32, maxNumWorkers), 1)
	//chunkSize := entitiesLength / numWorkers
	//
	//wg.Add(numWorkers)
	//
	//for i := 0; i < numWorkers; i++ {
	//	startIndex := i * chunkSize
	//	endIndex := startIndex + chunkSize - 1
	//	if i == numWorkers-1 { // have to set endIndex to entites lenght, if last worker
	//		endIndex = entitiesLength
	//	}
	//
	//	go func(start int, end int) {
	//		defer wg.Done()
	//
	//		for _, entityA := range entities[start:end] {
	//			collider := s.GenericCollider.GetUnsafe(entityA)
	//
	//			switch collider.Shape {
	//			case stdcomponents.BoxColliderShape:
	//				s.boxToXCollision(entityA, collisionChan)
	//			case stdcomponents.CircleColliderShape:
	//				s.circleToXCollision(entityA)
	//			default:
	//				panic("Unknown collider shape")
	//			}
	//		}
	//	}(startIndex, endIndex)
	//}
	//
	//// Wait for workers and close collision channel
	//wg.Wait()
	//close(collisionChan)
	//<-doneChan // Wait for result collector
	//
	////s.GenericCollider.EachEntity(func(entity ecs.Entity) bool {
	////	collider := s.GenericCollider.GetUnsafe(entity)
	////
	////	switch collider.Shape {
	////	case stdcomponents.BoxColliderShape:
	////		s.boxToXCollision(entity)
	////	case stdcomponents.CircleColliderShape:
	////		s.circleToXCollision(entity)
	////	default:
	////		panic("Unknown collider shape")
	////	}
	////	return true
	////})
	//
	//s.processExitStates()
}
func (s *CollisionDetectionGridSystem) Destroy() {}

func (s *CollisionDetectionGridSystem) registerCollision(entityA, entityB ecs.Entity, posX, posY float32) {
	pair := CollisionPair{entityA, entityB}.Normalize()

	s.currentCollisions[pair] = struct{}{}

	// Create proxy entity for new collisions
	if _, exists := s.activeCollisions[pair]; !exists {
		proxy := s.EntityManager.Create()
		s.Collisions.Create(proxy, stdcomponents.Collision{E1: pair.E1, E2: pair.E2, State: stdcomponents.CollisionStateEnter})
		s.Positions.Create(proxy, stdcomponents.Position{XY: vectors.Vec2{X: posX, Y: posY}})
		s.activeCollisions[pair] = proxy
	} else {
		s.Collisions.GetUnsafe(s.activeCollisions[pair]).State = stdcomponents.CollisionStateStay
		s.Positions.GetUnsafe(s.activeCollisions[pair]).XY.X = posX
		s.Positions.GetUnsafe(s.activeCollisions[pair]).XY.Y = posY
	}
}

func (s *CollisionDetectionGridSystem) processExitStates() {
	for pair, proxy := range s.activeCollisions {
		if _, exists := s.currentCollisions[pair]; !exists {
			collision := s.Collisions.GetUnsafe(proxy)
			if collision.State == stdcomponents.CollisionStateExit {
				delete(s.activeCollisions, pair)
				s.EntityManager.Delete(proxy)
			} else {
				collision.State = stdcomponents.CollisionStateExit
			}
		}
	}
}

func (s *CollisionDetectionGridSystem) boxToXCollision(entityA ecs.Entity, collisionChan chan<- CollisionEvent) {
	//position1 := s.Positions.GetUnsafe(entityA)
	//spatialIndex1 := s.entityToCell[entityA]
	//genericCollider1 := s.GenericCollider.GetUnsafe(entityA)
	//
	//var nearByIndexes = [9]stdcomponents.SpatialHash{
	//	{spatialIndex1.X, spatialIndex1.Y},
	//	{spatialIndex1.X - 1, spatialIndex1.Y},
	//	{spatialIndex1.X + 1, spatialIndex1.Y},
	//	{spatialIndex1.X, spatialIndex1.Y - 1},
	//	{spatialIndex1.X, spatialIndex1.Y + 1},
	//	{spatialIndex1.X - 1, spatialIndex1.Y - 1},
	//	{spatialIndex1.X + 1, spatialIndex1.Y + 1},
	//	{spatialIndex1.X - 1, spatialIndex1.Y + 1},
	//	{spatialIndex1.X + 1, spatialIndex1.Y - 1},
	//}
	//
	//for _, spatialIndex := range nearByIndexes {
	//	bucket := s.spatialBuckets[spatialIndex]
	//	for _, entityB := range bucket {
	//		if entityA >= entityB {
	//			continue // Skip duplicate checks
	//		}
	//
	//		// Broad Phase
	//		genericCollider2 := s.GenericCollider.GetUnsafe(entityB)
	//		if genericCollider1.Mask&(1<<genericCollider2.Layer) == 0 &&
	//			genericCollider2.Mask&(1<<genericCollider1.Layer) == 0 {
	//			continue
	//		}
	//
	//		// Narrow Phase
	//		position2 := s.Positions.GetUnsafe(entityB)
	//
	//		switch genericCollider2.Shape {
	//		case stdcomponents.BoxColliderShape:
	//			// Inside boxToXCollision
	//			aabbA := s.aabbs[entityA]
	//			aabbB := s.aabbs[entityB]
	//
	//			// Check X-axis first
	//			if aabbA.Right < aabbB.Left || aabbA.Left > aabbB.Right {
	//				continue
	//			}
	//			// Then Y-axis
	//			if aabbA.Bottom < aabbB.Top || aabbA.Top > aabbB.Bottom {
	//				continue
	//			}
	//
	//			posX := (position1.XY.X + position2.XY.X) / 2
	//			posY := (position1.XY.Y + position2.XY.Y) / 2
	//			collisionChan <- CollisionEvent{entityA, entityB, vectors.Vec2{posX, posY}, vectors.Vec2{posX, posY}, 0}
	//		case stdcomponents.CircleColliderShape:
	//			panic("Circle-Box collision not implemented")
	//		default:
	//			panic("Unknown collision shape")
	//		}
	//	}
	//}
}

func (s *CollisionDetectionGridSystem) circleToXCollision(entityA ecs.Entity) {
	panic("Circle-X collision not implemented")
}
