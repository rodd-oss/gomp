/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

// BitArray is a structure to manage an array of uint64 values as a bit array.
type BitArray []uint64

// New creates a new BitArray with the given number of bits.
func NewBitArray(size uint64) BitArray {
	words := (size + 63) / 64
	return make(BitArray, words)
}

// ensureIndex ensures the given bit index is accessible, resizing the array if necessary.
func (b *BitArray) ensureIndex(index uint64) {
	wordsNeeded := (index / 64) + 1
	if uint64(len(*b)) < wordsNeeded {
		newData := make(BitArray, wordsNeeded)
		copy(newData, *b)
		*b = newData
	}
}

// Size calculates the total number of bits in the BitArray.
func (b BitArray) Size() uint64 {
	return uint64(len(b)) * 64
}

// Set sets the bit at the given index to 1, resizing if necessary.
func (b *BitArray) Set(index uint64) {
	b.ensureIndex(index)
	(*b)[index/64] |= 1 << (index % 64)
}

// Clear clears the bit at the given index (sets it to 0), resizing if necessary.
func (b *BitArray) Clear(index uint64) {
	b.ensureIndex(index)
	(*b)[index/64] &^= 1 << (index % 64)
}

// Toggle toggles the bit at the given index, resizing if necessary.
func (b *BitArray) Toggle(index uint64) {
	b.ensureIndex(index)
	(*b)[index/64] ^= 1 << (index % 64)
}

// IsSet checks if the bit at the given index is set (1). Automatically resizes if the index is out of bounds.
func (b *BitArray) IsSet(index uint64) bool {
	b.ensureIndex(index)
	return ((*b)[index/64] & (1 << (index % 64))) != 0
}

// Resize adjusts the BitArray size to accommodate a specific number of bits.
// Shrinks the underlying array if the new size is smaller.
func (b *BitArray) Resize(newSize uint64) {
	newWords := (newSize + 63) / 64
	currentWords := uint64(len(*b))
	if newWords > currentWords {
		newData := make(BitArray, newWords)
		copy(newData, *b)
		*b = newData
	} else if newWords < currentWords {
		*b = (*b)[:newWords]
	}
}
