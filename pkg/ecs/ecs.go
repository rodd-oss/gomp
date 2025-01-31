/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"sync"
	"time"
)

type WorldID uint32

const (
	PREALLOC_DELETED_ENTITIES uint32 = 1 << 10
)

func CreateWorld(title string) World {
	id := generateWorldID()
	maskSet := NewSparseSet[ComponentBitArray256, Entity]()

	world := World{
		ID:                  id,
		Title:               title,
		wg:                  new(sync.WaitGroup),
		mx:                  new(sync.Mutex),
		deletedEntityIDs:    make([]Entity, 0, PREALLOC_DELETED_ENTITIES),
		entityComponentMask: &maskSet,
		lastUpdateAt:        time.Now(),
		lastFixedUpdateAt:   time.Now(),
		components:          make(map[ComponentID]AnyComponentManagerPtr),
	}

	return world
}

func CreateComponentService[T any](id ComponentID) ComponentService[T] {
	component := ComponentService[T]{
		id:       id,
		managers: make(map[*World]*ComponentManager[T]),
	}

	return component
}

func CreateSystemService[T AnySystemControllerPtr](controller T, dependsOn ...AnySystemServicePtr) SystemService[T] {
	return SystemService[T]{
		initValue: controller,
		dependsOn: dependsOn,
		instances: make(map[*World]*SystemServiceInstance),
	}
}
