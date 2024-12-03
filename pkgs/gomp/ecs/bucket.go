/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Bucket[T any] struct {
	id   int
	size int
	data []T
}

// [1,2][][][]
func NewBucket[T any](size int) Bucket[T] {
	return Bucket[T]{
		data: make([]T, 0, size),
		size: 0,
	}
}

func (b *Bucket[T]) init(id int, size int) {
	b.data = make([]T, 0, size)
	b.id = id
}

func (b *Bucket[T]) Len() int {
	return b.size
}

func (b *Bucket[T]) CapLeft() int {
	return cap(b.data) - len(b.data)
}

func (b *Bucket[T]) Append(value T) (int, *T) {
	index := b.size
	if index >= len(b.data) {
		b.data = append(b.data, value)
	} else {
		b.data[index] = value
	}

	b.size++
	return index, &b.data[index]
}

func (b *Bucket[T]) Get(index int) (data T) {
	return b.data[index]
}

func (b *Bucket[T]) Exists(index int) bool {
	return index < b.size && index >= 0
}

func (b *Bucket[T]) Set(index int, value T, emptyValue T) *T {
	var delta = index - (len(b.data) - 1)

	if delta > 0 {
		b.data = append(b.data, make([]T, delta)...)
	}

	for i := range delta {
		b.data[index-i] = emptyValue
	}

	b.data[index] = value
	b.size = index + 1

	return &b.data[index]
}

func (b *Bucket[T]) Swap(i, j int) {
	if i == j {
		return
	}

	b.data[i], b.data[j] = b.data[j], b.data[i]
}

func (b *Bucket[T]) SoftReduce() {
	b.size--
}

func (b *Bucket[T]) Clean() {
	b.data = b.data[:b.size]
}
