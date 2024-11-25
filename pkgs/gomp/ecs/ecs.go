/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type ECS struct {
	Entities SparseSet[Entity, EntityID]

	nextEntityID    EntityID
	nextComponentID ComponentID
}

func New() ECS {
	ecs := ECS{
		Entities: NewSparseSet[Entity, EntityID](),

		nextEntityID:    0,
		nextComponentID: 0,
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

func (e *ECS) CreateEntity() *Entity {
	entity := Entity{}
	entity.ID = e.generateEntityID()
	entity.ComponentsMask = NewBitArray(64)

	e.Entities.Add(entity.ID, entity)
	return e.Entities.Get(entity.ID)
}
