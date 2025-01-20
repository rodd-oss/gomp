/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/negrel/assert"
)

// ================
// Contracts
// ================

type AnyComponentServicePtr interface {
	register(*World, ComponentID) AnyComponentManagerPtr
}

type AnyComponentManagerPtr interface {
	registerComponentMask(mask *ComponentManager[big.Int])
	getId() ComponentID
	Remove(EntityID)
	Clean()
	Has(EntityID) bool
}

// ================
// Service
// ================

type ComponentService[T any] struct {
	id       ComponentID
	managers map[*World]*ComponentManager[T]
}

func (c *ComponentService[T]) GetManager(world *World) *ComponentManager[T] {
	manager, ok := c.managers[world]
	assert.True(ok, fmt.Sprintf("Component <%T> is not registered in <%s> world", c, world.Title))
	return manager
}

func (c *ComponentService[T]) register(world *World, id ComponentID) AnyComponentManagerPtr {
	newManager := ComponentManager[T]{
		mx: new(sync.Mutex),

		components: NewPagedArray[T](),
		entities:   NewPagedArray[EntityID](),
		lookup:     NewPagedMap[EntityID, int32](),

		maskComponent: world.entityComponentMask,
		id:            id,
		isInitialized: true,
	}

	c.managers[world] = &newManager

	return &newManager
}

// ================
// Service
// ================

type ComponentManager[T any] struct {
	mx *sync.Mutex

	components *PagedArray[T]
	entities   *PagedArray[EntityID]
	lookup     *PagedMap[EntityID, int32]

	maskComponent *SparseSet[ComponentBitArray256, EntityID]
	id            ComponentID
	isInitialized bool
}

func (c *ComponentManager[T]) getId() ComponentID {
	return c.id
}

func (c *ComponentManager[T]) registerComponentMask(mask *ComponentManager[big.Int]) {
}

//=====================================
//=====================================
//=====================================

func (c *ComponentManager[T]) Create(entity EntityID, value T) (component *T) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(entity != -1, "INVALID ENTITY!")
	assert.False(c.Has(entity), "Only one of component per enity allowed!")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	var index = c.components.Len()

	c.lookup.Set(entity, index)
	c.entities.Append(entity)
	component = c.components.Append(value)

	mask := c.maskComponent.GetPtr(entity)
	mask.Set(c.id)

	return component
}

func (c *ComponentManager[T]) Get(entity EntityID) (component *T) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(entity != -1, "INVALID ENTITY!")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	return c.components.Get(index)
}

func (c *ComponentManager[T]) Remove(entity EntityID) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(entity != -1, "INVALID ENTITY!")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

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
	mask := c.maskComponent.GetPtr(entity)
	mask.Unset(c.id)

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) Has(entity EntityID) bool {
	_, ok := c.lookup.Get(entity)
	return ok
}

func (c *ComponentManager[T]) All(yield func(EntityID, *T) bool) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	c.components.All(func(i int32, d *T) bool {
		assert.True(d != nil)
		entity := c.entities.Get(i)
		assert.True(entity != nil)
		entId := *entity
		shouldContinue := yield(entId, d)
		return shouldContinue
	})

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) AllParallel(yield func(EntityID, *T) bool) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	c.components.AllParallel(func(i int32, t *T) bool {
		entity := c.entities.Get(i)
		assert.True(entity != nil)
		entId := *entity
		shouldContinue := yield(entId, t)
		return shouldContinue
	})

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) AllData(yield func(*T) bool) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	c.components.AllData(yield)

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) AllDataParallel(yield func(*T) bool) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	c.components.AllDataParallel(yield)

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) Len() int32 {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	return c.components.Len()
}

func (c *ComponentManager[T]) Clean() {
	c.maskComponent.Clean()
	// c.components.Clean()
	// c.entities.Clean()
}
