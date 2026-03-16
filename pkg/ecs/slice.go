/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Саратов Рай Donated 1 000 RUB

Thank you for your support!
*/

package ecs

import (
	"sync"

	"github.com/negrel/assert"
)

// Slice tricks https://ueokande.github.io/go-slice-tricks/

type Slice[T any] struct {
	data []T
}

func NewSlice[T any](size int) (a Slice[T]) {
	a.data = make([]T, 0, size)
	return a
}

func (a *Slice[T]) Len() int {
	return len(a.data)
}

func (a *Slice[T]) Get(index int) *T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < len(a.data), "index out of range")

	return &a.data[index]
}

func (a *Slice[T]) GetValue(index int) T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < len(a.data), "index out of range")

	return a.data[index]
}

func (a *Slice[T]) Set(index int, value T) *T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < len(a.data), "index out of range")

	a.data[index] = value

	return &a.data[index]
}

func (a *Slice[T]) Append(values ...T) []T {
	a.data = append(a.data, values...)
	disCap := 1 << (FastIntLog2(uint64(len(a.data))) + 1)
	if cap(a.data) > disCap {
		a.data = a.data[:len(a.data):disCap]
	}
	return a.data
}

func (a *Slice[T]) SoftReduce() {
	assert.True(len(a.data) > 0, "Len is already 0")
	a.data = a.data[:len(a.data)-1]
}

func (a *Slice[T]) Reset() {
	a.data = a.data[:0]
}

func (a *Slice[T]) Copy(fromIndex, toIndex int) {
	assert.True(fromIndex >= 0, "index out of range")
	assert.True(fromIndex < len(a.data), "index out of range")
	from := a.Get(fromIndex)

	assert.True(toIndex >= 0, "index out of range")
	assert.True(toIndex < len(a.data), "index out of range")
	to := a.Get(toIndex)

	*to = *from
}

func (a *Slice[T]) Swap(i, j int) (newI, NewJ *T) {
	assert.True(i >= 0, "index out of range")
	assert.True(i < len(a.data), "index out of range")
	x := a.Get(i)

	assert.True(j >= 0, "index out of range")
	assert.True(j < len(a.data), "index out of range")
	y := a.Get(j)

	*x, *y = *y, *x
	return x, y
}

func (a *Slice[T]) Last() *T {
	index := len(a.data) - 1
	assert.True(index >= 0, "index out of range")

	return a.Get(index)
}

func (a *Slice[T]) Raw() []T {
	return a.data
}

func (a *Slice[T]) getPageIdAndIndex(index int) (int, int) {
	pageId := index >> pageSizeShift
	assert.True(pageId < len(a.data), "index out of range")

	index %= pageSize
	assert.True(index < pageSize, "index out of range")

	return pageId, index
}

func (a *Slice[T]) Each() func(yield func(int, *T) bool) {
	return func(yield func(int, *T) bool) {
		for j := len(a.data) - 1; j >= 0; j-- {
			if !yield(j, &a.data[j]) {
				return
			}
		}
	}
}

func (a *Slice[T]) EachData() func(yield func(*T) bool) {
	return func(yield func(*T) bool) {
		for j := len(a.data) - 1; j >= 0; j-- {
			if !yield(&a.data[j]) {
				return
			}
		}
	}
}

func (a *Slice[T]) EachDataValue() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for j := len(a.data) - 1; j >= 0; j-- {
			if !yield(a.data[j]) {
				return
			}
		}
	}
}

func (a *Slice[T]) EachParallel(numWorkers int) func(yield func(int, *T, int) bool) {
	return func(yield func(int, *T, int) bool) {
		assert.True(numWorkers > 0)
		var chunkSize = len(a.data) / numWorkers
		var wg sync.WaitGroup

		wg.Add(numWorkers)
		for workedId := 0; workedId < numWorkers; workedId++ {
			startIndex := workedId * chunkSize
			endIndex := startIndex + chunkSize - 1
			if workedId == numWorkers-1 { // have to set endIndex to entities length, if last worker
				endIndex = len(a.data)
			}
			go func(start int, end int) {
				defer wg.Done()
				for i := range a.data[start:end] {
					if !yield(i, &a.data[i+startIndex], workedId) {
						return
					}
				}
			}(startIndex, endIndex)
		}
		wg.Wait()
	}
}

func (a *Slice[T]) EachDataValueParallel(numWorkers int) func(yield func(T, int) bool) {
	return func(yield func(T, int) bool) {
		assert.True(numWorkers > 0)
		var chunkSize = len(a.data) / numWorkers
		var wg sync.WaitGroup

		wg.Add(numWorkers)
		for workedId := 0; workedId < numWorkers; workedId++ {
			startIndex := workedId * chunkSize
			endIndex := startIndex + chunkSize - 1
			if workedId == numWorkers-1 { // have to set endIndex to entities length, if last worker
				endIndex = len(a.data)
			}
			go func(start int, end int) {
				defer wg.Done()
				for i := range a.data[start:end] {
					if !yield(a.data[i+startIndex], workedId) {
						return
					}
				}
			}(startIndex, endIndex)
		}
		wg.Wait()
	}
}

func (a *Slice[T]) EachDataParallel(numWorkers int) func(yield func(*T, int) bool) {
	return func(yield func(*T, int) bool) {
		assert.True(numWorkers > 0)
		var chunkSize = len(a.data) / numWorkers
		var wg sync.WaitGroup

		wg.Add(numWorkers)
		for workedId := 0; workedId < numWorkers; workedId++ {
			startIndex := workedId * chunkSize
			endIndex := startIndex + chunkSize - 1
			if workedId == numWorkers-1 { // have to set endIndex to entities length, if last worker
				endIndex = len(a.data)
			}
			go func(start int, end int) {
				defer wg.Done()
				for i := range a.data[start:end] {
					if !yield(&a.data[i+startIndex], workedId) {
						return
					}
				}
			}(startIndex, endIndex)
		}
		wg.Wait()
	}
}
