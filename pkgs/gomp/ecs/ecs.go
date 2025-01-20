/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"sync"
)

type WorldID uint

const (
	PREALLOC_DELETED_ENTITIES uint32 = 1 << 10
)

var nextWorldId WorldID = 0

func generateWorldID() WorldID {
	id := nextWorldId
	nextWorldId++
	return id
}

func CreateWorld(title string) World {
	id := generateWorldID()
	maskSet := NewSparseSet[ComponentBitArray256, EntityID]()

	ecs := World{
		ID:                  id,
		Title:               title,
		wg:                  new(sync.WaitGroup),
		mx:                  new(sync.Mutex),
		deletedEntityIDs:    make([]EntityID, 0, PREALLOC_DELETED_ENTITIES),
		entityComponentMask: &maskSet,
	}

	return ecs
}

func CreateComponentService[T any](id ComponentID) ComponentService[T] {
	component := ComponentService[T]{
		id:       id,
		managers: make(map[*World]*ComponentManager[T]),
	}

	return component
}
