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
)

type EntityManager struct {
	lastId Entity
	size   uint32

	components          map[ComponentID]AnyComponentManagerPtr
	deletedEntityIDs    []Entity
	entityComponentMask *SparseSet[ComponentBitArray256, Entity]
	mx                  sync.Mutex

	patch Patch
}

type Patch []ComponentPatch

func (w *EntityManager) RegisterComponentServices(componentPtr ...AnyComponentServicePtr) {
	for i := 0; i < len(componentPtr); i++ {
		w.components[componentPtr[i].getId()] = componentPtr[i].register(w)
	}
}

func (w *EntityManager) Create() Entity {
	var newId = w.generateEntityID()

	w.entityComponentMask.Set(newId, ComponentBitArray256{})
	w.size++

	return newId
}

func (w *EntityManager) Delete(entity Entity) {
	mask := w.entityComponentMask.GetPtr(entity)
	if mask == nil {
		panic(fmt.Sprintf("Entity %d does not exist", entity))
	}

	for i := range mask.AllSet {
		w.components[i].Remove(entity)
	}

	w.entityComponentMask.SoftDelete(entity)
	w.deletedEntityIDs = append(w.deletedEntityIDs, entity)
	w.size--
}

func (w *EntityManager) Clean() {
	for i := range w.components {
		w.components[i].Clean()
	}
}

func (w *EntityManager) Size() uint32 {
	return w.size
}

func (w *EntityManager) LastId() Entity {
	return w.lastId
}

func (w *EntityManager) Destroy() {
	w.Clean()
}

func (w *EntityManager) PatchGet() Patch {
	patch := w.patch

	for _, component := range w.components {
		if !component.IsTrackingChanges() {
			continue
		}

		w.patch[component.getId()] = component.PatchGet()
	}

	return patch
}

func (w *EntityManager) PatchApply(patch Patch) {
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

func (w *EntityManager) PatchReset() {
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

func (w *EntityManager) init() {
	w.patch = make(Patch, len(w.components))
}

func (w *EntityManager) generateEntityID() (newId Entity) {
	if len(w.deletedEntityIDs) == 0 {
		newId = Entity(atomic.AddUint32((*entityType)(&w.lastId), 1))
	} else {
		newId = w.deletedEntityIDs[len(w.deletedEntityIDs)-1]
		w.deletedEntityIDs = w.deletedEntityIDs[:len(w.deletedEntityIDs)-1]
	}
	return newId
}
