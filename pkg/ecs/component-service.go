/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package ecs

import (
	"fmt"
	"github.com/negrel/assert"
	"sync"
)

type AnyComponentServicePtr interface {
	register(*EntityManager) AnyComponentManagerPtr
	getId() ComponentID
}

type ComponentService[T any] struct {
	id       ComponentID
	managers map[*EntityManager]*ComponentManager[T]
}

func (c *ComponentService[T]) GetManager(world *EntityManager) *ComponentManager[T] {
	manager, ok := c.managers[world]
	assert.True(ok, fmt.Sprintf("Component <%T> is not registered in world", c))
	return manager
}

func (c *ComponentService[T]) register(world *EntityManager) AnyComponentManagerPtr {
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
