/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

type ECSID uint

const (
	PREALLOC_DELETED_ENTITIES uint32 = 1 << 10
)

type World struct {
	ID    ECSID
	Title string

	tick         int
	lastEntityID EntityID

	updateSystems       [][]AnyUpdateSystem[World]
	drawSystems         [][]AnyDrawSystem[World]
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

func CreateWorld(title string) World {
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

func (e *World) RegisterComponentTypes(component_ptr ...AnyComponentTypePtr[World]) {
	for i := 0; i < len(component_ptr); i++ {
		e.components = append(e.components, component_ptr[i].register(e, ComponentID(i)))
	}
}

func (e *World) RegisterUpdateSystems() *UpdateSystemBuilder[World] {
	return &UpdateSystemBuilder[World]{
		world:   e,
		systems: &e.updateSystems,
	}
}

func (e *World) RegisterDrawSystems() *DrawSystemBuilder[World] {
	return &DrawSystemBuilder[World]{
		ecs:     e,
		systems: &e.drawSystems,
	}
}

func (e *World) RunUpdateSystems() error {
	for i := range e.updateSystems {
		// If systems are sequantial, we dont spawn goroutines
		if len(e.updateSystems[i]) == 1 {
			e.updateSystems[i][0].Run(e)
			continue
		}

		e.wg.Add(len(e.updateSystems[i]))
		for j := range e.updateSystems[i] {
			// TODO prespawn goroutines for systems with MAX_N channels, where MAX_N is max number of parallel systems
			go func(system AnyUpdateSystem[World], e *World) {
				defer e.wg.Done()
				system.Run(e)
			}(e.updateSystems[i][j], e)
		}
		e.wg.Wait()
	}

	e.tick++
	e.Clean()

	return nil
}

func (e *World) RunDrawSystems(screen *ebiten.Image) {
	for i := range e.drawSystems {
		// If systems are sequantial, we dont spawn goroutines
		if len(e.drawSystems[i]) == 1 {
			e.drawSystems[i][0].Run(e, screen)
			continue
		}

		e.wg.Add(len(e.drawSystems[i]))
		for j := range e.drawSystems[i] {
			// TODO prespawn goroutines for systems with MAX_N channels, where MAX_N is max number of parallel systems
			go func(system AnyDrawSystem[World], e *World, screen *ebiten.Image) {
				defer e.wg.Done()
				system.Run(e, screen)
			}(e.drawSystems[i][j], e, screen)
		}
		e.wg.Wait()
	}
}

func (e *World) CreateEntity(title string) EntityID {
	var newId = e.generateEntityID()

	e.entityComponentMask.Set(newId, ComponentBitArray256{})

	return newId
}

func (e *World) DestroyEntity(entityId EntityID) {
	mask := e.entityComponentMask.GetPtr(entityId)
	if mask == nil {
		panic(fmt.Sprintf("Entity %d does not exist", entityId))
	}

	for i := range mask.AllSet {
		e.components[i].Remove(entityId)
	}

	e.entityComponentMask.SoftDelete(entityId)
	e.deletedEntityIDs = append(e.deletedEntityIDs, entityId)
}

func (e *World) Clean() {
	for i := range e.components {
		e.components[i].Clean()
	}
}

func (e *World) generateEntityID() (newId EntityID) {
	if len(e.deletedEntityIDs) == 0 {
		newId = EntityID(atomic.AddInt32((*int32)(&e.lastEntityID), 1))
	} else {
		newId = e.deletedEntityIDs[len(e.deletedEntityIDs)-1]
		e.deletedEntityIDs = e.deletedEntityIDs[:len(e.deletedEntityIDs)-1]
	}
	return newId
}
