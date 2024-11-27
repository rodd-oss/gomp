/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "fmt"

type Component[T any] struct {
	IDs       map[*ECS]ComponentID
	Instances map[*ECS]*SparseSet[ComponentInstance[T], EntityID]
}

type ComponentInstance[T any] struct {
	Entity *Entity
	Data   T
}

func CreateComponent[T any]() Component[T] {
	component := Component[T]{}
	component.IDs = make(map[*ECS]ComponentID, 1)
	component.Instances = make(map[*ECS]*SparseSet[ComponentInstance[T], EntityID], 1)

	return component
}

func (c *Component[T]) Set(entity *Entity, data T) *T {
	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%s> entity", c, entity.ecs.Title, entity.Title))
	}

	entity.ComponentsMask.Set(uint64(c.IDs[entity.ecs]))
	var instance ComponentInstance[T]
	instance.Entity = entity
	instance.Data = data
	return &c.Instances[entity.ecs].Set(entity.ID, instance).Data
}

func (c *Component[T]) Get(entity *Entity) *T {
	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%s> entity", c, entity.ecs.Title, entity.Title))
	}

	if c.Instances[entity.ecs].Get(entity.ID) == nil {
		return nil
	}

	return &c.Instances[entity.ecs].Get(entity.ID).Data
}

func (c *Component[T]) Each(ecs *ECS, callback func(*Entity, *T)) {
	arr := c.Instances[ecs].dense
	for _, b := range arr.buckets {
		for i := range b.data {
			callback(b.data[i].Entity, &b.data[i].Data)
		}
	}
}

func (c *Component[T]) register(ecs *ECS) {
	c.IDs[ecs] = ecs.generateComponentID()
	set := NewSparseSet[ComponentInstance[T], EntityID](1000000)
	c.Instances[ecs] = &set
}
