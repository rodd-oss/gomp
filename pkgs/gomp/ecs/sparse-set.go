/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

func NewSparseSet[TData any, TKey EntityID | ComponentID | ECSID | int](buckets uint32, bucketSize uint32) SparseSet[TData, TKey] {
	set := SparseSet[TData, TKey]{
		deleted: make(map[TKey]int),
	}
	set.sparse.initialBucketsCount = buckets
	set.dense.initialBucketsCount = buckets
	set.sparse.initialBucketSize = bucketSize
	set.dense.initialBucketSize = bucketSize

	return set
}

type SparseSet[TData any, TKey EntityID | ComponentID | ECSID | int] struct {
	// TODO: refactor map to a slice with using of a deletedSparseElements slice
	sparse  Collection[int]
	dense   Collection[TData]
	deleted map[TKey]int
}

func (s *SparseSet[TData, TKey]) Set(id TKey, data TData) *TData {
	if pos, ok := s.deleted[id]; ok {
		delete(s.deleted, id)
		return s.dense.Set(pos, data)
	}

	pos := s.sparse.Get(int(id))
	if pos != nil {
		return s.dense.Set(*pos, data)
	}
	idx, r := s.dense.Append(data)
	s.sparse.Set(int(id), idx)
	return r
}

func (s *SparseSet[TData, TKey]) Get(id TKey) *TData {
	if s.isDeleted(id) {
		return nil
	}

	idx := s.sparse.Get(int(id))
	if idx == nil {
		return nil
	}
	return s.dense.Get(*idx)
}

func (s *SparseSet[TData, TKey]) Delete(id TKey) {
	// TODO optimize me
	idx := s.sparse.Get(int(id))
	if idx == nil {
		return
	}
	s.deleted[id] = *idx
}

func (s *SparseSet[TData, TKey]) isDeleted(id TKey) bool {
	_, deleted := s.deleted[id]
	return deleted
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
		idx := c.Len() - id
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

func (b *Bucket[T]) Cap() int {
	return cap(b.data) - len(b.data)
}

func (b *Bucket[T]) Add(t T) *T {
	b.data = append(b.data, t)
	return &b.data[len(b.data)-1]
}
