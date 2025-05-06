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
	"math/bits"
)

const (
	uintShift = 7 - 64/bits.UintSize
)

func NewComponentBitTable(maxComponentsLen int) ComponentBitTable {
	bitsetSize := ((maxComponentsLen - 1) / bits.UintSize) + 1
	return ComponentBitTable{
		bitsetsBook:  make([][]uint, 0, initialBookSize),
		entitiesBook: make([]*entityArray, 0, initialBookSize),
		lookup:       NewPagedMap[Entity, int](),
		bitsetSize:   bitsetSize,
		pageSize:     bitsetSize * pageSize,
	}
}

type ComponentBitTable struct {
	bitsetsBook  [][]uint
	entitiesBook []*entityArray
	lookup       PagedMap[Entity, int]
	length       int
	bitsetSize   int
	pageSize     int
}

type entityArray [pageSize]Entity

func (b *ComponentBitTable) Create(entity Entity) {
	assert.False(b.lookup.Has(entity), "entity already exists")

	b.extend()
	bitsId := b.length
	b.lookup.Set(entity, bitsId)
	pageId, entityId := b.getPageIDAndEntityIndex(bitsId)
	b.entitiesBook[pageId][entityId] = entity
	b.length++
}

func (b *ComponentBitTable) Delete(entity Entity) {
	bitsetIndex, ok := b.lookup.Get(entity)
	assert.True(ok, "entity not found")

	// Get the index of the last entity
	lastIndex := b.length - 1

	// If this is not the last entity, swap with the last one
	if bitsetIndex != lastIndex {
		lastPageId, lastEntityId := b.getPageIDAndEntityIndex(lastIndex)
		lastBitsetId := lastEntityId * b.bitsetSize
		deletePageId, deleteEntityId := b.getPageIDAndEntityIndex(bitsetIndex)
		deleteBitsetId := deleteEntityId * b.bitsetSize

		// Copy bitset from last entity to the deleted entity's position
		for i := 0; i < b.bitsetSize; i++ {
			b.bitsetsBook[deletePageId][deleteBitsetId+i] = b.bitsetsBook[lastPageId][lastBitsetId+i]
			b.bitsetsBook[lastPageId][lastBitsetId+i] = 0
		}

		// Get the last entity and update its position in lookup
		lastEntity := b.entitiesBook[lastPageId][lastEntityId]
		b.entitiesBook[deletePageId][deleteEntityId] = lastEntity
		b.lookup.Set(lastEntity, bitsetIndex)
	}

	b.lookup.Delete(entity)
	b.length--
}

// Set sets the bit at the given index to 1.
func (b *ComponentBitTable) Set(entity Entity, componentId ComponentId) {
	bitsId, ok := b.lookup.Get(entity)
	assert.True(ok, "entity not found")

	pageId, bitsetId := b.getPageIDAndBitsetIndex(bitsId)
	offset := int(componentId) >> uintShift
	b.bitsetsBook[pageId][bitsetId+offset] |= 1 << (componentId % bits.UintSize)
}

// Unset clears the bit at the given index (sets it to 0).
func (b *ComponentBitTable) Unset(entity Entity, componentId ComponentId) {
	bitsId, ok := b.lookup.Get(entity)
	assert.True(ok, "entity not found")

	pageId, bitsetId := b.getPageIDAndBitsetIndex(bitsId)
	offset := int(componentId) >> uintShift
	b.bitsetsBook[pageId][bitsetId+offset] &= ^(1 << (componentId % bits.UintSize))
}

func (b *ComponentBitTable) Test(entity Entity, componentId ComponentId) bool {
	bitsId, ok := b.lookup.Get(entity)
	if !ok {
		return false
	}
	pageId, bitsetId := b.getPageIDAndBitsetIndex(bitsId)
	offset := int(componentId) >> uintShift
	return (b.bitsetsBook[pageId][bitsetId+offset] & (1 << (componentId % bits.UintSize))) != 0
}

func (b *ComponentBitTable) AllSet(entity Entity, yield func(ComponentId) bool) {
	bitsId, ok := b.lookup.Get(entity)
	if !ok {
		return
	}
	pageId, bitsetId := b.getPageIDAndBitsetIndex(bitsId)
	bitset := b.bitsetsBook[pageId][bitsetId : bitsetId+b.bitsetSize]
	for i, set := range bitset {
		j := 0
		for set != 0 {
			if set&1 == 1 {
				if !yield(ComponentId(i*bits.UintSize + j)) {
					return
				}
			}
			set >>= 1
			j++
		}
	}
}

func (b *ComponentBitTable) extend() {
	lastChunkId, lastEntityId := b.getPageIDAndEntityIndex(b.length)
	if lastChunkId == len(b.bitsetsBook) && lastEntityId == 0 {
		b.bitsetsBook = append(b.bitsetsBook, make([]uint, b.pageSize))
		b.entitiesBook = append(b.entitiesBook, &entityArray{})
	}
}

func (b *ComponentBitTable) getPageIDAndBitsetIndex(index int) (int, int) {
	return index >> pageSizeShift, (index & pageSizeMask) * b.bitsetSize
}

func (b *ComponentBitTable) getPageIDAndEntityIndex(index int) (int, int) {
	return index >> pageSizeShift, index & pageSizeMask
}
