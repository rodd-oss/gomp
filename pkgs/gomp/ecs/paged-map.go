/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"github.com/negrel/assert"
)

const (
	page_size_shift int32 = 10
	page_size       int32 = 1 << page_size_shift
	book_size       int32 = 1 << 10
)

type MapPage[K EntityID, V any] map[K]V
type PagedMap[K EntityID, V any] struct {
	len  int32
	book []SlicePage[MapValue[V]]
}
type MapValue[V any] struct {
	value V
	ok    bool
}

func NewPagedMap[K EntityID, V any]() *PagedMap[K, V] {
	return &PagedMap[K, V]{
		book: make([]SlicePage[MapValue[V]], book_size),
	}
}

func (m *PagedMap[K, V]) Get(key K) (value V, ok bool) {
	page_id, index := m.getPageIdAndIndex(key)
	if page_id >= cap(m.book) {
		return value, false
	}
	page := m.book[page_id]
	if page.data == nil {
		return value, false
	}
	if index >= cap(page.data) {
		return value, false
	}
	d := page.data[index]
	return d.value, d.ok
}

func (m *PagedMap[K, V]) Set(key K, value V) {
	page_id, index := m.getPageIdAndIndex(key)
	if page_id >= cap(m.book) {
		// extend the pages slice
		new_pages := make([]SlicePage[MapValue[V]], cap(m.book)*2)
		m.book = append(m.book, new_pages...)
		m.Set(key, value)
		return
	}
	page := m.book[page_id]
	if page.data == nil {
		page.data = make([]MapValue[V], page_size)
		m.book[page_id] = page
	}
	d := &page.data[index]
	if !d.ok {
		m.len++
		d.ok = true
	}
	d.value = value
}

func (m *PagedMap[K, V]) Delete(key K) {
	page_id, index := m.getPageIdAndIndex(key)
	// Do not attempt to delete a value that does not exist
	assert.True(page_id < cap(m.book))
	page := &m.book[page_id]
	// Do not attempt to delete a value that does not exist
	assert.True(page != nil)
	page.data[index].ok = false
	m.len--
}

func (m *PagedMap[K, V]) getPageIdAndIndex(key K) (page_id int, index int) {
	page_id = int(key) >> page_size_shift
	index = int(int32(key) % page_size)
	return
}

func (m *PagedMap[K, V]) Len() int32 {
	return m.len
}
