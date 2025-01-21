/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file deveopment:
-===-===-===-===-===-===-===-===-===-===

<- MuTaToR Donated 250 RUB
<- Бодрящий член Donated 100 RUB
<- Plambirik Donated 5 000 RUB
<- Бодрящий член Donated 100 RUB
<- MuTaToR Donated 250 RUB
<- ksana_pro Donated 100 RUB
<- Skomich Donated 250 RUB
<- MuTaToR Donated 250 RUB
<- Бодрящий член Donated 100 RUB
<- мой код полная хуйня Donated 251 RUB
<- ksana_pro Donated 100 RUB
<- дубина Donated 250 RUB
<- WoWnik Donated 100 RUB
<- Vorobyan Donated 100 RUB
<- MuTaToR Donated 250 RUB
<- Мандовожка милана Donated 100 RUB
<- ksana_pro Donated 100 RUB
<- Зритель Donated 250 RUB
<- Ричард Donated 100 RUB
<- ksana_pro Donated 100 RUB
<- Ksana_pro Donated 100 RUB
<- Ksana_pro Donated 100 RUB
<- Ksana_pro Donated 100 RUB
<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package ecs

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type World struct {
	ID    WorldID
	Title string

	tick         int
	lastEntityID EntityID

	size          uint32
	shouldDestroy bool

	systems             [][]*SystemServiceInstance
	components          []AnyComponentManagerPtr
	deletedEntityIDs    []EntityID
	entityComponentMask *SparseSet[ComponentBitArray256, EntityID]
	wg                  *sync.WaitGroup
	mx                  *sync.Mutex

	lastUpdateAt      time.Time
	lastFixedUpdateAt time.Time
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
		w.wg.Wait()
	}

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
	duration := time.Second / time.Duration(tickrate)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for !w.ShouldDestroy() {
		needFixedUpdate := true
		for needFixedUpdate {
			select {
			default:
				needFixedUpdate = false
			case <-ticker.C:
				w.runSystemFunction(SystemFunctionFixedUpdate)
				w.lastFixedUpdateAt = time.Now()
			}
		}
		w.runSystemFunction(systemFunctionUpdate)
		w.lastUpdateAt = time.Now()
	}
}

func (w *World) DtUpdate() time.Duration {
	return time.Since(w.lastUpdateAt)
}

func (w *World) DtFixedUpdate() time.Duration {
	return time.Since(w.lastFixedUpdateAt)
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

var nextWorldId WorldID = 0

func generateWorldID() WorldID {
	id := nextWorldId
	nextWorldId++
	return id
}
