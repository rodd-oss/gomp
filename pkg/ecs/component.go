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
	register(*World) AnyComponentManagerPtr
	getId() ComponentID
}

type AnyComponentManagerPtr interface {
	registerComponentMask(mask *ComponentManager[big.Int])
	getId() ComponentID
	Remove(Entity)
	Clean()
	Has(Entity) bool
	PatchAdd(Entity)
	PatchGet() ComponentPatch
	PatchApply(patch ComponentPatch)
	PatchReset()
	IsTrackingChanges() bool
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

func (c *ComponentService[T]) register(world *World) AnyComponentManagerPtr {
	newManager := ComponentManager[T]{
		mx: new(sync.Mutex),

		components: NewPagedArray[T](),
		entities:   NewPagedArray[Entity](),
		lookup:     NewPagedMap[Entity, int32](),

		maskComponent: world.entityComponentMask,
		id:            c.id,
		isInitialized: true,

		TrackChanges:    false,
		createdEntities: NewPagedArray[Entity](),
		patchedEntities: NewPagedArray[Entity](),
		deletedEntities: NewPagedArray[Entity](),
	}

	c.managers[world] = &newManager

	return &newManager
}

func (c *ComponentService[T]) getId() ComponentID {
	return c.id
}

// ================
// Service
// ================

type ComponentManager[T any] struct {
	mx *sync.Mutex

	components *PagedArray[T]
	entities   *PagedArray[Entity]
	lookup     *PagedMap[Entity, int32]

	maskComponent *SparseSet[ComponentBitArray256, Entity]
	id            ComponentID
	isInitialized bool

	// Patch

	TrackChanges    bool // Enable TrackChanges to track changes and add them to patch
	createdEntities *PagedArray[Entity]
	patchedEntities *PagedArray[Entity]
	deletedEntities *PagedArray[Entity]

	encoder func([]T) []byte
	decoder func([]byte) []T
}

// ComponentChanges with byte encoded Components
type ComponentChanges struct {
	Len        int32
	Components []byte
	Entities   []Entity
}

// ComponentPatch with byte encoded Created, Patched and Deleted components
type ComponentPatch struct {
	ID      ComponentID
	Created ComponentChanges
	Patched ComponentChanges
	Deleted ComponentChanges
}

func (c *ComponentManager[T]) getId() ComponentID {
	return c.id
}

func (c *ComponentManager[T]) registerComponentMask(*ComponentManager[big.Int]) {
}

//=====================================
//=====================================
//=====================================

func (c *ComponentManager[T]) Create(entity Entity, value T) (component *T) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.False(c.Has(entity), "Only one of component per entity allowed!")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	var index = c.components.Len()

	c.lookup.Set(entity, index)
	c.entities.Append(entity)
	component = c.components.Append(value)

	mask := c.maskComponent.GetPtr(entity)
	mask.Set(c.id)

	c.createdEntities.Append(entity)

	return component
}

func (c *ComponentManager[T]) Get(entity Entity) (component *T) {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	return c.components.Get(index)
}

func (c *ComponentManager[T]) Set(entity Entity, value T) *T {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	component := c.components.Set(index, value)

	c.patchedEntities.Append(entity)

	return component
}

func (c *ComponentManager[T]) Remove(entity Entity) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")

	index, _ := c.lookup.Get(entity)

	lastIndex := c.components.Len() - 1
	if index < lastIndex {
		// Swap the dead element with the last one
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

	c.deletedEntities.Append(entity)

	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) Has(entity Entity) bool {
	_, ok := c.lookup.Get(entity)
	return ok
}

// Patches

func (c *ComponentManager[T]) PatchAdd(entity Entity) {
	assert.True(c.TrackChanges)

	c.patchedEntities.Append(entity)
}

func (c *ComponentManager[T]) PatchGet() ComponentPatch {
	assert.True(c.TrackChanges)

	patch := ComponentPatch{
		ID:      c.id,
		Created: c.getChangesBinary(c.createdEntities),
		Patched: c.getChangesBinary(c.patchedEntities),
		Deleted: c.getChangesBinary(c.deletedEntities),
	}

	return patch
}

func (c *ComponentManager[T]) PatchApply(patch ComponentPatch) {
	assert.True(c.TrackChanges)
	assert.True(patch.ID == c.id)
	assert.True(c.decoder != nil)

	var components []T

	created := patch.Created
	components = c.decoder(created.Components)
	for i := range created.Len {
		c.Create(created.Entities[i], components[i])
	}

	patched := patch.Patched
	components = c.decoder(patched.Components)
	for i := range patched.Len {
		c.Set(patched.Entities[i], components[i])
	}

	deleted := patch.Deleted
	components = c.decoder(deleted.Components)
	for i := range deleted.Len {
		c.Remove(deleted.Entities[i])
	}
}

func (c *ComponentManager[T]) PatchReset() {
	assert.True(c.TrackChanges)

	c.createdEntities.Reset()
	c.patchedEntities.Reset()
	c.deletedEntities.Reset()
}

func (c *ComponentManager[T]) getChangesBinary(source *PagedArray[Entity]) ComponentChanges {
	changesLen := source.Len()

	components := make([]T, 0, changesLen)
	entities := make([]Entity, 0, changesLen)

	source.AllData(func(e *Entity) bool {
		assert.True(e != nil)
		entId := *e
		assert.True(c.Has(entId))
		components = append(components, *c.Get(entId))
		entities = append(entities, entId)
		return true
	})

	assert.True(c.encoder != nil)

	componentsBinary := c.encoder(components)

	return ComponentChanges{
		Len:        changesLen,
		Components: componentsBinary,
		Entities:   entities,
	}
}

func (c *ComponentManager[T]) SetEncoder(function func(components []T) []byte) *ComponentManager[T] {
	c.encoder = function
	return c
}

func (c *ComponentManager[T]) SetDecoder(function func(data []byte) []T) *ComponentManager[T] {
	c.decoder = function
	return c
}

func (c *ComponentManager[T]) IsTrackingChanges() bool {
	return c.TrackChanges
}

// Iterators

func (c *ComponentManager[T]) All(yield func(Entity, *T) bool) {
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

func (c *ComponentManager[T]) AllParallel(yield func(Entity, *T) bool) {
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

// Utils

func (c *ComponentManager[T]) Len() int32 {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")

	return c.components.Len()
}

func (c *ComponentManager[T]) Clean() {
	c.maskComponent.Clean()
	// c.components.Clean()
	// c.entities.Clean()
}
