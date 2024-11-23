/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

const NULL_INDEX = -1

func NewSparseSet[TKey EntityID, TData any]() SparseSet[TKey, TData] {
	return SparseSet[TKey, TData]{
		sparse: make(map[TKey]int),
		dense:  make([]TData, 0),
	}
}

type SparseSet[TKey EntityID, TData any] struct {
	sparse map[TKey]int
	dense  []TData
}

func (s *SparseSet[TKey, TData]) Add(id TKey, data TData) {
	s.sparse[id] = len(s.dense)
	s.dense = append(s.dense, data)
}

func (s *SparseSet[TKey, TData]) Get(id TKey) *TData {
	i, ok := s.sparse[id]
	if !ok || i == NULL_INDEX {
		return nil
	}

	return &s.dense[i]
}

func (s *SparseSet[TKey, TData]) Delete(id TKey) {

	i, ok := s.sparse[id]
	if !ok {
		return
	}

	var lastEntity = TKey(len(s.sparse) - 1)
	s.dense[i] = s.dense[len(s.dense)]

	s.sparse[id] = NULL_INDEX
	s.sparse[lastEntity] = i

	s.dense = s.dense[:len(s.dense)-1]
}
