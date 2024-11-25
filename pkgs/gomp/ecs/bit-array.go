/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

// BitArray is a structure to manage an array of uint64 values as a bit array.
type BitArray struct {
	data []uint64
	size uint // total number of bits
}

// New creates a new BitArray with the given number of bits.
func NewBitArray(size uint) BitArray {
	words := (size + 63) / 64
	return BitArray{
		data: make([]uint64, words),
		size: size,
	}
}

// ensureIndex ensures the given bit index is within bounds.
func (b *BitArray) ensureIndex(index uint) {
	if index >= b.size {
		b.Resize(index)
	}
}

// Set sets the bit at the given index to 1.
func (b *BitArray) Set(index uint) {
	b.ensureIndex(index)
	b.data[index/64] |= 1 << (index % 64)
}

// Clear clears the bit at the given index (sets it to 0).
func (b *BitArray) Clear(index uint) {
	b.ensureIndex(index)
	b.data[index/64] &^= 1 << (index % 64)
}

// Toggle toggles the bit at the given index.
func (b *BitArray) Toggle(index uint) {
	b.ensureIndex(index)
	b.data[index/64] ^= 1 << (index % 64)
}

// IsSet checks if the bit at the given index is set (1).
func (b *BitArray) IsSet(index uint) bool {
	b.ensureIndex(index)
	return (b.data[index/64] & (1 << (index % 64))) != 0
}

// Resize resizes the BitArray to a new size. Data is preserved for overlapping bits.
func (b *BitArray) Resize(newSize uint) {
	newWords := (newSize + 63) / 64
	if int(newWords) != len(b.data) {
		newData := make([]uint64, newWords)
		copy(newData, b.data)
		b.data = newData
	}
	b.size = newSize
}

// Size returns the total number of bits in the BitArray.
func (b *BitArray) Size() uint {
	return b.size
}
