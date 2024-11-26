/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "slices"

const NULL_INDEX int = -1

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int](size int) SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{
		sparse: make([]int, 0, size),
		dense:  make([]TData, 0, size),
		size:   size,
	}

	for i := 0; i < size; i++ {
		set.sparse = append(set.sparse, NULL_INDEX)
	}

	return set
}

type SparseSet[TData any, TKey EntityID | ComponentID | ECSID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse []int
	dense  []TData
	size   int
}

func (s *SparseSet[TData, TKey]) Set(id TKey, data TData) *TData {
	if int(id) >= len(s.sparse) {
		s.sparse = append(s.sparse, slices.Repeat([]int{NULL_INDEX}, s.size)...)
	}

	if id == 10000000 {

	}

	if s.sparse[id] == NULL_INDEX {
		s.sparse[id] = len(s.dense)
		s.dense = append(s.dense, data)
	} else {
		s.dense[s.sparse[id]] = data
	}

	return &s.dense[s.sparse[id]]
}

func (s *SparseSet[TData, TKey]) Get(id TKey) *TData {
	if int(id) >= len(s.sparse) {
		return nil
	}

	if s.sparse[id] == NULL_INDEX {
		return nil
	}

	return &s.dense[s.sparse[id]]
}

func (s *SparseSet[TData, TKey]) Delete(id TKey) {
	if int(id) >= len(s.sparse) {
		return
	}

	i := s.sparse[id]

	var lastEntity = TKey(len(s.sparse) - 1)
	s.dense[i] = s.dense[len(s.dense)]

	s.sparse[id] = NULL_INDEX
	s.sparse[lastEntity] = i

	s.dense = s.dense[:len(s.dense)-1]
}
