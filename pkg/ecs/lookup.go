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

import "sync"

const (
	IndexBits = 10
	IndexMask = 1<<IndexBits - 1
)

type LookupMap[T any] struct {
	data map[int]map[int]map[int]T
	mx   sync.RWMutex
}

func NewLookupMap[T any]() *LookupMap[T] {
	return &LookupMap[T]{
		data: make(map[int]map[int]map[int]T),
	}
}

func (m *LookupMap[T]) Get(key uint32) (T, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	index0, index1, index2 := m.calcIndex(key)
	var ret T

	v0, ok := m.data[index0]
	if !ok {
		return ret, false
	}
	v1, ok := v0[index1]
	if !ok {
		return ret, false
	}
	ret, ok = v1[index2]
	return ret, ok
}

func (m *LookupMap[T]) Set(key uint32, value T) {
	m.mx.Lock()
	defer m.mx.Unlock()

	index0, index1, index2 := m.calcIndex(key)

	v0, ok := m.data[index0]
	if !ok {
		v0 = make(map[int]map[int]T)
		m.data[index0] = v0
	}
	v1, ok := v0[index1]
	if !ok {
		v1 = make(map[int]T)
		v0[index1] = v1
	}
	v1[index2] = value
}

func (m *LookupMap[T]) calcIndex(key uint32) (int, int, int) {
	return int(key >> 20 & IndexMask), int(key >> 10 & IndexMask), int(key & IndexMask)
}

const (
	k  = 1024*1024*1024*4 - 1
	i0 = k >> 30 & IndexMask
	i1 = k >> 20 & IndexMask
	i2 = k >> 10 & IndexMask
	i3 = k & IndexMask
)
