/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Hininn Donated 2 000 RUB
<- Сосисочник Паша Donated 77 RUB
<- mitwelve Donated 1 000 RUB

Thank you for your support!
*/

package ecs

import (
	"github.com/negrel/assert"
	"gomp/pkg/worker"
	"sync"
)

// ================
// Contracts
// ================

type ComponentId uint16
type AnyComponentList interface{}
type AnyComponentListPtr interface{}

type AnyComponentManagerPtr interface {
	Id() ComponentId
	Delete(Entity)
	Clean()
	Has(Entity) bool
	PatchAdd(Entity)
	PatchGet() ComponentPatch
	PatchApply(patch ComponentPatch)
	PatchReset()
	IsTrackingChanges() bool
	registerEntityManager(*EntityManager)
	registerWorkerPool(*worker.Pool)
}

// ================
// Service
// ================

var _ AnyComponentManagerPtr = &ComponentManager[any]{}

func NewComponentManager[T any](id ComponentId) ComponentManager[T] {
	newManager := ComponentManager[T]{
		components: NewPagedArray[T](),
		entities:   NewPagedArray[Entity](),
		lookup:     NewPagedMap[Entity, int](),

		id:            id,
		isInitialized: true,

		TrackChanges:    false,
		createdEntities: NewPagedArray[Entity](),
		patchedEntities: NewPagedArray[Entity](),
		deletedEntities: NewPagedArray[Entity](),
	}

	return newManager
}

type ComponentManager[T any] struct {
	components PagedArray[T]
	entities   PagedArray[Entity]
	lookup     PagedMap[Entity, int]

	entityManager           *EntityManager
	entityComponentBitTable *ComponentBitTable

	id            ComponentId
	isInitialized bool

	pool *worker.Pool

	// Patch

	TrackChanges    bool // Enable TrackChanges to track changes and add them to patch
	createdEntities PagedArray[Entity]
	patchedEntities PagedArray[Entity]
	deletedEntities PagedArray[Entity]

	encoder func([]T) []byte
	decoder func([]byte) []T

	mx sync.Mutex
}

// ComponentChanges with byte encoded Components
type ComponentChanges struct {
	Len        int
	Components []byte
	Entities   []Entity
}

// ComponentPatch with byte encoded Created, Patched and Deleted components
type ComponentPatch struct {
	ID      ComponentId
	Created ComponentChanges
	Patched ComponentChanges
	Deleted ComponentChanges
}

func (c *ComponentManager[T]) Id() ComponentId {
	return c.id
}

func (c *ComponentManager[T]) registerEntityManager(entityManager *EntityManager) {
	c.entityManager = entityManager
	c.entityComponentBitTable = &entityManager.componentBitTable
}

func (c *ComponentManager[T]) registerWorkerPool(pool *worker.Pool) {
	c.pool = pool
}

//=====================================
//=====================================
//=====================================

func (c *ComponentManager[T]) Create(entity Entity, value T) (component *T) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.False(c.Has(entity), "Only one of component per entity allowed!")
	c.assertBegin()
	defer c.assertEnd()

	var index = c.components.Len()

	c.lookup.Set(entity, index)
	c.entities.Append(entity)
	component = c.components.Append(value)

	c.entityComponentBitTable.Set(entity, c.id)

	//c.createdEntities.Append(entity)

	return component
}

/*
GetUnsafe - is not thread safe. DO NOT store the pointer to the value anywhere, because it might be changed anytime with Create or Delete operations.
*/
func (c *ComponentManager[T]) GetUnsafe(entity Entity) (component *T) {
	assert.True(c.isInitialized, "ComponentManager should be created with NewComponentManager()")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	return c.components.Get(index)
}

func (c *ComponentManager[T]) Get(entity Entity) (component T, ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with NewComponentManager()")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return component, false
	}

	return c.components.GetValue(index), true
}

