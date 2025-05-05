/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- SeniorOverflow Donated 500 RUB
<- SeniorOverflow Donated 1 000 RUB
<- Монтажёр сука Donated 100 RUB
<- Монтажёр сука Donated 100 RUB

Thank you for your support!
*/

package stdsystems

import (
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"image/color"
	"math/rand"
	"time"

	"github.com/negrel/assert"
)

const (
	collidersPerCell = 8
)

func NewCollisionSetupSystem() CollisionSetupSystem {
	return CollisionSetupSystem{}
}

type CollisionSetupSystem struct {
	EntityManager                       *ecs.EntityManager
	Positions                           *stdcomponents.PositionComponentManager
	Rotations                           *stdcomponents.RotationComponentManager
	Scales                              *stdcomponents.ScaleComponentManager
	Tints                               *stdcomponents.TintComponentManager
	GenericCollider                     *stdcomponents.GenericColliderComponentManager
	BoxColliders                        *stdcomponents.BoxColliderComponentManager
	CircleColliders                     *stdcomponents.CircleColliderComponentManager
	PolygonColliders                    *stdcomponents.PolygonColliderComponentManager
	Collisions                          *stdcomponents.CollisionComponentManager
	SpatialHash                         *stdcomponents.SpatialHashComponentManager
	AABB                                *stdcomponents.AABBComponentManager
	ColliderSleepStateComponentManager  *stdcomponents.ColliderSleepStateComponentManager
	BvhTreeComponentManager             *stdcomponents.BvhTreeComponentManager
	CollisionGridComponentManager       *stdcomponents.CollisionGridComponentManager
	CollisionChunkComponentManager      *stdcomponents.CollisionChunkComponentManager
	CollisionCellComponentManager       *stdcomponents.CollisionCellComponentManager
	CollisionGridMemberComponentManager *stdcomponents.CollisionGridMemberComponentManager

	gridLookup map[stdcomponents.CollisionLayer]ecs.Entity
	Engine     *core.Engine

	clearCellAccumulator []ecs.PagedArray[ecs.Entity]

	memberListPool stdcomponents.MemberListPool
}

func (s *CollisionSetupSystem) Init() {
	s.memberListPool = stdcomponents.NewMemberListPool(s.Engine.Pool())
	s.clearCellAccumulator = make([]ecs.PagedArray[ecs.Entity], s.Engine.Pool().NumWorkers())
	for i := range s.Engine.Pool().NumWorkers() {
		s.clearCellAccumulator[i] = ecs.NewPagedArray[ecs.Entity]()
	}
	s.gridLookup = make(map[stdcomponents.CollisionLayer]ecs.Entity, 64)
}

func (s *CollisionSetupSystem) Run(dt time.Duration) {
	s.setup()
}

