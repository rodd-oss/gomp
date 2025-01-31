/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"iter"
	"math/big"
	"sync"

	"github.com/negrel/assert"
)

type ComponentManagerX[T any] struct {
	mx            *sync.Mutex
	components    *PagedArray[T]
	entities      *PagedArray[Entity]
	lookup        *PagedMap[Entity, int32]
	ID            ComponentID
	isInitialized bool
}

type ComponentSeviceX[T any] struct {
	id       ComponentID
	managers map[*World]*ComponentManagerX[T]
}

func (m *ComponentSeviceX[T]) Instance(world *World) *ComponentManagerX[T] {
	instance, ok := m.managers[world]
	assert.True(ok)
	return instance
}

func CreateComponentServiceX[T any](id ComponentID) *ComponentSeviceX[T] {
	return &ComponentSeviceX[T]{
		id:       id,
		managers: make(map[*World]*ComponentManagerX[T]),
	}
}

func (c *ComponentManagerX[T]) registerComponentMask(mask *ComponentManagerX[big.Int]) {
	// c.worldMask = mask
}

func (c *ComponentManagerX[T]) getId() ComponentID {
	return c.ID
}

func (c *ComponentManagerX[T]) Create(entity Entity, value T) (returnValue *T) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// Only one of component per enity allowed!
	assert.False(c.Has(entity))

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())

	var index = c.components.Len()

	c.lookup.Set(entity, index)
	c.entities.Append(entity)
	return c.components.Append(value)
}

func (c *ComponentManagerX[T]) Get(entity Entity) *T {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	return c.components.Get(index)
}

func (c *ComponentManagerX[T]) Remove(entity Entity) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// ENTITY HAS NO COMPONENT!
	assert.True(c.Has(entity))

	index, _ := c.lookup.Get(entity)

	lastIndex := c.components.Len() - 1
	if index < lastIndex {
		// Swap the the dead element with the last one
		c.components.Swap(index, lastIndex)
		newSwappedEntityId, _ := c.entities.Swap(index, lastIndex)
		assert.True(newSwappedEntityId != nil)

		// Update the lookup table
		c.lookup.Set(*newSwappedEntityId, index)
	}

	// Shrink the container
	c.components.SoftReduce()
	c.entities.SoftReduce()

	c.lookup.Delete(entity)

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())
}

func (c *ComponentManagerX[T]) Has(entity Entity) bool {
	_, ok := c.lookup.Get(entity)
	return ok
}

func (c *ComponentManagerX[T]) All(yield func(Entity, *T) bool) {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())

	nextData, stopData := iter.Pull(c.components.AllData)
	defer stopData()

	nextEntity, stopEntity := iter.Pull(c.entities.AllData)
	defer stopEntity()

	for {
		data, ok := nextData()
		if !ok {
			break
		}
		entityId, ok := nextEntity()
		if !ok {
			break
		}
		assert.True(entityId != nil)
		entId := *entityId
		shouldContinue := yield(entId, data)
		if !shouldContinue {
			break
		}
	}
}

func (c *ComponentManagerX[T]) AllParallel(yield func(Entity, *T) bool) {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())

	c.components.AllParallel(func(i int32, t *T) bool {
		entId := c.entities.Get(i)
		assert.True(entId != nil)
		return yield(*entId, t)
	})
}

func (c *ComponentManagerX[T]) AllData(yield func(*T) bool) {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())

	c.components.AllData(yield)
}

func (c *ComponentManagerX[T]) AllDataParallel(yield func(*T) bool) {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	// Entity Count must always be the same as the number of components!
	assert.True(c.entities.Len() == c.components.Len())
	assert.True(c.components.Len() == c.lookup.Len())

	c.components.AllDataParallel(yield)
}

func (c *ComponentManagerX[T]) Len() int32 {
	// ComponentManager must be initialized with CreateComponentManager()
	assert.True(c.isInitialized)

	return c.components.Len()
}

func (c *ComponentManagerX[T]) Clean() {
	// c.components.Clean()
	// c.entities.Clean()
}
