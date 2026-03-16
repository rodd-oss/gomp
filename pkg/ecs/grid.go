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

package ecs

import (
	"math"
	"sync"
)

// GridKey is {cellX, cellY} key for SpatialGrid
type GridKey = [2]int

type GridCell struct {
	Entities []Entity
	lookup   map[Entity]int
	Key      GridKey
}

/*
SpatialGrid

# To maximize performance:

- Set cellSize to ~2x average collision radius

- Call Compact() during level transitions

- Use GetNearbyCells() for broadphase collision

- Pair with spatial hashing for very large worlds
*/
type SpatialGrid struct {
	cells          []GridCell
	keys           []GridKey
	lookupByKey    map[GridKey]int
	lookupByEntity map[Entity]int

	cellSize float64
	mx       sync.RWMutex

	updated map[GridKey]struct{}
}

func NewSpatialGrid(cellSize float64) SpatialGrid {
	return SpatialGrid{
		cells:          make([]GridCell, 0, PREALLOC_DEFAULT),
		keys:           make([]GridKey, 0, PREALLOC_DEFAULT),
		lookupByKey:    make(map[GridKey]int, PREALLOC_DEFAULT),
		lookupByEntity: make(map[Entity]int, PREALLOC_DEFAULT),
		cellSize:       cellSize,
	}
}

// GetCell returns the grid cell for a given position.
func (g *SpatialGrid) GetCell(positionX, positionY float64) (*GridCell, GridKey, bool) {
	cellX := int(math.Floor(positionX / g.cellSize))
	cellY := int(math.Floor(positionY / g.cellSize))
	key := GridKey{cellX, cellY}
	cellIndex, exists := g.lookupByKey[key]
	if !exists {
		return nil, key, false
	}
	return &g.cells[cellIndex], key, true
}

// AddEntity adds an entity to the grid.
func (g *SpatialGrid) AddEntity(entity Entity, positionX, positionY float64) {
	g.mx.Lock()
	defer g.mx.Unlock()

	cellX := int(math.Floor(positionX / g.cellSize))
	cellY := int(math.Floor(positionY / g.cellSize))

	key := GridKey{cellX, cellY}

	cellIndex, ok := g.lookupByKey[key]
	if !ok {
		cellIndex = len(g.cells)
		g.cells = append(g.cells, GridCell{
			Entities: make([]Entity, 0, PREALLOC_DEFAULT),
			lookup:   make(map[Entity]int, PREALLOC_DEFAULT),
			Key:      key,
		})
	}
	cell := &g.cells[cellIndex]

	cell.Entities = append(cell.Entities, entity)
	cell.lookup[entity] = len(cell.Entities) - 1

	g.lookupByEntity[entity] = cellIndex
	g.lookupByKey[key] = cellIndex
}

// GetEntities returns the Entities in a given cell.
func (g *SpatialGrid) GetEntities(cellX, cellY int) []Entity {
	key := GridKey{cellX, cellY}
	cellIndex, ok := g.lookupByKey[key]
	if !ok {
		return nil
	}
	return g.cells[cellIndex].Entities
}

func (g *SpatialGrid) GetCellByKey(key GridKey) (*GridCell, GridKey, bool) {
	cellIndex, ok := g.lookupByKey[key]
	if !ok {
		return nil, key, false
	}
	return &g.cells[cellIndex], key, true
}

// UpdateEntity updates state of an entity in the grid.
func (g *SpatialGrid) UpdateEntity(entity Entity) {
	g.updated[g.keys[g.lookupByEntity[entity]]] = struct{}{}
}

// RemoveEntity removes an entity from the grid.
func (g *SpatialGrid) RemoveEntity(entity Entity) {
	cellIndex, exists := g.lookupByEntity[entity]
	if !exists {
		return
	}

	cell := &g.cells[cellIndex]

	last := len(cell.Entities) - 1
	entityIndex := cell.lookup[entity]

	if entityIndex != last {
		cell.Entities[entityIndex] = cell.Entities[last]
		cell.lookup[cell.Entities[last]] = entityIndex
	}

	cell.Entities = cell.Entities[:last]
	delete(cell.lookup, entity)
	delete(g.lookupByEntity, entity)

	if len(cell.Entities) == 0 {
		delete(g.lookupByKey, GridKey{cellIndex, cellIndex})
		g.cells[cellIndex] = GridCell{} // Release memory
		g.keys[cellIndex] = GridKey{}
	}
}

