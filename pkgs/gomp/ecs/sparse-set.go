/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

const NULL_INDEX = -1

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int](buckets uint32, bucketSize uint32) SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{}
	set.sparse.init(int(buckets), int(bucketSize), NULL_INDEX)
	set.dense.init(int(buckets), int(bucketSize), DenseElement[TData]{})

	return set
}

type DenseElement[TData any] struct {
	index int
	value TData
}

type SparseSet[TData any, TKey EntityID | ComponentID | ECSID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse Collection[int]
	dense  Collection[DenseElement[TData]]
}

func (s *SparseSet[TData, TKey]) Set(id TKey, data TData) *TData {
	var element = DenseElement[TData]{
		index: int(id),
		value: data,
	}

	pos := s.sparse.Get(int(id))
	if pos != NULL_INDEX {
		return &s.dense.Set(pos, element).value
	}

	idx, r := s.dense.Append(element)
	s.sparse.Set(int(id), idx)

	return &r.value
}

func (s *SparseSet[TData, TKey]) Get(id TKey) (data TData, ok bool) {
	index := s.sparse.Get(int(id))

	if index == NULL_INDEX {
		return data, false
	}

	return s.dense.Get(index).value, true
}

func (s *SparseSet[TData, TKey]) SoftDelete(id TKey) {
	idx := int(id)

	indexx := s.sparse.Get(idx)
	if indexx == NULL_INDEX {
		return
	}

	lastDenseId, lastDenseElement := s.dense.Last()
	backEntityId := lastDenseElement.index
	s.dense.Swap(indexx, lastDenseId)
	s.sparse.Set(backEntityId, indexx)

	s.sparse.Set(idx, NULL_INDEX)

	s.dense.SoftReduce()
}

func (s *SparseSet[TData, TKey]) Clean() {
	s.dense.Clean()
}
