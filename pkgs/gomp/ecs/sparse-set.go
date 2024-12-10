/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "iter"

type SparseSet[TData any, TKey EntityID | ComponentID | ECSID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse     *ChunkMap[int]
	denseData  *ChunkArray[TData]
	denseIndex *ChunkArray[TKey]
}

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int]() SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{}
	set.sparse = NewChunkMap[int](5, 10)
	set.denseData = NewChunkArray[TData](5, 10)
	set.denseIndex = NewChunkArray[TKey](5, 10)

	return set
}

func (s *SparseSet[TData, TKey]) Set(id TKey, data TData) *TData {
	pos, ok := s.sparse.Get(int(id))
	if ok {
		d, _ := s.denseData.Set(pos, data)
		return d
	}

	idx, r := s.denseData.Append(data)
	s.denseIndex.Append(id)
	s.sparse.Set(int(id), idx)

	return r
}

func (s *SparseSet[TData, TKey]) Get(id TKey) (data TData, ok bool) {
	index, ok := s.sparse.Get(int(id))

	if !ok {
		return data, false
	}

	el := s.denseData.GetPtr(index)
	if el == nil {
		return data, false
	}

	return *el, true
}

func (s *SparseSet[TData, TKey]) GetPtr(id TKey) *TData {
	index, ok := s.sparse.Get(int(id))
	if !ok {
		return nil
	}

	return s.denseData.GetPtr(index)
}

func (s *SparseSet[TData, TKey]) All() iter.Seq2[TKey, *TData] {
	return s.yielderAll
}

func (s *SparseSet[TData, TKey]) yielderAll(yield func(TKey, *TData) bool) {
	var indexBuffer = s.denseIndex.buffer
	var denseData = s.denseData

	for i, v := range denseData.All() {
		if !yield(indexBuffer[i.page].data[i.local], v) {
			return
		}
	}
}

func (s *SparseSet[TData, TKey]) SoftDelete(id TKey) {
	idx := int(id)

	indexx, ok := s.sparse.Get(idx)
	if !ok {
		return
	}

	lastDenseId, backEntityId, ok := s.denseIndex.Last()
	if !ok {
		return
	}

	s.denseData.Swap(indexx, lastDenseId)
	s.denseIndex.Swap(indexx, lastDenseId)

	s.sparse.Set(int(backEntityId), indexx)

	s.sparse.Delete(idx)

	s.denseData.SoftReduce()
	s.denseIndex.SoftReduce()
}

func (s *SparseSet[TData, TKey]) Clean() {
	s.denseData.Clean()
	s.denseIndex.Clean()
}