// MoveEntity updates the position of an entity in the grid.
func (g *SpatialGrid) MoveEntity(entity Entity, positionX, positionY float64) {
	oldCellIndex := g.lookupByEntity[entity]
	oldCell := &g.cells[oldCellIndex]

	newCellX := int(math.Floor(positionX / g.cellSize))
	newCellY := int(math.Floor(positionY / g.cellSize))
	newKey := GridKey{newCellX, newCellY}

	// If the entity hasn't changed cells, no need to move it
	if g.keys[oldCellIndex] == newKey {
		g.updated[newKey] = struct{}{}
		return
	}

	// Remove from old cell
	last := len(oldCell.Entities) - 1
	entityIndex := oldCell.lookup[entity]
	if entityIndex != last {
		oldCell.Entities[entityIndex] = oldCell.Entities[last]
		oldCell.lookup[oldCell.Entities[last]] = entityIndex
	}
	oldCell.Entities = oldCell.Entities[:last]
	delete(oldCell.lookup, entity)

	// Add to new cell
	newCellIndex, ok := g.lookupByKey[newKey]
	if !ok {
		newCellIndex = len(g.cells)
		g.cells = append(g.cells, GridCell{
			Entities: make([]Entity, 0, PREALLOC_DEFAULT),
			lookup:   make(map[Entity]int, PREALLOC_DEFAULT),
		})
		g.lookupByKey[newKey] = newCellIndex
		g.keys = append(g.keys, newKey)
	}
	newCell := &g.cells[newCellIndex]
	newCell.Entities = append(newCell.Entities, entity)
	newCell.lookup[entity] = len(newCell.Entities) - 1

	// Update entity lookup
	g.lookupByEntity[entity] = newCellIndex
	g.updated[newKey] = struct{}{}
}

func (g *SpatialGrid) GetNearbyCells(key GridKey) []GridKey {
	return []GridKey{
		{key[0] - 1, key[1] - 1}, {key[0], key[1] - 1}, {key[0] + 1, key[1] - 1},
		{key[0] - 1, key[1]}, key, {key[0] + 1, key[1]},
		{key[0] - 1, key[1] + 1}, {key[0], key[1] + 1}, {key[0] + 1, key[1] + 1},
	}
}

func (g *SpatialGrid) GetNearbyEntities(key GridKey) []Entity {
	entities := make([]Entity, 0, 32)
	for _, k := range g.GetNearbyCells(key) {
		if cell, _, exists := g.GetCellByKey(k); exists {
			entities = append(entities, cell.Entities...)
		}
	}
	return entities
}

func (g *SpatialGrid) Compact() {
	g.mx.Lock()
	defer g.mx.Unlock()

	newCells := make([]GridCell, 0, len(g.cells))
	newKeys := make([]GridKey, 0, len(g.keys))
	newLookup := make(map[GridKey]int, len(g.lookupByKey))

	for i, cell := range g.cells {
		if len(cell.Entities) > 0 {
			newLookup[g.keys[i]] = len(newCells)
			newCells = append(newCells, cell)
			newKeys = append(newKeys, g.keys[i])
		}
	}

	// Rebuild entity lookup
	g.lookupByEntity = make(map[Entity]int, len(g.lookupByEntity))
	for i, cell := range newCells {
		for _, e := range cell.Entities {
			g.lookupByEntity[e] = i
		}
	}

	g.cells = newCells
	g.keys = newKeys
	g.lookupByKey = newLookup
}

func (g *SpatialGrid) BoundingBoxesIntersect(x1, y1, x2, y2 float64) bool {
	// Fast AABB check using spatial grid cell size
	dx := math.Abs(x1 - x2)
	dy := math.Abs(y1 - y2)
	return dx < g.cellSize*2 && dy < g.cellSize*2
}

func (g *SpatialGrid) GetActiveCells() []GridCell {
	return g.cells
}
