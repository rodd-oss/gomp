/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type SparseSet[TData any, TKey EntityID | ComponentID | WorldID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse     *ChunkMap[int]
	denseData  *ChunkArray[TData]
	denseIndex *ChunkArray[TKey]
}

func NewSparseSet[TData any, TKey EntityID | ComponentID | WorldID | int]() SparseSet[TData, TKey] {
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

	el := s.denseData.Get(index)
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

	return s.denseData.Get(index)
}

func (s *SparseSet[TData, TKey]) All(yield func(TKey, *TData) bool) {
	var denseData = s.denseData

	denseData.All(func(i ChunkArrayIndex, value *TData) bool {
		key := s.denseIndex.Get(i.globalOffset + i.local)
		if key == nil || value == nil {
			return true
		}
		return yield(*key, value)
	})
}

func (s *SparseSet[TData, TKey]) AllParallel(yield func(TKey, *TData) bool) {
	var denseData = s.denseData

	denseData.AllParallel(func(i ChunkArrayIndex, value *TData) bool {
		key := s.denseIndex.Get(i.globalOffset + i.local)
		if key == nil || value == nil {
			return true
		}
		return yield(*key, value)
	})
}

func (s *SparseSet[TData, TKey]) AllData(yield func(*TData) bool) {
	var denseData = s.denseData
	denseData.All(func(i ChunkArrayIndex, value *TData) bool {
		return yield(value)
	})
}

func (s *SparseSet[TData, TKey]) AllDataParallel(yield func(*TData) bool) {
	var denseData = s.denseData
	denseData.AllDataParallel(yield)
}

func (s *SparseSet[TData, TKey]) SoftDelete(id TKey) {
	//Get dense array index of the element to be deleted
	dataIndex, ok := s.sparse.Get(int(id))
	if !ok {
		return
	}

	// Get the last dense id
	lastDenseId, backSparseId, ok := s.denseIndex.Last()
	if !ok {
		return
	}

	s.denseData.Swap(lastDenseId, dataIndex)
	s.denseIndex.Swap(lastDenseId, dataIndex)
	s.sparse.SwapData(int(backSparseId), int(id))
	s.sparse.Delete(int(id))

	s.denseData.SoftReduce()
	s.denseIndex.SoftReduce()
}

func (s *SparseSet[TData, TKey]) Clean() {
	s.denseData.Clean()
	s.denseIndex.Clean()
}

func (s *SparseSet[TData, TKey]) Len() int {
	return s.denseData.Len()
}
