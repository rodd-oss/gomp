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

	systems             [][]*SystemServiceInstance
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

func (w *World) RegisterComponentServices(component_ptr ...AnyComponentServicePtr) {
	for i := 0; i < len(component_ptr); i++ {
		w.components = append(w.components, component_ptr[i].register(w, ComponentID(i)))
	}
}

func (e *World) RegisterSystems() *SystemBuilder {
	return &SystemBuilder{
		world:   e,
		systems: &e.systems,
	}
}

func (w *World) FixedUpdate() error {
	w.runSystemFunction(SystemFunctionFixedUpdate)
	return nil
}

func (w *World) runSystemFunction(method SystemFunctionMethod) error {
	for i := range w.systems {
		parallel := w.systems[i]

		if len(parallel) == 1 {
			controller := parallel[0].controller
			switch method {
			case systemFunctionInit:
				controller.Init(w)
			case systemFunctionUpdate:
				controller.Update(w)
			case SystemFunctionFixedUpdate:
				controller.FixedUpdate(w)
			case systemFunctionDestroy:
				controller.Destroy(w)
			}
			continue
		}

		for j := range parallel {
			parallel[j].PrepareWg()
		}

		w.wg.Add(len(parallel))
		for j := range parallel {
			controller := parallel[j]

			switch method {
			case systemFunctionInit:
				controller.asyncInit()
			case systemFunctionUpdate:
				controller.asyncUpdate()
			case SystemFunctionFixedUpdate:
				controller.asyncFixedUpdate()
			case systemFunctionDestroy:
				controller.asyncDestroy()
			}
		}
	}
	w.wg.Wait()

	if method == SystemFunctionFixedUpdate {
		w.tick++
	}

	w.Clean()

	return nil
}

func (w *World) CreateEntity(title string) EntityID {
	var newId = w.generateEntityID()

	w.entityComponentMask.Set(newId, ComponentBitArray256{})
	w.size++

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

func (w *World) Clean() {
	for i := range w.components {
		w.components[i].Clean()
	}
}

func (w *World) Size() uint32 {
	return w.size
}

func (w *World) LastEntityID() EntityID {
	return w.lastEntityID
}

func (w *World) ShouldDestroy() bool {
	return w.shouldDestroy
}

func (w *World) SetShouldDestroy(value bool) {
	w.shouldDestroy = value
}

func (w *World) Destroy() {
	w.runSystemFunction(systemFunctionDestroy)
	w.wg.Wait()
	w.Clean()
}

func (w *World) Run(tickrate uint) {
	ticker := time.NewTicker(time.Second / time.Duration(tickrate))

	for !w.ShouldDestroy() {
		select {
		case <-ticker.C:
			w.runSystemFunction(SystemFunctionFixedUpdate)

			if len(ticker.C) > 0 {
				<-ticker.C
				log.Println("Skipping tick")
			}
		default:
			w.runSystemFunction(systemFunctionUpdate)
		}
	}
}

func (w *World) generateEntityID() (newId EntityID) {
	if len(w.deletedEntityIDs) == 0 {
		newId = EntityID(atomic.AddInt32((*int32)(&w.lastEntityID), 1))
	} else {
		newId = w.deletedEntityIDs[len(w.deletedEntityIDs)-1]
		w.deletedEntityIDs = w.deletedEntityIDs[:len(w.deletedEntityIDs)-1]
	}
	return newId
}
