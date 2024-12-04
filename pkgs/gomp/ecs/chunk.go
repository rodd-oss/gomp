/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "iter"

type ChunkArray[T any] struct {
	buffer               []ChunkArrayElement[T]
	first                *ChunkArrayElement[T]
	current              *ChunkArrayElement[T]
	last                 *ChunkArrayElement[T]
	size                 int
	initialChunkCapacity int
	bufferInitialCap     int
	bufferSizeIndex      int
}

func NewChunkArray[T any](bufferCapacity int, chunkCapacity int) (arr *ChunkArray[T]) {
	arr = new(ChunkArray[T])
	arr.bufferInitialCap = bufferCapacity
	arr.initialChunkCapacity = chunkCapacity
	arr.buffer = make([]ChunkArrayElement[T], bufferCapacity)
	arr.bufferSizeIndex = 0

	chunk := arr.makeChunk(chunkCapacity)
	chunk.parent = arr

	arr.first = chunk
	arr.current = chunk
	arr.last = chunk

	return arr
}

func (a *ChunkArray[T]) Len() int {
	return a.size
}

func (a *ChunkArray[T]) Get(index int) (T, bool) {
	return a.first.Get(index)
}

func (a *ChunkArray[T]) Set(index int, value T) bool {
	return a.first.Set(index, value)
}

func (a *ChunkArray[T]) Append(value T) int {
	a.size++
	return a.current.Append(value)
}

func (a *ChunkArray[T]) SoftReduce() {
	a.size--
	a.current.SoftReduce()
}

func (a *ChunkArray[T]) Clean() {
	a.last.Clean()
}

func (a *ChunkArray[T]) Swap(i, j int) {
	x, _ := a.Get(i)
	y, _ := a.Get(j)

	a.Set(j, x)
	a.Set(i, y)
}

func (a *ChunkArray[T]) extendBuffer() {
	a.bufferInitialCap += a.bufferInitialCap
	a.buffer = append(a.buffer, make([]ChunkArrayElement[T], a.bufferInitialCap)...)
}

func (a *ChunkArray[T]) makeChunk(cap int) *ChunkArrayElement[T] {
	if a.bufferSizeIndex >= len(a.buffer) {
		a.extendBuffer()
	}

	chunk := &a.buffer[a.bufferSizeIndex]
	chunk.parent = a
	chunk.data = make([]T, 0, cap)
	a.bufferSizeIndex++

	a.current = chunk
	a.last = chunk

	return chunk
}

func (a *ChunkArray[T]) Iter() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		for i := range a.buffer {
			var chunk = &a.buffer[i]
			var offset = i * a.bufferInitialCap

			for j := range chunk.data {
				if !yield(offset+j, &chunk.data[j]) {
					return
				}
			}
		}
	}
}

// ======
// ======
// ======

type ChunkArrayElement[T any] struct {
	next   *ChunkArrayElement[T]
	prev   *ChunkArrayElement[T]
	parent *ChunkArray[T]
	size   int
	data   []T
}

func (c *ChunkArrayElement[T]) Get(index int) (data T, ok bool) {
	if index >= c.size {
		if c.next == nil {
			return data, false
		}

		return c.next.Get(index - c.size)
	}

	data = c.data[index]
	return data, true
}

func (c *ChunkArrayElement[T]) Set(index int, value T) (ok bool) {
	if index >= c.size {
		if c.next == nil {
			return false
		}

		return c.next.Set(index-c.size, value)
	}

	c.data[index] = value
	return true
}

func (c *ChunkArrayElement[T]) Append(value T) int {
	var index = c.size

	if index < len(c.data) {
		c.data[index] = value
		c.size++
		return index
	}

	if index < cap(c.data) {
		c.data = append(c.data, value)
		c.size++
		return index
	}

	if c.next == nil {
		parent := c.parent
		parent.initialChunkCapacity = parent.initialChunkCapacity * 2
		chunk := c.parent.makeChunk(parent.initialChunkCapacity)
		chunk.prev = c
		c.next = chunk
	}

	return c.next.Append(value)
}

func (c *ChunkArrayElement[T]) SoftReduce() {
	if c.size > 0 {
		c.size--
		return
	}

	if c.prev == nil {
		return
	}

	c.parent.current = c.prev
	c.prev.SoftReduce()
}

func (c *ChunkArrayElement[T]) Clean() {
	c.data = c.data[:c.size]

	if len(c.data) == 0 {
		if c.next != nil {
			c.parent.last = c
			c.next = nil
		}

		if c.prev != nil {
			c.prev.Clean()
		}
	}
}
