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
	"math/big"
	"math/bits"
)

const ComponentBitsetPreallocate = 1024

func NewComponentBitSet() ComponentBitSet {
	return ComponentBitSet{
		bits:     make([]BitSet, 0, ComponentBitsetPreallocate),
		entities: make([]Entity, 0, ComponentBitsetPreallocate),
		lookup:   make(map[Entity]int, ComponentBitsetPreallocate),
	}
}

const (
	a = 1243123123 & 0
)

type BitSet = big.Int

type ComponentBitSet struct {
	bits     []BitSet
	entities []Entity
	lookup   map[Entity]int
}

func (b *ComponentBitSet) Get(entity Entity) BitSet {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")
	return b.bits[bitsId]
}

// Set sets the bit at the given index to 1.
func (b *ComponentBitSet) Set(entity Entity, componentId ComponentId) {
	bitsId, ok := b.lookup[entity]
	if !ok {
		bitsId = len(b.bits)
		b.lookup[entity] = bitsId
		b.entities = append(b.entities, entity)
		b.bits = append(b.bits, BitSet{})
	}

	bitSet := &b.bits[bitsId]
	bitSet.SetBit(bitSet, int(componentId), 1)
}

// Unset clears the bit at the given index (sets it to 0).
func (b *ComponentBitSet) Unset(entity Entity, componentId ComponentId) {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")
	bitSet := &b.bits[bitsId]
	bitSet.SetBit(bitSet, int(componentId), 0)
}

// Toggle toggles the bit at the given index.
func (b *ComponentBitSet) Toggle(entity Entity, componentId ComponentId) {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")
	bitSet := &b.bits[bitsId]
	bitSet.SetBit(bitSet, int(componentId), 1-bitSet.Bit(int(componentId)))
}

// IsSet checks if the bit at the given index is set (1). Automatically resizes if the index is out of bounds.
func (b *ComponentBitSet) IsSet(entity Entity, componentId ComponentId) bool {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")
	bitSet := &b.bits[bitsId]

	return bitSet.Bit(int(componentId)) == 1
}

func (b *ComponentBitSet) Delete(entity Entity) {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")

	lastIndex := len(b.bits) - 1
	if bitsId < lastIndex {
		// swap the dead element with the last one
		b.bits[bitsId], b.bits[lastIndex] = b.bits[lastIndex], b.bits[bitsId]
		b.entities[bitsId] = b.entities[lastIndex]

		// update lookup table
		b.lookup[b.entities[bitsId]] = bitsId
	}

	b.bits = b.bits[:lastIndex]
	b.entities = b.entities[:lastIndex]
	delete(b.lookup, entity)
}

func (b *ComponentBitSet) FilterByMask(mask BitSet, yield func(entity Entity) bool) {
	bitsLen := len(b.bits)
	var cmpBitset BitSet
	var zeroBitSet BitSet

	for i := range bitsLen {
		bitSet := &b.bits[i]
		if cmpBitset.And(&cmpBitset, &zeroBitSet).And(bitSet, &mask).Cmp(&mask) != 0 {
			continue
		}

		if !yield(b.entities[i]) {
			return
		}
	}
}

func (b *ComponentBitSet) AllSet(entity Entity, yield func(ComponentId) bool) {
	bitsId, ok := b.lookup[entity]
	assert.True(ok, "entity not found")

	bitSet := b.bits[bitsId].Bits()
	var id ComponentId

	for i := range bitSet {
		v := uint(bitSet[i])
		for v != 0 {
			index := bits.TrailingZeros(v)
			v &^= 1 << index
			id = ComponentId(i*bits.UintSize + index)
			if !yield(id) {
				return
			}
		}
	}
}
