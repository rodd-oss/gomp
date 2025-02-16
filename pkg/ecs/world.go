/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
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

func init() {
	nextWorldId.Store(-1)
}

var nextWorldId atomic.Int32

func generateWorldID() WorldID {
	id := nextWorldId.Add(1)
	return WorldID(id)
}

type World struct {
	ID    WorldID
	Title string

	tick         int
	lastEntityID Entity

	size          uint32
	shouldDestroy bool

	systems             [][]*SystemServiceInstance
	components          map[ComponentID]AnyComponentManagerPtr
	deletedEntityIDs    []Entity
	entityComponentMask *SparseSet[ComponentBitArray256, Entity]
	wg                  *sync.WaitGroup
	mx                  *sync.Mutex

	patch WorldPatch

	lastUpdateAt      time.Time
	lastFixedUpdateAt time.Time
}

type WorldPatch []ComponentPatch

func (w *World) RegisterComponentServices(componentPtr ...AnyComponentServicePtr) {
	for i := 0; i < len(componentPtr); i++ {
		w.components[componentPtr[i].getId()] = componentPtr[i].register(w)
	}
}

func (w *World) RegisterSystems() *SystemBuilder {
	return &SystemBuilder{
		world:   w,
		systems: &w.systems,
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

func (w *World) CreateEntity(title string) Entity {
	var newId = w.generateEntityID()

	w.entityComponentMask.Set(newId, ComponentBitArray256{})
	w.size++

	return newId
}

func (w *World) DestroyEntity(entityId Entity) {
	mask := w.entityComponentMask.GetPtr(entityId)
	if mask == nil {
		panic(fmt.Sprintf("Entity %d does not exist", entityId))
	}

	for i := range mask.AllSet {
		w.components[i].Remove(entityId)
	}

	w.entityComponentMask.SoftDelete(entityId)
	w.deletedEntityIDs = append(w.deletedEntityIDs, entityId)
	w.size--
}

func (w *World) Clean() {
	for i := range w.components {
		w.components[i].Clean()
	}
}

func (w *World) Size() uint32 {
	return w.size
}

func (w *World) LastEntityID() Entity {
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

func (w *World) PatchGet() WorldPatch {
	patch := w.patch

	for _, component := range w.components {
		if !component.IsTrackingChanges() {
			continue
		}

		w.patch[component.getId()] = component.PatchGet()
	}

	return patch
}

func (w *World) PatchApply(patch WorldPatch) {
	for _, componentPatch := range patch {
		component := w.components[componentPatch.ID]
		if component == nil {
			panic(fmt.Sprintf("Component %d does not exist", componentPatch.ID))
		}

		if !component.IsTrackingChanges() {
			continue
		}

		component.PatchApply(componentPatch)
	}
}

func (w *World) PatchReset() {
	for i, component := range w.components {
		if component == nil {
			panic(fmt.Sprintf("Component %d does not exist", i))
		}

		if !component.IsTrackingChanges() {
			continue
		}

		component.PatchReset()
	}
}

func (w *World) init() {
	w.patch = make(WorldPatch, len(w.components))
}

func (w *World) Run(tickrate uint) {
	w.init()

	duration := time.Second / time.Duration(tickrate)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	var t time.Time

	for !w.ShouldDestroy() {
		needFixedUpdate := true
		for needFixedUpdate {
			select {
			default:
				needFixedUpdate = false
			case <-ticker.C:
				t = time.Now()
				w.runSystemFunction(SystemFunctionFixedUpdate)
				w.lastFixedUpdateAt = t
			}
		}
		t = time.Now()
		w.runSystemFunction(systemFunctionUpdate)
		w.lastUpdateAt = t
	}
}

func (w *World) DtUpdate() time.Duration {
	return time.Since(w.lastUpdateAt)
}

func (w *World) DtFixedUpdate() time.Duration {
	return time.Since(w.lastFixedUpdateAt)
}

func (w *World) generateEntityID() (newId Entity) {
	if len(w.deletedEntityIDs) == 0 {
		newId = Entity(atomic.AddUint32((*entityType)(&w.lastEntityID), 1))
	} else {
		newId = w.deletedEntityIDs[len(w.deletedEntityIDs)-1]
		w.deletedEntityIDs = w.deletedEntityIDs[:len(w.deletedEntityIDs)-1]
	}
	return newId
}
