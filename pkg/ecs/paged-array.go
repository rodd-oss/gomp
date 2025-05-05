/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"github.com/negrel/assert"
	"gomp/pkg/worker"
	"sync"
)

/*
Go Slice tips and tricks - https://ueokande.github.io/go-slice-tricks/
*/

func NewPagedArray[T any]() (a PagedArray[T]) {
	a.book = make([]*ArrayPage[T], initialBookSize)
	for i := 0; i < initialBookSize; i++ {
		a.book[i] = &ArrayPage[T]{}
	}
	a.edTasks = make([]EachDataTask[T], initialBookSize)
	a.edvTasks = make([]EachDataValueTask[T], initialBookSize)
	a.edvpTasks = make([]EachDataValueParallelTask[T], pageSize)
	a.edpTasks = make([]EachDataParallelTask[T], pageSize)
	return a
}

type SlicePage[T any] struct {
	len  int
	data []T
}

type PagedArray[T any] struct {
	book             []*ArrayPage[T]
	currentPageIndex int
	len              int
	wg               sync.WaitGroup

	// Cache
	edvTasks  []EachDataValueTask[T]
	edTasks   []EachDataTask[T]
	edvpTasks []EachDataValueParallelTask[T]
	edpTasks  []EachDataParallelTask[T]
}

type ArrayPage[T any] struct {
	len  int
	data [pageSize]T
}

func (a *PagedArray[T]) Len() int {
	return a.len
}

func (a *PagedArray[T]) Get(index int) *T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < a.len, "index out of range")

	pageId, index := a.getPageIdAndIndex(index)
	page := a.book[pageId]

	return &(page.data[index])
}

func (a *PagedArray[T]) GetValue(index int) T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < a.len, "index out of range")

	pageId, index := a.getPageIdAndIndex(index)
	page := a.book[pageId]

	return page.data[index]
}

func (a *PagedArray[T]) Set(index int, value T) *T {
	assert.True(index >= 0, "index out of range")
	assert.True(index < a.len, "index out of range")

	pageId, index := a.getPageIdAndIndex(index)
	page := a.book[pageId]

	page.data[index] = value

	return &page.data[index]
}

func (a *PagedArray[T]) extend() {
	oldLen := len(a.book)
	newLen := oldLen * 2
	newBooks := make([]*ArrayPage[T], newLen)
	copy(newBooks, a.book)
	for i := oldLen; i < newLen; i++ {
		newBooks[i] = &ArrayPage[T]{}
	}
	a.book = newBooks

	newEdvTasks := make([]EachDataValueTask[T], newLen)
	copy(newEdvTasks, a.edvTasks)
	a.edvTasks = newEdvTasks

	newEdTasks := make([]EachDataTask[T], newLen)
	copy(newEdTasks, a.edTasks)
	a.edTasks = newEdTasks

	newEdvpTasks := make([]EachDataValueParallelTask[T], pageSize)
	copy(newEdvpTasks, a.edvpTasks)
	a.edvpTasks = newEdvpTasks

	newEdpTasks := make([]EachDataParallelTask[T], pageSize)
	copy(newEdpTasks, a.edpTasks)
	a.edpTasks = newEdpTasks
}

func (a *PagedArray[T]) Append(value T) *T {
	var result *T
	if a.currentPageIndex >= len(a.book) {
		a.extend()
	}

	page := a.book[a.currentPageIndex]

	if page.len == pageSize {
		a.currentPageIndex++
		if a.currentPageIndex >= len(a.book) {
			a.extend()
		}
		page = a.book[a.currentPageIndex]
	}
	page.data[page.len] = value
	result = &page.data[page.len]
	page.len++
	a.len++
	return result
}

func (a *PagedArray[T]) AppendMany(values ...T) *T {
	var result *T
	for i := range values {
		value := values[i]
		if a.currentPageIndex >= len(a.book) {
			a.extend()
		}

		page := a.book[a.currentPageIndex]

		if page.len == pageSize {
			a.currentPageIndex++
			if a.currentPageIndex >= len(a.book) {
				a.extend()
			}
			page = a.book[a.currentPageIndex]
		}
		page.data[page.len] = value
		result = &page.data[page.len]
		page.len++
		a.len++
	}

	return result
}

