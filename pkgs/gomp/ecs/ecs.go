/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type ECSID uint

const (
	MAX_COMPONENTS = 128
	ALLOC_CHUNK    = 1_000_000
)

type ECS struct {
	ID                  ECSID
	Title               string
	Entities            SparseSet[Entity, EntityID]
	EntityComponentMask []BitArray

	nextEntityID    EntityID
	nextComponentID ComponentID
}

type AnyComponentPtr interface {
	register(ecs *ECS)
}

var nextId ECSID = 0

func generateECSID() ECSID {
	id := nextId
	nextId++
	return id
}

func New(title string) ECS {
	ecs := ECS{
		ID:                  generateECSID(),
		Title:               title,
		Entities:            NewSparseSet[Entity, EntityID](ALLOC_CHUNK),
		EntityComponentMask: make([]BitArray, ALLOC_CHUNK),

		nextEntityID:    0,
		nextComponentID: 0,
	}

	for i := 0; i < ALLOC_CHUNK; i++ {
		ecs.EntityComponentMask[i] = NewBitArray(MAX_COMPONENTS)
	}

	return ecs
}

func (e *ECS) generateComponentID() ComponentID {
	id := e.nextComponentID
	e.nextComponentID++
	return id
}

func (e *ECS) generateEntityID() EntityID {
	id := e.nextEntityID
	e.nextEntityID++
	return id
}

func (e *ECS) RegisterComponents(component_ptr ...AnyComponentPtr) {
	for i := 0; i < len(component_ptr); i++ {
		component_ptr[i].register(e)
	}
}

func (e *ECS) CreateEntity(title string) *Entity {
	var entity = Entity{
		ID:    e.generateEntityID(),
		Title: title,
		ecs:   e,
	}

	// if len(e.EntityComponentMask) <= int(entity.ID) {
	// 	e.EntityComponentMask = append(e.EntityComponentMask, make([]BitArray, ALLOC_CHUNK)...)
	// 	for i := int(entity.ID); i < ALLOC_CHUNK; i++ {
	// 		e.EntityComponentMask[i] = NewBitArray(MAX_COMPONENTS)
	// 	}
	// }
	// entity.ComponentsMask = e.EntityComponentMask[entity.ID]

	return e.Entities.Set(entity.ID, entity)
}
