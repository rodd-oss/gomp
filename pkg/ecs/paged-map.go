/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

const (
	pageSizeShift   = 10
	pageSize        = 1 << pageSizeShift
	pageSizeMask    = pageSize - 1
	initialBookSize = 32 // Starting with a small initial book size
)

type PagedMap[K Entity | SharedComponentInstanceId, V any] struct {
	len  int
	book []SlicePage[MapValue[V]]
}

type MapValue[V any] struct {
	value V
	ok    bool
}

func NewPagedMap[K Entity | SharedComponentInstanceId, V any]() PagedMap[K, V] {
	return PagedMap[K, V]{
		book: make([]SlicePage[MapValue[V]], 0, initialBookSize),
	}
}

func (m *PagedMap[K, V]) Get(key K) (value V, ok bool) {
	pageID, index := m.getPageIDAndIndex(key)
	if pageID >= len(m.book) {
		return value, false
	}
	page := &m.book[pageID]
	if page.data == nil {
		return value, false
	}
	d := &page.data[index]
	return d.value, d.ok
}

func (m *PagedMap[K, V]) Set(key K, value V) {
	pageID, index := m.getPageIDAndIndex(key)
	if pageID >= len(m.book) {
		m.expandBook(pageID + 1)
	}
	page := &m.book[pageID]
	if page.data == nil {
		page.data = make([]MapValue[V], pageSize)
	}
	entry := &page.data[index]
	if !entry.ok {
		m.len++
		entry.ok = true
	}
	entry.value = value
}

func (m *PagedMap[K, V]) Delete(key K) {
	pageID, index := m.getPageIDAndIndex(key)
	if pageID >= len(m.book) {
		return
	}
	page := &m.book[pageID]
	if page.data == nil {
		return
	}
	entry := &page.data[index]
	if entry.ok {
		entry.ok = false
		m.len--
	}
}

func (m *PagedMap[K, V]) Has(key K) bool {
	pageID, index := m.getPageIDAndIndex(key)
	if pageID >= len(m.book) {
		return false
	}
	page := &m.book[pageID]
	if page.data == nil {
		return false
	}
	return page.data[index].ok
}

func (m *PagedMap[K, V]) getPageIDAndIndex(key K) (pageID int, index int) {
	return int(key) >> pageSizeShift, int(key) & pageSizeMask
}

func (m *PagedMap[K, V]) expandBook(minLen int) {
	if minLen <= cap(m.book) {
		m.book = m.book[:minLen]
		return
	}
	newCap := minLen
	if newCap < 2*cap(m.book) {
		newCap = 2 * cap(m.book)
	}
	newBook := make([]SlicePage[MapValue[V]], minLen, newCap)
	copy(newBook, m.book)
	m.book = newBook
}

func (m *PagedMap[K, V]) Len() int {
	return m.len
}
