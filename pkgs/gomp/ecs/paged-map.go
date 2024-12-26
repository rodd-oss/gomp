/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "github.com/negrel/assert"

const (
	page_size_shift = 12
	page_size       = 1 << page_size_shift
	paged_map_size  = 1 << 10
)

type Page[K EntityID, V any] map[K]V
type PagedMap[K EntityID, V any] struct {
	pages []Page[K, V]
	len   int
}

func NewPagedMap[K EntityID, V any]() *PagedMap[K, V] {
	return &PagedMap[K, V]{
		pages: make([]Page[K, V], paged_map_size),
	}
}

func (m *PagedMap[K, V]) Get(key K) (value V, ok bool) {
	page_id := m.getPageId(key)
	if page_id >= cap(m.pages) {
		return value, false
	}
	page := m.pages[page_id]
	if page == nil {
		return value, false
	}
	v, ok := page[key]
	return v, ok
}

func (m *PagedMap[K, V]) Set(key K, value V) {
	page_id := m.getPageId(key)
	if page_id >= cap(m.pages) {
		// extend the pages slice
		new_pages := make([]Page[K, V], cap(m.pages)*2)
		copy(new_pages, m.pages)
		m.pages = new_pages
	}
	page := m.pages[page_id]
	if page == nil {
		page = make(Page[K, V], page_size)
		m.pages[page_id] = page
	}
	_, ok := page[key]
	if !ok {
		m.len++
	}
	page[key] = value
}

func (m *PagedMap[K, V]) Delete(key K) {
	page_id := m.getPageId(key)
	// Do not attempt to delete a value that does not exist
	assert.True(page_id < cap(m.pages))
	page := m.pages[page_id]
	// Do not attempt to delete a value that does not exist
	assert.True(page != nil)
	delete(page, key)
	m.len--
}

func (m *PagedMap[K, V]) getPageId(key K) int {
	id := key >> page_size_shift
	return int(id)
}

func (m *PagedMap[K, V]) Len() int {
	return m.len
}
