/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Фрея Donated 2 000 RUB

Thank you for your support!
*/

package stdsystems

import (
	"github.com/negrel/assert"
	gjk "gomp/pkg/collision"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"math/bits"
	"time"
)

func NewCollisionDetectionSystem() CollisionDetectionSystem {
	return CollisionDetectionSystem{}
}

type CollisionDetectionSystem struct {
	EntityManager                      *ecs.EntityManager
	Positions                          *stdcomponents.PositionComponentManager
	Rotations                          *stdcomponents.RotationComponentManager
	Scales                             *stdcomponents.ScaleComponentManager
	GenericCollider                    *stdcomponents.GenericColliderComponentManager
	BoxColliders                       *stdcomponents.BoxColliderComponentManager
	CircleColliders                    *stdcomponents.CircleColliderComponentManager
	PolygonColliders                   *stdcomponents.PolygonColliderComponentManager
	Collisions                         *stdcomponents.CollisionComponentManager
	SpatialIndex                       *stdcomponents.SpatialHashComponentManager
	AABB                               *stdcomponents.AABBComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	BvhTreeComponentManager            *stdcomponents.BvhTreeComponentManager
	CollisionGridComponentManager      *stdcomponents.CollisionGridComponentManager
	CollisionChunkComponentManager     *stdcomponents.CollisionChunkComponentManager
	CollisionCellComponentManager      *stdcomponents.CollisionCellComponentManager
	Engine                             *core.Engine

	collisionEventAcc []ecs.PagedArray[CollisionEvent]
	activeCollisions  map[CollisionPair]ecs.Entity // Maps collision pairs to proxy entities
	currentCollisions map[CollisionPair]struct{}
	gridLookup        map[stdcomponents.CollisionLayer]*stdcomponents.CollisionGrid
	potentialEntities []ecs.PagedArray[ecs.Entity]
}

func (s *CollisionDetectionSystem) Init() {
	s.gridLookup = make(map[stdcomponents.CollisionLayer]*stdcomponents.CollisionGrid)
	s.collisionEventAcc = make([]ecs.PagedArray[CollisionEvent], s.Engine.Pool().NumWorkers())
	s.potentialEntities = make([]ecs.PagedArray[ecs.Entity], s.Engine.Pool().NumWorkers())
	for i := 0; i < s.Engine.Pool().NumWorkers(); i++ {
		s.collisionEventAcc[i] = ecs.NewPagedArray[CollisionEvent]()
		s.potentialEntities[i] = ecs.NewPagedArray[ecs.Entity]()
	}
	s.activeCollisions = make(map[CollisionPair]ecs.Entity)
}

func (s *CollisionDetectionSystem) Run(dt time.Duration) {
	s.CollisionGridComponentManager.EachComponent()(func(grid *stdcomponents.CollisionGrid) bool {
		s.gridLookup[grid.Layer] = grid
		return true
	})

	s.currentCollisions = make(map[CollisionPair]struct{})

	if s.GenericCollider.Len() == 0 {
		return
	}

	s.GenericCollider.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		potentialEntities := &s.potentialEntities[workerId]
		s.broadPhase(entity, potentialEntities)
		if potentialEntities.Len() == 0 {
			return
		}
		s.narrowPhase(entity, potentialEntities, workerId)
		potentialEntities.Reset()
	})

	s.registerCollisionEvents()
	s.processExitStates()

	for i := 0; i < len(s.collisionEventAcc); i++ {
		s.collisionEventAcc[i].Reset()
	}
}

func (s *CollisionDetectionSystem) Destroy() {
	s.collisionEventAcc = nil
	s.activeCollisions = nil
	s.currentCollisions = nil
	s.gridLookup = nil
}

func (s *CollisionDetectionSystem) broadPhase(entityA ecs.Entity, potentialEntities *ecs.PagedArray[ecs.Entity]) {
	colliderA := s.GenericCollider.GetUnsafe(entityA)

	// Early exit for sleeping colliders (moved aabb access after sleep check)
	if colliderA.AllowSleep && s.ColliderSleepStateComponentManager.Has(entityA) {
		return
	}

	aabbPtr := s.AABB.GetUnsafe(entityA)
	assert.NotNil(aabbPtr)
	bb := *aabbPtr

	// Direct layer bitmask iteration
	mask := colliderA.Mask

	// Iterate only set bits in mask
	for mask != 0 {
		// Get the least significant raised bit position
		layer := stdcomponents.CollisionLayer(bits.TrailingZeros32(uint32(mask)))
		mask ^= 1 << layer // Clear the processed bit

		// Direct grid access
		grid := s.gridLookup[layer]
		assert.NotNil(grid)

		// Query grid
		minSpatialCellIndex := grid.GetCellIndex(bb.Min)
		maxSpatialCellIndex := grid.GetCellIndex(bb.Max)
		// make a list of all spatial indexes that intersect the aabb
		// get cells that intersect the aabb by spatial indexes
		for i := minSpatialCellIndex.X; i <= maxSpatialCellIndex.X; i++ {
			for j := minSpatialCellIndex.Y; j <= maxSpatialCellIndex.Y; j++ {
				cellEntity, exists := grid.CellMap.Get(stdcomponents.SpatialCellIndex{X: i, Y: j})
				if !exists {
					continue
				}

				cell := s.CollisionCellComponentManager.GetUnsafe(cellEntity)
				assert.NotNil(cell)

				members := cell.Members.Members

				for _, entityB := range members {
					if entityB != entityA {
						potentialEntities.Append(entityB)
					}
				}
			}
		}
	}
}

