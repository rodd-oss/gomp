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

const (
	PREALLOC_DEFAULT uint32 = 1 << 10
)

func NewEntityManager() EntityManager {
	entityManager := EntityManager{
		deletedEntityIDs: make([]Entity, 0, PREALLOC_DEFAULT),
		components:       make(map[ComponentId]AnyComponentManagerPtr),
	}

	return entityManager
}

type EntityManager struct {
	lastId Entity
	size   uint32

	groups            map[string][]Entity
	components        map[ComponentId]AnyComponentManagerPtr
	deletedEntityIDs  []Entity
	componentBitTable ComponentBitTable
	mx                sync.Mutex

	patch Patch
}

type Patch []ComponentPatch

func (e *EntityManager) Create() Entity {
	e.mx.Lock()
	defer e.mx.Unlock()
	var newId = e.generateEntityID()
	e.componentBitTable.Create(newId)

	e.size++

	return newId
}

func (e *EntityManager) Delete(entity Entity) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.componentBitTable.AllSet(entity, func(id ComponentId) bool {
		e.components[id].Delete(entity)
		return true
	})
	e.componentBitTable.Delete(entity)

	e.deletedEntityIDs = append(e.deletedEntityIDs, entity)
	e.size--
}

func (e *EntityManager) Clean() {
	for i := range e.components {
		e.components[i].Clean()
	}
}

func (e *EntityManager) Size() uint32 {
	return e.size
}

func (e *EntityManager) LastId() Entity {
	return e.lastId
}

func (e *EntityManager) Destroy() {
	e.Clean()
}

func (e *EntityManager) PatchGet() Patch {
	patch := e.patch

	for _, component := range e.components {
		if !component.IsTrackingChanges() {
			continue
		}

		e.patch[component.Id()] = component.PatchGet()
	}

	return patch
}

func (e *EntityManager) PatchApply(patch Patch) {
	for _, componentPatch := range patch {
		component := e.components[componentPatch.ID]
		if component == nil {
			panic(fmt.Sprintf("Component %d does not exist", componentPatch.ID))
		}

		if !component.IsTrackingChanges() {
			continue
		}

		component.PatchApply(componentPatch)
	}
}

func (e *EntityManager) PatchReset() {
	for i, component := range e.components {
		if component == nil {
			panic(fmt.Sprintf("Component %d does not exist", i))
		}

		if !component.IsTrackingChanges() {
			continue
		}

		component.PatchReset()
	}
}

func (e *EntityManager) init() {
	e.componentBitTable = NewComponentBitTable(len(e.components))
	e.patch = make(Patch, len(e.components))
}

func (e *EntityManager) generateEntityID() (newId Entity) {
	if len(e.deletedEntityIDs) == 0 {
		newId = Entity(atomic.AddUint32((*entityType)(&e.lastId), 1))
	} else {
		newId = e.deletedEntityIDs[len(e.deletedEntityIDs)-1]
		e.deletedEntityIDs = e.deletedEntityIDs[:len(e.deletedEntityIDs)-1]
	}
	return newId
}

func (e *EntityManager) registerComponent(c AnyComponentManagerPtr) {
	e.components[c.Id()] = c
}
