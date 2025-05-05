/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"gomp/pkg/worker"
	"sync"

	"github.com/negrel/assert"
)

type SharedComponentInstanceId uint16

var _ AnyComponentManagerPtr = &SharedComponentManager[any]{}

func NewSharedComponentManager[T any](id ComponentId) SharedComponentManager[T] {
	newManager := SharedComponentManager[T]{
		components:          NewPagedArray[T](),
		instances:           NewPagedArray[SharedComponentInstanceId](),
		instanceToComponent: NewPagedMap[SharedComponentInstanceId, int](),
		entityToComponent:   NewPagedMap[Entity, int](),
		entities:            NewPagedArray[Entity](),
		references:          NewPagedArray[SharedComponentInstanceId](),
		lookup:              NewPagedMap[Entity, int](),

		id:            id,
		isInitialized: true,

		TrackChanges:    false,
		createdEntities: NewPagedArray[Entity](),
		patchedEntities: NewPagedArray[Entity](),
		deletedEntities: NewPagedArray[Entity](),
	}

	return newManager
}

type SharedComponentManager[T any] struct {
	mx                  sync.Mutex
	components          PagedArray[T]
	instances           PagedArray[SharedComponentInstanceId]
	instanceToComponent PagedMap[SharedComponentInstanceId, int] // value is components array index
	entityToComponent   PagedMap[Entity, int]

	// for itter returning 2 values
	entities   PagedArray[Entity]
	references PagedArray[SharedComponentInstanceId]
	lookup     PagedMap[Entity, int]

	entityManager           *EntityManager
	entityComponentBitTable *ComponentBitTable
	pool                    *worker.Pool

	id            ComponentId
	isInitialized bool

	// Patch

	TrackChanges    bool // Enable TrackChanges to track changes and add them to patch
	createdEntities PagedArray[Entity]
	patchedEntities PagedArray[Entity]
	deletedEntities PagedArray[Entity]

	encoder func([]T) []byte
	decoder func([]byte) []T
}

func (c *SharedComponentManager[T]) registerWorkerPool(pool *worker.Pool) {
	c.pool = pool
}

func (c *SharedComponentManager[T]) PatchAdd(entity Entity) {
	//TODO implement me
	panic("implement me")
}

func (c *SharedComponentManager[T]) PatchGet() ComponentPatch {
	//TODO implement me
	panic("implement me")
}

func (c *SharedComponentManager[T]) PatchApply(patch ComponentPatch) {
	//TODO implement me
	panic("implement me")
}

func (c *SharedComponentManager[T]) PatchReset() {
	//TODO implement me
	panic("implement me")
}

func (c *SharedComponentManager[T]) Id() ComponentId {
	return c.id
}

func (c *SharedComponentManager[T]) registerEntityManager(entityManager *EntityManager) {
	c.entityManager = entityManager
	c.entityComponentBitTable = &entityManager.componentBitTable
}

//=====================================
//=====================================
//=====================================

// Create an instance of shared component
func (c *SharedComponentManager[T]) Create(instanceId SharedComponentInstanceId, value T) *T {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.assertBegin()
	defer c.assertEnd()

	componentIndex := c.components.Len()
	component := c.components.Append(value)
	c.instances.Append(instanceId)
	c.instanceToComponent.Set(instanceId, componentIndex)

	return component
}

func (c *SharedComponentManager[T]) Get(entity Entity) (component *T) {
	assert.True(c.isInitialized, "SharedComponentManager should be created with SharedNewComponentManager()")
	index, exists := c.entityToComponent.Get(entity)
	if !exists {
		return nil
	}
	component = c.components.Get(index)
	return component
}

func (c *SharedComponentManager[T]) GetComponentByInstance(instanceId SharedComponentInstanceId) (component *T) {
	assert.True(c.isInitialized, "SharedComponentManager should be created with SharedNewComponentManager()")
	index, exists := c.instanceToComponent.Get(instanceId)
	if !exists {
		return nil
	}
	component = c.components.Get(index)
	return component
}

func (c *SharedComponentManager[T]) GetInstanceByEntity(entity Entity) (SharedComponentInstanceId, bool) {
	assert.True(c.isInitialized, "SharedComponentManager should be created with SharedNewComponentManager()")
	index, exists := c.lookup.Get(entity)
	if !exists {
		return 0, false
	}
	instanceId := *c.instances.Get(index)
	return instanceId, true
}

