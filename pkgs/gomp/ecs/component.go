/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"sync"
)

type AnyComponentPtr interface {
	register(*ECS)
	SoftRemove(*Entity)
	Clean(*ECS)
}

type Component[T any] struct {
	IDs       []ComponentID
	instances *ChunkMap[SparseSet[T, EntityID]]

	wg *sync.WaitGroup
	mx *sync.Mutex
}

func CreateComponent[T any]() Component[T] {
	component := Component[T]{}
	component.IDs = make([]ComponentID, 0, 5)
	component.instances = NewChunkMap[SparseSet[T, EntityID]](2, 5)
	component.wg = new(sync.WaitGroup)
	component.mx = new(sync.Mutex)

	return component
}

func (c *Component[T]) Get(entity *Entity) (T, bool) {
	// if value, ok := c.Instances[entity.ecs.ID]; ok {
	value, _ := c.instances.Get(int(entity.ecsID))

	instance, ok := value.Get(entity.ID)

	return instance, ok
	// }

	// panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
}

func (c *Component[T]) Set(entity *Entity, data T) *T {
	// if value, ok := c.Instances[entity.ecs.id]; ok {
	value, _ := c.instances.Get(int(entity.ecsID))

	// entity.ComponentsMask.Set(uint64(c.IDs[entity.ecs]))
	var newinstance = value.Set(entity.ID, data)

	return newinstance
	// }

	// panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
}

func (c *Component[T]) SoftRemove(entity *Entity) {
	// if _, ok := c.Instances[entity.ecs.ID]; !ok {
	// 	panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
	// }
	value, _ := c.instances.Get(int(entity.ecsID))

	value.SoftDelete(entity.ID)
	// entity.ComponentsMask.Unset(uint64(c.IDs[entity.ecs]))
}

func (c *Component[T]) Clean(ecs *ECS) {
	value, _ := c.instances.Get(int(ecs.ID))
	value.Clean()
}

// To use more threads we need to prespawn goroutines for each component
// var threads = runtime.NumCPU() * 2

// TODO EachParallel()
// func (c *Component[T]) Each(ecs *ECS, callback func(*Entity, T)) {
// 	ecsInstances, _ := c.instances.Get(int(ecs.ID))

// 	for _, instance := range ecsInstances.Iter() {
// 		callback(instance.Entity, instance.Data)
// 	}

// 	// for _, s := range c.Instances {
// 	// 	s.Clean()
// 	// }
// }

func (c *Component[T]) Instances(ecs *ECS) SparseSet[T, EntityID] {
	ecsInstances, _ := c.instances.Get(int(ecs.ID))

	return ecsInstances
}

// func (c *Component[T]) EachParallel(ecs *ECS, callback func(*Entity, *T)) {
// 	arr := c.Instances[ecs.ID].denseData.buckets
// 	for _, b := range arr {
// 		c.parallelCallback(callback, b.data)
// 	}

// 	for _, s := range c.Instances {
// 		s.Clean()
// 	}
// }

// func (c *Component[T]) parallelCallback(callback func(*Entity, *T), data []ComponentInstance[T]) {
// 	for j := 0; j < len(data); j++ {
// 		callback(data[j].Entity, &data[j].Data)
// 	}
// }

func (c *Component[T]) register(ecs *ECS) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.IDs = append(c.IDs, ecs.generateComponentID())
	set := NewSparseSet[T, EntityID]()
	c.instances.Set(int(ecs.ID), set)
}