func (s *CollisionSetupSystem) Destroy() {
	s.gridLookup = nil
}
func (s *CollisionSetupSystem) setup() {
	// Prepare grids
	s.CollisionGridComponentManager.EachEntity()(func(entity ecs.Entity) bool {
		grid := s.CollisionGridComponentManager.GetUnsafe(entity)
		assert.NotNil(grid)

		s.gridLookup[grid.Layer] = entity

		if grid.CreateCellsAccumulator == nil {
			grid.CreateCellsAccumulator = make([]ecs.GenMap[stdcomponents.SpatialCellIndex, struct{}], s.Engine.Pool().NumWorkers())
			for i := range grid.CreateCellsAccumulator {
				grid.CreateCellsAccumulator[i] = ecs.NewGenMap[stdcomponents.SpatialCellIndex, struct{}](1024)
			}
		}

		return true
	})

	// Create grid member components
	var gridMemberAccumulator = make([][]ecs.Entity, s.Engine.Pool().NumWorkers())
	s.GenericCollider.ProcessEntities(func(entity ecs.Entity, id worker.WorkerId) {
		if s.CollisionGridMemberComponentManager.Has(entity) {
			return
		}
		gridMemberAccumulator[id] = append(gridMemberAccumulator[id], entity)
	})
	for i := range gridMemberAccumulator {
		acc := gridMemberAccumulator[i]
		for j := range acc {
			entity := acc[j]
			s.CollisionGridMemberComponentManager.Create(entity, stdcomponents.CollisionGridMember{})
		}
	}

	// Update grid member component
	s.GenericCollider.ProcessEntities(func(entity ecs.Entity, id worker.WorkerId) {
		collider := s.GenericCollider.GetUnsafe(entity)
		assert.NotNil(collider)

		gridMember := s.CollisionGridMemberComponentManager.GetUnsafe(entity)
		assert.NotNil(gridMember)

		gridEntity, exists := s.gridLookup[collider.Layer]
		assert.True(exists)

		gridMember.Grid = gridEntity
	})

	// Create spatial hash components
	var spatialHashAccumulator = make([][]ecs.Entity, s.Engine.Pool().NumWorkers())
	s.GenericCollider.ProcessEntities(func(entity ecs.Entity, id worker.WorkerId) {
		if s.SpatialHash.Has(entity) {
			return
		}
		spatialHashAccumulator[id] = append(spatialHashAccumulator[id], entity)
	})
	for i := range spatialHashAccumulator {
		acc := spatialHashAccumulator[i]
		for j := range acc {
			entity := acc[j]
			s.SpatialHash.Create(entity, stdcomponents.SpatialHash{})
		}
		spatialHashAccumulator[i] = spatialHashAccumulator[i][:0]
	}

	// Calculate spatial hash for each collider
	s.CollisionGridMemberComponentManager.ProcessEntities(func(entity ecs.Entity, id worker.WorkerId) {
		gridMember := s.CollisionGridMemberComponentManager.GetUnsafe(entity)
		assert.NotNil(gridMember)

		grid := s.CollisionGridComponentManager.GetUnsafe(gridMember.Grid)
		assert.NotNil(grid)

		spatialHash := s.SpatialHash.GetUnsafe(entity)
		assert.NotNil(spatialHash)

		bb := s.AABB.GetUnsafe(entity)
		assert.NotNil(bb)

		newSpatialHash := grid.CalculateSpatialHash(*bb)
		spatialHash.Min = newSpatialHash.Min
		spatialHash.Max = newSpatialHash.Max
	})

	// Accumulate spatial hash that need to be created
	s.SpatialHash.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		spatialHash := s.SpatialHash.GetUnsafe(entity)
		assert.NotNil(spatialHash)

		gridMember := s.CollisionGridMemberComponentManager.GetUnsafe(entity)
		assert.NotNil(gridMember)

		grid := s.CollisionGridComponentManager.GetUnsafe(gridMember.Grid)
		assert.NotNil(grid)

		for i := spatialHash.Min.X; i <= spatialHash.Max.X; i++ {
			for j := spatialHash.Min.Y; j <= spatialHash.Max.Y; j++ {
				index := stdcomponents.SpatialCellIndex{X: i, Y: j}
				if grid.CellMap.Has(index) {
					continue
				}
				grid.CreateCellsAccumulator[workerId].Set(index, struct{}{})
			}
		}
	})

	// Create requested cells
	s.CollisionGridComponentManager.EachEntityParallel(func(gridEntity ecs.Entity, workerId worker.WorkerId) {
		grid := s.CollisionGridComponentManager.GetUnsafe(gridEntity)
		assert.NotNil(grid)
		for i := range grid.CreateCellsAccumulator {
			set := &grid.CreateCellsAccumulator[i]
			set.Clear()
			for cellIndex := range set.Each() {
				if grid.CellMap.Has(cellIndex) {
					continue
				}
				cellEntity := s.EntityManager.Create()
				s.CollisionCellComponentManager.Create(cellEntity, stdcomponents.CollisionCell{
					Index: cellIndex,
					Layer: grid.Layer,
					Grid:  gridEntity,
					Size:  grid.CellSize,
				})
				s.Positions.Create(cellEntity, stdcomponents.Position{
					XY: cellIndex.ToVec2().Scale(grid.CellSize),
				})
				const colorbase uint8 = 120
				s.Tints.Create(cellEntity, color.RGBA{
					R: colorbase + uint8(rand.Intn(int(255-colorbase))),
					G: colorbase + uint8(rand.Intn(int(255-colorbase))),
					B: colorbase + uint8(rand.Intn(int(255-colorbase))),
					A: 70,
				})
				grid.CellMap.Set(cellIndex, cellEntity)
			}
			grid.CreateCellsAccumulator[i].Reset()
		}
	})
	s.CollisionCellComponentManager.ProcessComponents(func(cell *stdcomponents.CollisionCell, id worker.WorkerId) {
		if cell.Members != nil {
			return
		}
		cell.Members = s.memberListPool.Get()
	})

	// Remove entities from cells
	s.CollisionCellComponentManager.ProcessComponents(func(cell *stdcomponents.CollisionCell, id worker.WorkerId) {
		memberList := cell.Members
		assert.NotNil(memberList)
		for i := len(memberList.Members) - 1; i >= 0; i-- {
			member := memberList.Members[i]
			spatialHash := s.SpatialHash.GetUnsafe(member)
			assert.NotNil(spatialHash)

			// Check if entity is in cell
			if cell.Index.X >= spatialHash.Min.X && cell.Index.X <= spatialHash.Max.X {
				if cell.Index.Y >= spatialHash.Min.Y && cell.Index.Y <= spatialHash.Max.Y {
					continue
				}
			}

			// Entity should be removed from cell
			memberList.Members = append(memberList.Members[:i], memberList.Members[i+1:]...)
			memberList.Lookup.Delete(member)
		}
		memberList.Lookup.Clear()
	})

	// Distribute entities in grid cells
	s.SpatialHash.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		spatialHash := s.SpatialHash.GetUnsafe(entity)
		assert.NotNil(spatialHash)

		gridMember := s.CollisionGridMemberComponentManager.GetUnsafe(entity)
		assert.NotNil(gridMember)

		grid := s.CollisionGridComponentManager.GetUnsafe(gridMember.Grid)
		assert.NotNil(grid)

		for i := spatialHash.Min.X; i <= spatialHash.Max.X; i++ {
			for j := spatialHash.Min.Y; j <= spatialHash.Max.Y; j++ {
				cellEntity, exists := grid.CellMap.Get(stdcomponents.SpatialCellIndex{X: i, Y: j})
				assert.True(exists)

				cell := s.CollisionCellComponentManager.GetUnsafe(cellEntity)
				assert.NotNil(cell)
				members := cell.Members
				assert.NotNil(members)

				if members.Has(entity) {
					continue
				}

				members.InputAcc[workerId] = append(members.InputAcc[workerId], entity)
			}
		}
	})

	//Build grid cells
	s.CollisionCellComponentManager.ProcessComponents(func(cell *stdcomponents.CollisionCell, workerId worker.WorkerId) {
		members := cell.Members
		for i := range members.InputAcc {
			acc := members.InputAcc[i]
			for j := range acc {
				members.Members = append(members.Members, acc[j])
				members.Lookup.Set(acc[j], len(members.Members)-1)
			}
			members.InputAcc[i] = acc[:0]
		}
	})

	//Remove empty cells
	s.CollisionCellComponentManager.ProcessEntities(func(cellEntity ecs.Entity, workerId worker.WorkerId) {
		cell := s.CollisionCellComponentManager.GetUnsafe(cellEntity)
		assert.NotNil(cell)

		if len(cell.Members.Members) != 0 {
			return
		}
		cell.Members.Reset()
		s.memberListPool.Put(cell.Members)
		cell.Members = nil
		s.clearCellAccumulator[workerId].Append(cellEntity)
	})
	for i := range s.clearCellAccumulator {
		v := &s.clearCellAccumulator[i]
		v.EachDataValue()(func(cellEntity ecs.Entity) bool {
			cell := s.CollisionCellComponentManager.GetUnsafe(cellEntity)
			assert.NotNil(cell)

			grid := s.CollisionGridComponentManager.GetUnsafe(cell.Grid)
			assert.NotNil(grid)

			grid.CellMap.Delete(cell.Index)
			return true
		})
		v.EachDataValue()(func(cellEntity ecs.Entity) bool {
			s.EntityManager.Delete(cellEntity)
			return true
		})
		v.Reset()
	}
	s.CollisionGridComponentManager.EachComponentParallel(func(grid *stdcomponents.CollisionGrid, workerId worker.WorkerId) {
		grid.CellMap.Clear()
	})
}
func (s *CollisionSetupSystem) detection() {
}
