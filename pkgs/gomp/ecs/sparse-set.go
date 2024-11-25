/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

func NewSparseSet[TData any, TKey EntityID | int]() SparseSet[TData, TKey] {
	return SparseSet[TData, TKey]{
		sparse: make(map[TKey]int),
		dense:  make([]TData, 0),
	}
}

type SparseSet[TData any, TKey EntityID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse map[TKey]int
	dense  []TData
}

func (s *SparseSet[TData, TKey]) Add(id TKey, data TData) {
	s.sparse[id] = len(s.dense)
	s.dense = append(s.dense, data)
}

func (s *SparseSet[TData, TKey]) Get(id TKey) *TData {
	i, ok := s.sparse[id]
	if !ok {
		return nil
	}

	return &s.dense[i]
}

func (s *SparseSet[TData, TKey]) Delete(id TKey) {
	i, ok := s.sparse[id]
	if !ok {
		return
	}

	var lastEntity = TKey(len(s.sparse) - 1)
	s.dense[i] = s.dense[len(s.dense)]

	delete(s.sparse, id)
	s.sparse[lastEntity] = i

	s.dense = s.dense[:len(s.dense)-1]
}
