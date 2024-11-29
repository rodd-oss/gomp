/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"fmt"
	"sync"
)

type AnyComponentPtr interface {
	register(*ECS)
	Remove(*Entity)
}

type Component[T any] struct {
	IDs       map[*ECS]ComponentID
	Instances map[*ECS]*SparseSet[ComponentInstance[T], EntityID]

	wg *sync.WaitGroup
	mx *sync.Mutex
}

type ComponentInstance[T any] struct {
	Entity *Entity
	Data   T
}

func CreateComponent[T any]() Component[T] {
	component := Component[T]{}
	component.IDs = make(map[*ECS]ComponentID, 1)
	component.Instances = make(map[*ECS]*SparseSet[ComponentInstance[T], EntityID], 1)
	component.wg = new(sync.WaitGroup)
	component.mx = new(sync.Mutex)

	return component
}

func (c *Component[T]) Get(entity *Entity) *T {
	c.mx.Lock()
	defer c.mx.Unlock()

	if entity == nil {
		panic("Entity is nil")
	}

	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
	}

	if c.Instances[entity.ecs].Get(entity.ID) == nil {
		return nil
	}

	return &c.Instances[entity.ecs].Get(entity.ID).Data
}

func (c *Component[T]) Set(entity *Entity, data T) *T {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
	}

	// entity.ComponentsMask.Set(uint64(c.IDs[entity.ecs]))
	var instance ComponentInstance[T]
	instance.Entity = entity
	instance.Data = data
	return &c.Instances[entity.ecs].Set(entity.ID, instance).Data
}

func (c *Component[T]) Remove(entity *Entity) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
	}

	c.Instances[entity.ecs].Delete(entity.ID)
	// entity.ComponentsMask.Unset(uint64(c.IDs[entity.ecs]))
}

// To use more threads we need to prespawn goroutines for each component
// var threads = runtime.NumCPU() * 2

// TODO EachParallel()
func (c *Component[T]) Each(ecs *ECS, callback func(*Entity, *T)) {
	arr := c.Instances[ecs].dense.buckets
	for _, b := range arr {
		c.parallelCallback(callback, b.data)
	}
}

func (c *Component[T]) parallelCallback(callback func(*Entity, *T), data []DenseElement[ComponentInstance[T]]) {
	for j := 0; j < len(data); j++ {
		callback(data[j].value.Entity, &data[j].value.Data)
	}
}

func (c *Component[T]) register(ecs *ECS) {
	c.IDs[ecs] = ecs.generateComponentID()
	set := NewSparseSet[ComponentInstance[T], EntityID](PREALLOC_BUCKETS, PREALLOC_BUCKETS_SIZE)
	c.Instances[ecs] = &set
}