func (c *ComponentManager[T]) Set(entity Entity, value T) *T {
	c.mx.Lock()
	defer c.mx.Unlock()

	assert.True(c.isInitialized, "ComponentManager should be created with NewComponentManager()")

	index, ok := c.lookup.Get(entity)
	if !ok {
		return nil
	}

	component := c.components.Set(index, value)

	c.patchedEntities.Append(entity)

	return component
}

func (c *ComponentManager[T]) Delete(entity Entity) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.assertBegin()
	defer c.assertEnd()

	index, exists := c.lookup.Get(entity)
	assert.True(exists, "Entity does not have component")

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
	c.entityComponentBitTable.Unset(entity, c.id)

	//c.deletedEntities.Append(entity)
}

func (c *ComponentManager[T]) Has(entity Entity) bool {
	_, ok := c.lookup.Get(entity)
	return ok
}

func (c *ComponentManager[T]) Len() int {
	assert.True(c.isInitialized, "ComponentManager should be created with CreateComponentService()")
	return c.components.Len()
}

func (c *ComponentManager[T]) Clean() {
	// c.entityComponentBitSet.Clean()
	//c.components.Clean()
	// c.Entities.Clean()
}

// ========================================================
// Iterators
// ========================================================

func (c *ComponentManager[T]) EachComponent(yield func(entity *T) bool) {
	c.assertBegin()
	defer c.assertEnd()
	c.components.EachData(yield)
}

func (c *ComponentManager[T]) EachEntity(yield func(entity Entity) bool) {
	c.assertBegin()
	defer c.assertEnd()
	c.entities.EachDataValue(yield)
}

func (c *ComponentManager[T]) Each() func(yield func(entity Entity, component *T) bool) {
	c.assertBegin()
	defer c.assertEnd()
	return func(yield func(entity Entity, component *T) bool) {
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

func (c *ComponentManager[T]) ProcessEntities(handler func(Entity, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.entities.ProcessDataValue(handler, c.pool)
}

func (c *ComponentManager[T]) ProcessComponents(handler func(*T, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.components.ProcessData(handler, c.pool)
}

func (c *ComponentManager[T]) EachParallel(numWorkers int) func(yield func(Entity, *T, int) bool) {
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

func (c *ComponentManager[T]) EachEntityParallel(handler func(Entity, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.entities.EachDataValueParallel(handler, c.pool)
}

func (c *ComponentManager[T]) EachComponentParallel(handler func(*T, worker.WorkerId)) {
	c.assertBegin()
	defer c.assertEnd()
	c.components.EachDataParallel(handler, c.pool)
}

// ========================================================
// Patches
// ========================================================

func (c *ComponentManager[T]) PatchAdd(entity Entity) {
	assert.True(c.TrackChanges)
	c.patchedEntities.Append(entity)
}

func (c *ComponentManager[T]) PatchGet() ComponentPatch {
	assert.True(c.TrackChanges)
	patch := ComponentPatch{
		ID:      c.id,
		Created: c.getChangesBinary(&c.createdEntities),
		Patched: c.getChangesBinary(&c.patchedEntities),
		Deleted: c.getChangesBinary(&c.deletedEntities),
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
		c.Delete(deleted.Entities[i])
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

	source.EachData(func(e *Entity) bool {
		assert.True(e != nil)
		entId := *e
		assert.True(c.Has(entId))
		components = append(components, *c.GetUnsafe(entId))
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

// ========================================================
// Utils
// ========================================================

func (c *ComponentManager[T]) RawComponents(ptr []T) []T {
	return c.components.Raw(ptr)
}

func (c *ComponentManager[T]) RawEntities(ptr []Entity) []Entity {
	return c.entities.Raw(ptr)
}

func (c *ComponentManager[T]) assertBegin() {
	assert.True(c.isInitialized, "ComponentManager should be created with NewComponentManager()")
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}

func (c *ComponentManager[T]) assertEnd() {
	assert.True(c.components.Len() == c.lookup.Len(), "Lookup Count must always be the same as the number of components!")
	assert.True(c.entities.Len() == c.components.Len(), "Entity Count must always be the same as the number of components!")
}