// SoftReduce - removes the last element of the array
func (a *PagedArray[T]) SoftReduce() {
	assert.True(a.len > 0, "Len is already 0")

	page := a.book[a.currentPageIndex]
	assert.True(page.len > 0, "Len is already 0")

	page.len--
	a.len--

	if page.len == 0 && a.currentPageIndex != 0 {
		a.currentPageIndex--
	}
}

// Reset - resets the array to its initial state
func (a *PagedArray[T]) Reset() {
	for i := 0; i <= a.currentPageIndex; i++ {
		page := a.book[i]
		page.len = 0
	}

	a.currentPageIndex = 0
	a.len = 0
}

func (a *PagedArray[T]) Copy(fromIndex, toIndex int) {
	assert.True(fromIndex >= 0, "index out of range")
	assert.True(fromIndex < a.len, "index out of range")
	from := a.Get(fromIndex)

	assert.True(toIndex >= 0, "index out of range")
	assert.True(toIndex < a.len, "index out of range")
	to := a.Get(toIndex)

	*to = *from
}

func (a *PagedArray[T]) Swap(i, j int) (newI, NewJ *T) {
	assert.True(i >= 0, "index out of range")
	assert.True(i < a.len, "index out of range")
	x := a.Get(i)

	assert.True(j >= 0, "index out of range")
	assert.True(j < a.len, "index out of range")
	y := a.Get(j)

	*x, *y = *y, *x
	return x, y
}

func (a *PagedArray[T]) Last() *T {
	index := a.len - 1
	assert.True(index >= 0, "index out of range")

	return a.Get(index)
}

func (a *PagedArray[T]) Raw(result []T) []T {
	totalLen := a.len
	// Ensure the result has enough capacity and set the correct length.
	if cap(result) < totalLen {
		result = make([]T, totalLen)
	} else {
		result = result[:totalLen]
	}

	pos := 0
	for i := 0; i <= a.currentPageIndex; i++ {
		page := a.book[i]
		n := copy(result[pos:], page.data[:page.len])
		pos += n
	}

	return result
}

func (a *PagedArray[T]) getPageIdAndIndex(index int) (int, int) {
	return index >> pageSizeShift, index & pageSizeMask
}

// =========================
// Iterators
// =========================

func (a *PagedArray[T]) Each() func(yield func(int, *T) bool) {
	return func(yield func(int, *T) bool) {
		var page *ArrayPage[T]
		var indexOffset int

		book := a.book

		if a.len == 0 {
			return
		}

		for i := a.currentPageIndex; i >= 0; i-- {
			page = book[i]
			indexOffset = i << pageSizeShift

			for j := page.len - 1; j >= 0; j-- {
				if !yield(indexOffset+j, &page.data[j]) {
					return
				}
			}
		}
	}
}

func (a *PagedArray[T]) EachParallel(numWorkers int) func(yield func(int, *T, int) bool) {
	return func(yield func(int, *T, int) bool) {
		assert.True(numWorkers > 0)
		var chunkSize = a.len / numWorkers
		var wg sync.WaitGroup

		wg.Add(numWorkers)
		for workedId := 0; workedId < numWorkers; workedId++ {
			startIndex := workedId * chunkSize
			endIndex := startIndex + chunkSize - 1
			if workedId == numWorkers-1 { // have to set endIndex to entities length, if last worker
				endIndex = a.len
			}
			go func(start int, end int) {
				defer wg.Done()
				r := end - start
				for i := range r {
					if !yield(i, a.Get(i+startIndex), workedId) {
						return
					}
				}
			}(startIndex, endIndex)
		}
		wg.Wait()
	}
}

func (a *PagedArray[T]) EachData() func(yield func(*T) bool) {
	return func(yield func(*T) bool) {
		var page *ArrayPage[T]
		var book = a.book

		if a.len == 0 {
			return
		}

		for i := a.currentPageIndex; i >= 0; i-- {
			page = book[i]

			for j := page.len - 1; j >= 0; j-- {
				if !yield(&page.data[j]) {
					return
				}
			}
		}
	}
}

