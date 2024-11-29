/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

const NULL_INDEX = -1

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int](buckets uint32, bucketSize uint32) SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{}
	set.sparse.initialBucketsCount = buckets
	set.dense.initialBucketsCount = buckets
	set.sparse.initialBucketSize = bucketSize
	set.dense.initialBucketSize = bucketSize

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
	if !s.isDeleted(id) {
		pos := s.sparse.Get(int(id))
		if pos != nil {
			element := DenseElement[TData]{
				index: *pos,
				value: data,
			}
			return &s.dense.Set(*pos, element).value
		}
	}

	idx, r := s.dense.Append(DenseElement[TData]{})
	r.index = idx
	r.value = data

	s.sparse.Set(int(id), idx)
	return &r.value
}

func (s *SparseSet[TData, TKey]) Get(id TKey) *TData {
	if s.isDeleted(id) {
		return nil
	}

	idx := s.sparse.Get(int(id))
	if idx == nil {
		return nil
	}
	return &s.dense.Get(*idx).value
}

func (s *SparseSet[TData, TKey]) Delete(id TKey) {
	// TODO optimize me
	idx_ptr := s.sparse.Get(int(id))
	if idx_ptr == nil {
		return
	}

	idx := *idx_ptr
	if idx == NULL_INDEX {
		return
	}

	var lastData = *s.dense.Get(s.dense.Len() - 1)

	s.dense.Set(idx, lastData)
	s.sparse.Set(int(id), NULL_INDEX)
	s.sparse.Set(int(lastData.index), idx)
	s.dense.Pop()
}

func (s *SparseSet[TData, TKey]) isDeleted(id TKey) bool {
	idx := s.sparse.Get(int(id))
	if idx == nil || *idx == NULL_INDEX {
		return true
	}

	return false
}

type Collection[T any] struct {
	buckets []Bucket[T]
	last    *Bucket[T]
	count   int

	initialBucketsCount uint32
	initialBucketSize   uint32
}

type Bucket[T any] struct {
	data []T
}

func (c *Collection[T]) Append(obj T) (int, *T) {
	if c.last == nil {
		c.init(c.initialBucketsCount, c.initialBucketSize)
	}
	if c.last.Cap() < 1 {
		c.extend(c.initialBucketSize)
	}
	curr := c.count
	r := c.last.Add(obj)
	c.count++
	return curr, r
}

func (c *Collection[T]) init(buckets uint32, size uint32) {
	c.buckets = make([]Bucket[T], 1, buckets)
	c.buckets[0].data = make([]T, 0, size)
	c.last = &c.buckets[0]
}

func (c *Collection[T]) extend(size uint32) {
	c.buckets = append(c.buckets, Bucket[T]{data: make([]T, 0, size)})
	c.last = &c.buckets[len(c.buckets)-1]
}

func (c *Collection[T]) Len() (l int) {
	return c.count
}

func (c *Collection[T]) Get(id int) *T {
	for _, b := range c.buckets {
		if id >= len(b.data) {
			id -= len(b.data)
			continue
		}
		return &b.data[id]
	}
	return nil
}

func (c *Collection[T]) Set(id int, val T) *T {
	if id >= c.Len() {
		idx := id - c.Len()
		// TODO optimize so it would skip allocating intermediary buckets
		for i := 0; i < idx; i++ {
			var t T
			c.Append(t)
		}
		_, n := c.Append(val)
		return n
	}
	for _, b := range c.buckets {
		if id >= len(b.data) {
			id -= len(b.data)
			continue
		}
		b.data[id] = val
		return &b.data[id]
	}
	panic("out of bounds")
}

func (c *Collection[T]) Pop() T {
	var value T
	id := len(c.last.data) - 1
	if id < 0 {
		id = len(c.last.data) - 1
	}
	value, c.last.data = c.last.data[id], c.last.data[:id]
	if len(c.last.data) == 0 {
		c.buckets = c.buckets[:len(c.buckets)-1]
		c.last = &c.buckets[len(c.buckets)-1]
	}
	c.count--
	return value
}

func (b *Bucket[T]) Cap() int {
	return cap(b.data) - len(b.data)
}

func (b *Bucket[T]) Add(t T) *T {
	b.data = append(b.data, t)
	return &b.data[len(b.data)-1]
}
