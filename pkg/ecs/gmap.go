/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package ecs

import (
	"iter"

	"github.com/negrel/assert"
)

func NewGenMap[K comparable, V any](cap int) GenMap[K, V] {
	return GenMap[K, V]{
		data:       make(map[K]genMapValue[V], cap),
		generation: 1,
	}
}

type genMapValue[V any] struct {
	generation int32
	value      V
}

type GenMap[K comparable, V any] struct {
	generation int32
	data       map[K]genMapValue[V]
	len        int
}

func (m *GenMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.data[key]
	if !ok {
		return v.value, ok
	}
	if v.generation != m.generation {
		return v.value, false
	}
	return v.value, ok
}

func (m *GenMap[K, V]) Set(key K, value V) {
	if !m.Has(key) {
		m.len++
	}
	m.data[key] = genMapValue[V]{
		generation: m.generation,
		value:      value,
	}
}

func (m *GenMap[K, V]) Reset() {
	m.generation++
	m.len = 0
}

func (m *GenMap[K, V]) Delete(key K) {
	_, ok := m.data[key]
	assert.True(ok)
	m.data[key] = genMapValue[V]{
		generation: m.generation - 1,
	}
}

func (m *GenMap[K, V]) Has(key K) bool {
	v, ok := m.data[key]
	return ok && v.generation == m.generation
}

func (m *GenMap[K, V]) Len() int {
	return m.len
}

func (m *GenMap[K, V]) Each() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m.data {
			if v.generation != m.generation {
				continue
			}
			if !yield(k, v.value) {
				return
			}
		}
	}
}

func (m *GenMap[K, V]) Clear() {
	for i, v := range m.data {
		if v.generation == m.generation {
			continue
		}
		delete(m.data, i)
	}
}
