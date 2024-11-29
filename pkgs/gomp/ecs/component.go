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

type Component[T any] struct {
	IDs       map[*ECS]ComponentID
	Instances map[*ECS]*SparseSet[ComponentInstance[T], EntityID]

	wg            *sync.WaitGroup
	parallelCount int
	dataChunkLen  int
	mx            *sync.Mutex
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

func (c *Component[T]) Get(entity *Entity) *T {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.Instances[entity.ecs]; !ok {
		panic(fmt.Sprintf("Component <%T> is not registered in <%s> world for <%d> entity", c, entity.ecs.Title, entity.ID))
	}

	if c.Instances[entity.ecs].Get(entity.ID) == nil {
		return nil
	}

	return &c.Instances[entity.ecs].Get(entity.ID).Data
}

// To use more threads we need to prespawn goroutines for each component
// var threads = runtime.NumCPU() * 2
const threads = 2

func (c *Component[T]) Each(ecs *ECS, callback func(*Entity, *T)) {
	c.parallelCount = threads

	arr := c.Instances[ecs].dense
	for _, b := range arr.buckets {
		for {
			c.dataChunkLen = len(b.data) / c.parallelCount
			if c.dataChunkLen != 0 {
				break
			}

			c.parallelCount = c.parallelCount / 2
			if c.parallelCount == 0 {
				panic("Can't split data into chunks")
			}
		}

		if c.parallelCount == 1 {
			for j := 0; j < len(b.data); j++ {
				callback(b.data[j].Entity, &b.data[j].Data)
			}
			continue
		}

		c.wg.Add(c.parallelCount)
		for i := 0; i < c.parallelCount; i++ {
			go c.parallelCallback(callback, b.data, i)
		}
		c.wg.Wait()
	}
}

func (c *Component[T]) parallelCallback(callback func(*Entity, *T), data []ComponentInstance[T], i int) {
	for j := i * c.dataChunkLen; j < (i+1)*c.dataChunkLen; j++ {
		callback(data[j].Entity, &data[j].Data)
	}
	c.wg.Done()
}

func (c *Component[T]) register(ecs *ECS) {
	c.IDs[ecs] = ecs.generateComponentID()
	set := NewSparseSet[ComponentInstance[T], EntityID](PREALLOC_BUCKETS, PREALLOC_BUCKETS_SIZE)
	c.Instances[ecs] = &set
}
