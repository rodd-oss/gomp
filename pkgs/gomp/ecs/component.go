/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "fmt"

type Component[T any] struct {
	IDs       map[*ECS]ComponentID
	Instances map[*ECS]*SparseSet[T, EntityID]
}

func CreateComponent[T any]() Component[T] {
	component := Component[T]{}
	component.IDs = make(map[*ECS]ComponentID, 1)
	component.Instances = make(map[*ECS]*SparseSet[T, EntityID], 1)

	return component
}

func (c *Component[T]) Set(entity *Entity, data T) *T {
	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%s> entity", c, entity.ecs.Title, entity.Title))
	}

	entity.ComponentsMask.Set(uint(c.IDs[entity.ecs]))
	return c.Instances[entity.ecs].Set(entity.ID, data)
}

func (c *Component[T]) Get(entity *Entity) *T {
	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%s> entity", c, entity.ecs.Title, entity.Title))
	}

	return c.Instances[entity.ecs].Get(entity.ID)
}

func (c *Component[T]) Each(ecs *ECS, callback func(data *T)) {
	arr := c.Instances[ecs].dense
	for i := range arr {
		callback(&arr[i])
	}
}

func (c *Component[T]) Data(ecs *ECS) []T {
	return c.Instances[ecs].dense
}

func (c *Component[T]) register(ecs *ECS) {
	c.IDs[ecs] = ecs.generateComponentID()
	set := NewSparseSet[T, EntityID](1000000)
	c.Instances[ecs] = &set
}
