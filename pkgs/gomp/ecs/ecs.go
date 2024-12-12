/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "sync"

type ECSID uint

const (
	MAX_COMPONENTS               = 128
	PREALLOC_BUCKETS      uint32 = 32
	PREALLOC_BUCKETS_SIZE uint32 = 1_000_000
)

type ECS struct {
	ID                  ECSID
	Title               string
	Entities            *SparseSet[Entity, EntityID]
	EntityComponentMask []BitArray
	systems             [][]System
	components          []AnyComponentPtr

	tick            int
	nextEntityID    EntityID
	nextComponentID ComponentID
	wg              *sync.WaitGroup
	mx              *sync.Mutex
}

var nextId ECSID = 0

func generateECSID() ECSID {
	id := nextId
	nextId++
	return id
}

func New(title string) ECS {
	set := NewSparseSet[Entity, EntityID]()

	ecs := ECS{
		ID:       generateECSID(),
		Title:    title,
		Entities: &set,
		// EntityComponentMask: make([]BitArray, ALLOC_BUCKETS_SIZE),

		nextEntityID:    0,
		nextComponentID: 0,
		tick:            0,
		wg:              new(sync.WaitGroup),
		mx:              new(sync.Mutex),
	}

	// for i := 0; i < ALLOC_CHUNK; i++ {
	// 	ecs.EntityComponentMask[i] = NewBitArray(MAX_COMPONENTS)
	// }

	return ecs
}

func (e *ECS) RegisterComponents(component_ptr ...AnyComponentPtr) {
	for i := 0; i < len(component_ptr); i++ {
		e.components = append(e.components, component_ptr[i])
		component_ptr[i].register(e)
	}
}

func (e *ECS) RegisterSystems() *SystemBuilder {
	return &SystemBuilder{
		ecs: e,
	}
}

func (e *ECS) RunSystems() {
	for i := range e.systems {
		// If systems are sequantial, we dont spawn goroutines
		if len(e.systems[i]) == 1 {
			e.systems[i][0].Run(e)
			continue
		}

		e.wg.Add(len(e.systems[i]))
		for j := range e.systems[i] {
			// TODO prespawn goroutines for systems with MAX_N channels, where MAX_N is max number of parallel systems
			go runSystemAsync(e.systems[i][j], e)
		}
		e.wg.Wait()
	}

	// e.Entities.Clean()
	// for i := range e.components {
	// 	e.components[i].Clean(e)
	// }

	e.tick++
}

func (e *ECS) CreateEntity(title string) *Entity {
	e.mx.Lock()
	defer e.mx.Unlock()

	var entity = Entity{
		ID: e.generateEntityID(),
		// Title: title,
		ecsID:     e.ID,
		isDeleted: false,
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

func (e *ECS) SoftDestroyEntity(entity *Entity) {
	e.Entities.SoftDelete(entity.ID)

	for i := range e.components {
		e.components[i].SoftRemove(entity)
	}
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
