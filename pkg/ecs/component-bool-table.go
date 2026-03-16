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
	"github.com/negrel/assert"
)

//const (
//	pageSizeShift   = 10
//	pageSize        = 1 << pageSizeShift
//	initialBookSize = 1 // Starting with a small initial book size
//)

func NewComponentBoolTable(maxComponentsLen int) ComponentBoolTable {
	return ComponentBoolTable{
		bools:          make([][]bool, 0, initialBookSize),
		lookup:         make(map[Entity]int, initialBookSize),
		bytesArraySize: maxComponentsLen,
	}
}

type ComponentBoolTable struct {
	bools          [][]bool
	lookup         map[Entity]int
	length         int
	bytesArraySize int
}

func (b *ComponentBoolTable) Create(entity Entity) {
	boolsId, ok := b.lookup[entity]
	if !ok {
		b.extend()
		boolsId = b.length
		b.lookup[entity] = boolsId
		b.length += b.bytesArraySize
	}
}

// Set sets the bit at the given index to 1.
func (b *ComponentBoolTable) Set(entity Entity, componentId ComponentId) {
	boolsId, ok := b.lookup[entity]
	if !ok {
		b.extend()
		boolsId = b.length
		b.lookup[entity] = boolsId
		b.length += b.bytesArraySize
	}
	chunkId := boolsId >> pageSizeShift
	index := boolsId % pageSize
	b.bools[chunkId][index+int(componentId)] = true
}

// Unset clears the bit at the given index (sets it to 0).
func (b *ComponentBoolTable) Unset(entity Entity, componentId ComponentId) {
	boolsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")
	chunkId := boolsId >> pageSizeShift
	index := boolsId % pageSize
	b.bools[chunkId][index+int(componentId)] = false
}

func (b *ComponentBoolTable) Test(entity Entity, componentId ComponentId) bool {
	boolsId, ok := b.lookup[entity]
	if !ok {
		return false
	}
	chunkId := boolsId >> pageSizeShift
	index := boolsId % pageSize
	return b.bools[chunkId][index+int(componentId)]
}

func (b *ComponentBoolTable) extend() {
	lastChunkId := b.length >> pageSizeShift
	if lastChunkId == len(b.bools) && b.length%pageSize == 0 {
		b.bools = append(b.bools, make([]bool, b.bytesArraySize*pageSize))
	}
}

func (b *ComponentBoolTable) AllSet(entity Entity, yield func(ComponentId) bool) {
	boolsId, ok := b.lookup[entity]
	if !ok {
		return
	}
	chunkId := boolsId >> pageSizeShift
	index := boolsId % pageSize
	for i := 0; i < b.bytesArraySize; i++ {
		if b.bools[chunkId][index+i] {
			if !yield(ComponentId(i)) {
				return
			}
		}
	}
}
