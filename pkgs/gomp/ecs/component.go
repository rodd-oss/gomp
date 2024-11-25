/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Component[T any] struct {
	ecs       *ECS
	ID        ComponentID
	Instances SparseSet[T, EntityID]
}

func CreateComponent[T any](ecs *ECS) Component[T] {
	component := Component[T]{}

	component.ID = ecs.generateComponentID()
	component.Instances = NewSparseSet[T, EntityID]()
	component.ecs = ecs

	return component
}

func (c *Component[T]) Set(entity *Entity, data T) *T {
	c.Instances.Add(entity.ID, data)
	c.ecs.Entities.Get(entity.ID).ComponentsMask.Set(uint(c.ID))
	return c.Instances.Get(entity.ID)
}