func (a *PagedArray[T]) EachDataValue() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		var page *ArrayPage[T]
		var book = a.book

		if a.len == 0 {
			return
		}

		for i := a.currentPageIndex; i >= 0; i-- {
			page = book[i]

			for j := page.len - 1; j >= 0; j-- {
				if !yield(page.data[j]) {
					return
				}
			}
		}
	}
}

func (a *PagedArray[T]) ProcessDataValue(handler func(T, worker.WorkerId), pool *worker.Pool) {
	assert.NotNil(handler)
	assert.NotNil(pool)
	pool.GroupAdd(a.currentPageIndex + 1)
	for i := a.currentPageIndex; i >= 0; i-- {
		j := a.currentPageIndex - i
		a.edvTasks[j].page = a.book[i]
		a.edvTasks[j].f = handler
		pool.ProcessGroupTask(&a.edvTasks[j])
	}
	pool.GroupWait()
}

func (a *PagedArray[T]) ProcessData(handler func(*T, worker.WorkerId), pool *worker.Pool) {
	assert.NotNil(handler)
	assert.NotNil(pool)
	pool.GroupAdd(a.currentPageIndex + 1)
	for i := a.currentPageIndex; i >= 0; i-- {
		j := a.currentPageIndex - i
		a.edTasks[j].page = a.book[i]
		a.edTasks[j].f = handler
		pool.ProcessGroupTask(&a.edTasks[j])
	}
	pool.GroupWait()
}

func (a *PagedArray[T]) EachDataValueParallel(handler func(T, worker.WorkerId), pool *worker.Pool) {
	assert.NotNil(handler)
	assert.NotNil(pool)

	var page *ArrayPage[T]
	var book = a.book

	pool.GroupAdd(a.len)
	for i := a.currentPageIndex; i >= 0; i-- {
		page = book[i]
		n := (a.currentPageIndex - i) * pageSize
		for j := page.len - 1; j >= 0; j-- {
			m := n + (page.len - 1 - j)
			a.edvpTasks[m].page = page
			a.edvpTasks[m].index = m
			a.edvpTasks[m].f = handler
			pool.ProcessGroupTask(&a.edvpTasks[m])
		}
	}
	pool.GroupWait()
}

func (a *PagedArray[T]) EachDataParallel(handler func(*T, worker.WorkerId), pool *worker.Pool) {
	assert.NotNil(handler)
	assert.NotNil(pool)

	var page *ArrayPage[T]
	var book = a.book

	pool.GroupAdd(a.len)
	for i := a.currentPageIndex; i >= 0; i-- {
		page = book[i]
		n := (a.currentPageIndex - i) * pageSize
		for j := page.len - 1; j >= 0; j-- {
			m := n + (page.len - 1 - j)
			a.edpTasks[m].page = page
			a.edpTasks[m].index = m
			a.edpTasks[m].f = handler
			pool.ProcessGroupTask(&a.edpTasks[m])
		}
	}
	pool.GroupWait()
}

// =========================
// TASKS
// =========================

type EachDataValueTask[T any] struct {
	f    func(T, worker.WorkerId)
	page *ArrayPage[T]
}

func (t *EachDataValueTask[T]) Run(workerId worker.WorkerId) error {
	for i := 0; i < t.page.len; i++ {
		t.f(t.page.data[i], workerId)
	}
	return nil
}

type EachDataTask[T any] struct {
	f    func(*T, worker.WorkerId)
	page *ArrayPage[T]
}

func (t *EachDataTask[T]) Run(workerId worker.WorkerId) error {
	for i := 0; i < t.page.len; i++ {
		t.f(&t.page.data[i], workerId)
	}
	return nil
}

type EachDataValueParallelTask[T any] struct {
	f     func(T, worker.WorkerId)
	page  *ArrayPage[T]
	index int
}

func (t *EachDataValueParallelTask[T]) Run(workerId worker.WorkerId) error {
	t.f(t.page.data[t.index], workerId)
	return nil
}

type EachDataParallelTask[T any] struct {
	f     func(*T, worker.WorkerId)
	page  *ArrayPage[T]
	index int
}

func (t *EachDataParallelTask[T]) Run(workerId worker.WorkerId) error {
	t.f(&t.page.data[t.index], workerId)
	return nil
}
