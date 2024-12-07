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
	denseIndex *ChunkArray[int]
}

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int](buckets uint32, bucketSize uint32) SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{}
	set.sparse = NewChunkMap[int](8, 1024)
	set.denseData = NewChunkArray[TData](3, 2)
	set.denseIndex = NewChunkArray[int](3, 2)

	return set
}

type DenseElement[TData any] struct {
	index int
	value TData
}

func (s *SparseSet[TData, TKey]) Set(id TKey, data TData) *TData {
	pos, ok := s.sparse.Get(int(id))
	if ok {
		d, _ := s.denseData.Set(pos, data)
		return d
	}

	idx, r := s.denseData.Append(data)
	s.denseIndex.Append(int(id))
	s.sparse.Set(int(id), idx)

	return r
}

func (s *SparseSet[TData, TKey]) Get(id TKey) (data TData, ok bool) {
	index, ok := s.sparse.Get(int(id))

	if !ok {
		return data, false
	}

	el, ok := s.denseData.Get(index)
	if !ok {
		return data, false
	}

	return el, true
}

func (s *SparseSet[TData, TKey]) Iter() iter.Seq2[int, *TData] {
	return s.denseData.Iter()
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
	s.sparse.Set(backEntityId, indexx)

	s.sparse.Delete(idx)

	s.denseData.SoftReduce()
}

func (s *SparseSet[TData, TKey]) Clean() {
	s.denseData.Clean()
}