func (c *SharedComponentManager[T]) Set(entity Entity, instanceId SharedComponentInstanceId) *T {
	assert.True(c.isInitialized, "SharedComponentManager should be created with SharedNewComponentManager()")
	index, exists := c.lookup.Get(entity)
	if exists {
		c.references.Set(index, instanceId)
	} else {
		newIndex := c.entities.Len()
		c.entities.Append(entity)
		c.references.Append(instanceId)
		c.lookup.Set(entity, newIndex)
	}
	componentIndex, _ := c.instanceToComponent.Get(instanceId)
	c.entityToComponent.Set(entity, componentIndex)
	c.patchedEntities.Append(entity)
	c.entityComponentBitTable.Set(entity, c.id)
	return c.components.Get(componentIndex)
}

func (c *SharedComponentManager[T]) Delete(entity Entity) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.assertBegin()
	defer c.assertEnd()

	index, _ := c.lookup.Get(entity)

	lastIndex := c.references.Len() - 1
	if index < lastIndex {
		// Swap the dead element with the last one
		c.references.Swap(index, lastIndex)
		newSwappedEntityId, _ := c.entities.Swap(index, lastIndex)
		assert.True(newSwappedEntityId != nil)

		// Update the lookup table
		c.lookup.Set(*newSwappedEntityId, index)
	}

	// Shrink the container
	c.references.SoftReduce()
	c.entities.SoftReduce()

	c.lookup.Delete(entity)
	c.entityToComponent.Delete(entity)
	c.entityComponentBitTable.Unset(entity, c.id)

	c.deletedEntities.Append(entity)
}

func (c *SharedComponentManager[T]) Has(entity Entity) bool {
	_, ok := c.lookup.Get(entity)
	return ok
}

func (c *SharedComponentManager[T]) Len() int {
	assert.True(c.isInitialized, "SharedComponentManager should be created with CreateComponentService()")
	return c.entities.Len()
}

func (c *SharedComponentManager[T]) Clean() {
	// c.entityComponentBitSet.Clean()
	//c.components.Clean()
	// c.Entities.Clean()
}

// ========================================================
// Iterators
// ========================================================

func (c *SharedComponentManager[T]) EachComponent() func(yield func(*T) bool) {
	c.assertBegin()
	defer c.assertEnd()
	return c.components.EachData()
}

func (c *SharedComponentManager[T]) EachEntity() func(yield func(Entity) bool) {
	c.assertBegin()
	defer c.assertEnd()
	return c.entities.EachDataValue()
}

func (c *SharedComponentManager[T]) Each() func(yield func(Entity, *T) bool) {
	c.assertBegin()
	defer c.assertEnd()
	return func(yield func(Entity, *T) bool) {
		c.components.Each()(func(i int, d *T) bool {
			entity := c.entities.Get(i)
			entId := *entity
			shouldContinue := yield(entId, d)
			return shouldContinue
		})
	}
}

// ========================================================
// Iterators Parallel
// ========================================================

func (c *SharedComponentManager[T]) ProcessComponents(handler func(*T, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.components.ProcessData(handler, c.pool)
}

func (c *SharedComponentManager[T]) EachEntityParallel(handler func(Entity, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.entities.ProcessDataValue(handler, c.pool)
}

func (c *SharedComponentManager[T]) EachParallel(numWorkers int) func(yield func(Entity, *T, int) bool) {
	c.assertBegin()
	defer c.assertEnd()
	return func(yield func(Entity, *T, int) bool) {
		c.components.EachParallel(numWorkers)(func(i int, t *T, workerId int) bool {
			entity := c.entities.Get(i)
			entId := *entity
			shouldContinue := yield(entId, t, workerId)
			return shouldContinue
		})
	}
}

// ========================================================
// Patches
// ========================================================

func (c *SharedComponentManager[T]) IsTrackingChanges() bool {
	return c.TrackChanges
}

// ========================================================
// Utils
// ========================================================

func (c *SharedComponentManager[T]) RawComponents(ptr []T) {
	c.components.Raw(ptr)
}

func (c *SharedComponentManager[T]) assertBegin() {
	assert.True(c.isInitialized, "SharedComponentManager should be created with SharedNewComponentManager()")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *SharedComponentManager[T]) assertEnd() {
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}
