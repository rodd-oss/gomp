/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

const (
	PREALLOC_DELETED_ENTITIES uint32 = 1 << 10
)

func NewEntityManager() EntityManager {
	maskSet := NewSparseSet[ComponentBitArray256, Entity]()

	entityManager := EntityManager{
		deletedEntityIDs:    make([]Entity, 0, PREALLOC_DELETED_ENTITIES),
		components:          make(map[ComponentID]AnyComponentManagerPtr),
		entityComponentMask: &maskSet,
	}

	return entityManager
}

func CreateComponentService[T any](id ComponentID) ComponentService[T] {
	component := ComponentService[T]{
		id:       id,
		managers: make(map[*EntityManager]*ComponentManager[T]),
	}

	return component
}
