/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"fmt"
	"sync"
)

type ECSID uint

const (
	PREALLOC_DELETED_ENTITIES uint32 = 1 << 10
)

type World struct {
	ID    ECSID
	Title string

	tick         int
	nextEntityID EntityID

	systems             [][]System
	components          []AnyComponentInstancesPtr
	deletedEntityIDs    []EntityID
	entityComponentMask *SparseSet[ComponentBitArray256, EntityID]
	wg                  *sync.WaitGroup
	mx                  *sync.Mutex
}

var nextId ECSID = 0

func generateECSID() ECSID {
	id := nextId
	nextId++
	return id
}

func New(title string) World {
	id := generateECSID()
	maskSet := NewSparseSet[ComponentBitArray256, EntityID]()

	ecs := World{
		ID:                  id,
		Title:               title,
		wg:                  new(sync.WaitGroup),
		mx:                  new(sync.Mutex),
		deletedEntityIDs:    make([]EntityID, 0, PREALLOC_DELETED_ENTITIES),
		entityComponentMask: &maskSet,
	}

	return ecs
}

func (e *World) RegisterComponentTypes(component_ptr ...AnyComponentTypePtr) {
	for i := 0; i < len(component_ptr); i++ {
		e.components = append(e.components, component_ptr[i].register(e, ComponentID(i)))
	}
}

func (e *World) RegisterSystems() *SystemBuilder {
	return &SystemBuilder{
		ecs: e,
	}
}

func (e *World) RunSystems() {
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

func (e *World) CreateEntity(title string) EntityID {
	e.mx.Lock()
	defer e.mx.Unlock()

	var newId EntityID

	if len(e.deletedEntityIDs) == 0 {
		newId = e.generateEntityID()
	} else {
		newId = e.deletedEntityIDs[len(e.deletedEntityIDs)-1]
		e.deletedEntityIDs = e.deletedEntityIDs[:len(e.deletedEntityIDs)-1]
	}

	e.entityComponentMask.Set(newId, ComponentBitArray256{})

	return newId
}

func (e *World) SoftDestroyEntity(entityId EntityID) {
	e.mx.Lock()
	defer e.mx.Unlock()

	mask := e.entityComponentMask.GetPtr(entityId)
	if mask == nil {
		panic(fmt.Sprintf("Entity %d does not exist", entityId))
	}

	for i := range mask.AllSet {
		e.components[i].SoftRemove(entityId)
	}

	e.entityComponentMask.SoftDelete(entityId)
	e.deletedEntityIDs = append(e.deletedEntityIDs, entityId)
}

func (e *World) generateEntityID() EntityID {
	id := e.nextEntityID
	e.nextEntityID++
	return id
}
