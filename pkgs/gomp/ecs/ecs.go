/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
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

	size          uint32
	shouldDestroy bool

	systems             []AnySystemPtr[World]
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

func (e *World) RegisterComponents(component_ptr ...AnyComponentTypePtr[World]) {
	for i := 0; i < len(component_ptr); i++ {
		e.components = append(e.components, component_ptr[i].register(e, ComponentID(i)))
	}
}

func (e *World) RegisterSystems(systems ...AnySystemManagerPtr) {
	for i := range systems {
		e.systems = append(e.systems, systems[i].register(e))
		e.systems[i].Init(e)
	}
}

func (e *World) RunSystems() error {
	e.wg.Add(len(e.systems))

	for i := range e.systems {
		func(system AnySystemPtr[World], e *World) {
			defer e.wg.Done()
			system.Run(e)
		}(e.systems[i], e)
	}

	e.wg.Wait()

	e.tick++
	e.Clean()

	return nil
}

func (e *World) CreateEntity(title string) EntityID {
	var newId = e.generateEntityID()

	e.entityComponentMask.Set(newId, ComponentBitArray256{})
	e.size++

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
	e.size--
}

func (e *World) Clean() {
	for i := range e.components {
		e.components[i].Clean()
	}
}

func (e *World) Size() uint32 {
	return e.size
}

func (e *World) LastEntityID() EntityID {
	return e.lastEntityID
}

func (e *World) ShouldDestroy() bool {
	return e.shouldDestroy
}

func (e *World) SetShouldDestroy(value bool) {
	e.shouldDestroy = value
}

func (e *World) Destroy() {
	e.wg.Add(len(e.systems))

	for i := range e.systems {
		func(system AnySystemPtr[World], e *World) {
			defer e.wg.Done()
			system.Destroy(e)
		}(e.systems[i], e)
	}

	e.wg.Wait()

	e.Clean()
}

func (w *World) Run(tickrate uint) {
	var ticker *time.Ticker
	if tickrate > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(tickrate))
	} else {
		ticker = time.NewTicker(0)
	}

	for !w.ShouldDestroy() {
		select {
		case <-ticker.C:
			w.RunSystems()

			if len(ticker.C) > 0 {
				<-ticker.C
				log.Println("Skipping tick")
			}
		default:
		}
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
