/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Map[K Entity | SharedComponentInstanceId, V any] struct {
	book map[K]V
}

func NewMap[K Entity | SharedComponentInstanceId, V any]() Map[K, V] {
	return Map[K, V]{
		book: make(map[K]V, initialBookSize),
	}
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	value, ok = m.book[key]
	return value, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.book[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.book, key)
}

func (m *Map[K, V]) Len() int {
	return len(m.book)
}