func (s *CollisionDetectionSystem) narrowPhase(entityA ecs.Entity, potentialEntities *ecs.PagedArray[ecs.Entity], workerId worker.WorkerId) {
	posA := s.Positions.GetUnsafe(entityA)
	assert.NotNil(posA)

	colliderA := s.GenericCollider.GetUnsafe(entityA)
	assert.NotNil(colliderA)

	scaleA := s.Scales.GetUnsafe(entityA)
	assert.NotNil(scaleA)

	rotA := s.Rotations.GetUnsafe(entityA)
	assert.NotNil(rotA)

	aabbA := s.AABB.GetUnsafe(entityA)
	assert.NotNil(aabbA)

	colA := s.getGjkCollider(colliderA, entityA)

	circleA := s.CircleColliders.GetUnsafe(entityA)
	// Cache circleA properties if exists
	var radiusA float32
	if circleA != nil {
		radiusA = circleA.Radius * scaleA.XY.X
	}

	transformA := stdcomponents.Transform2d{
		Position: posA.XY,
		Rotation: rotA.Angle,
		Scale:    scaleA.XY,
	}

	for entityB := range potentialEntities.EachDataValue() {
		aabbB := s.AABB.GetUnsafe(entityB)
		assert.NotNil(aabbB)

		if !(aabbA.Max.X >= aabbB.Min.X && aabbA.Min.X <= aabbB.Max.X &&
			aabbA.Max.Y >= aabbB.Min.Y && aabbA.Min.Y <= aabbB.Max.Y) {
			continue
		}

		positionB := s.Positions.GetUnsafe(entityB)
		assert.NotNil(positionB)
		colliderB := s.GenericCollider.GetUnsafe(entityB)
		assert.NotNil(colliderB)
		scaleB := s.Scales.GetUnsafe(entityB)
		assert.NotNil(scaleB)
		rotationB := s.Rotations.GetUnsafe(entityB)
		assert.NotNil(rotationB)
		circleB := s.CircleColliders.GetUnsafe(entityB)

		// 1. FAST PATH: Circle-circle collision
		if circleA != nil && circleB != nil {
			posB := positionB.XY
			radiusB := circleB.Radius * scaleB.XY.X

			// Vector math with early exit
			dx := posB.X - transformA.Position.X
			dy := posB.Y - transformA.Position.Y
			sqDist := dx*dx + dy*dy
			sumRadii := radiusA + radiusB

			if sqDist < sumRadii*sumRadii {
				dist := float32(math.Sqrt(float64(sqDist)))
				assert.NotZero(dist)
				event := CollisionEvent{
					entityA:  entityA,
					entityB:  entityB,
					position: posA.XY,
					normal: vectors.Vec2{
						X: dx / dist,
						Y: dy / dist,
					},
					depth: sumRadii - dist,
				}
				s.collisionEventAcc[workerId].Append(event)
			}
			continue
		}

		// 2. GJK/EPA PATH
		test := gjk.New()
		transformB := stdcomponents.Transform2d{
			Position: positionB.XY,
			Rotation: rotationB.Angle,
			Scale:    scaleB.XY,
		}
		colB := s.getGjkCollider(colliderB, entityB)
		// Detect collision using GJK
		if test.CheckCollision(colA, colB, transformA, transformB) {
			// If collision detected, get penetration details using EPA
			normal, depth := test.EPA(colA, colB, transformA, transformB)
			position := posA.XY.Add(positionB.XY.Sub(posA.XY))
			s.collisionEventAcc[workerId].Append(CollisionEvent{
				entityA:  entityA,
				entityB:  entityB,
				position: position,
				normal:   normal,
				depth:    depth,
			})
		}
	}
}
func (s *CollisionDetectionSystem) registerCollisionEvents() {
	for i := range s.collisionEventAcc {
		events := &s.collisionEventAcc[i]
		events.EachData()(func(event *CollisionEvent) bool {
			pair := CollisionPair{event.entityA, event.entityB}
			s.currentCollisions[pair] = struct{}{}
			displacement := event.normal.Scale(event.depth)
			pos := event.position.Add(displacement)

			if _, exists := s.activeCollisions[pair]; !exists {
				proxy := s.EntityManager.Create()
				s.Collisions.Create(proxy, stdcomponents.Collision{
					E1:     pair.E1,
					E2:     pair.E2,
					State:  stdcomponents.CollisionStateEnter,
					Normal: event.normal,
					Depth:  event.depth,
				})

				s.Positions.Create(proxy, stdcomponents.Position{
					XY: vectors.Vec2{
						X: pos.X, Y: pos.Y,
					}})
				s.activeCollisions[pair] = proxy
			} else {
				proxy := s.activeCollisions[pair]
				collision := s.Collisions.GetUnsafe(proxy)
				position := s.Positions.GetUnsafe(proxy)
				collision.State = stdcomponents.CollisionStateStay
				collision.Depth = event.depth
				collision.Normal = event.normal
				position.XY.X = pos.X
				position.XY.Y = pos.Y
			}
			return true
		})
		events.Reset()
	}
}

func (s *CollisionDetectionSystem) processExitStates() {
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
func (s *CollisionDetectionSystem) getGjkCollider(collider *stdcomponents.GenericCollider, entity ecs.Entity) gjk.AnyCollider {
	switch collider.Shape {
	case stdcomponents.BoxColliderShape:
		return s.BoxColliders.GetUnsafe(entity)
	case stdcomponents.CircleColliderShape:
		return s.CircleColliders.GetUnsafe(entity)
	case stdcomponents.PolygonColliderShape:
		return s.PolygonColliders.GetUnsafe(entity)
	default:
		panic("unsupported collider shape")
	}
	return nil
}
