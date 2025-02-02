/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"math/bits"
)

// type BitArray []uint64

// // New creates a new BitArray with the given number of bits.
// func NewBitArray(size uint64) BitArray {
// 	words := (size + 63) / 64
// 	return make(BitArray, words)
// }

// // ensureIndex ensures the given bit index is accessible, resizing the array if necessary.
// func (b *BitArray) ensureIndex(index uint64) {
// 	wordsNeeded := (index / 64) + 1
// 	if uint64(len(*b)) < wordsNeeded {
// 		newData := make(BitArray, wordsNeeded)
// 		copy(newData, *b)
// 		*b = newData
// 	}
// }

// // Resize adjusts the BitArray size to accommodate a specific number of bits.
// // Shrinks the underlying array if the new size is smaller.
// func (b *BitArray) Resize(newSize uint64) {
// 	newWords := (newSize + 63) / 64
// 	currentWords := uint64(len(*b))
// 	if newWords > currentWords {
// 		newData := make(BitArray, newWords)
// 		copy(newData, *b)
// 		*b = newData
// 	} else if newWords < currentWords {
// 		*b = (*b)[:newWords]
// 	}
// }
// // Size calculates the total number of bits in the BitArray.
// func (b BitArray) Size() uint64 {
// 	return uint64(len(b)) * 64
// }

// ComponentBitArray256 is a structure to manage an array of uint64 values as a bit array.
const bit_array_size = 256 / bits.UintSize

type ComponentBitArray256 [bit_array_size]uint

// Set sets the bit at the given index to 1.
func (b *ComponentBitArray256) Set(index ComponentID) {
	b[index/bits.UintSize] |= 1 << (index % bits.UintSize)
}

// Unset clears the bit at the given index (sets it to 0).
func (b *ComponentBitArray256) Unset(index ComponentID) {
	b[index/bits.UintSize] &^= 1 << (index % bits.UintSize)
}

// Toggle toggles the bit at the given index.
func (b *ComponentBitArray256) Toggle(index ComponentID) {
	b[index/bits.UintSize] ^= 1 << (index % bits.UintSize)
}

// IsSet checks if the bit at the given index is set (1). Automatically resizes if the index is out of bounds.
func (b *ComponentBitArray256) IsSet(index ComponentID) bool {
	return (b[index/bits.UintSize] & (1 << (index % bits.UintSize))) != 0
}

// IncludesAll checks that all bits in other are set in b
func (b *ComponentBitArray256) IncludesAll(other ComponentBitArray256) bool {
	return ((b[0] & other[0]) == other[0]) &&
		((b[1] & other[1]) == other[1]) &&
		((b[2] & other[2]) == other[2]) &&
		((b[3] & other[3]) == other[3])
}

// IncludesAll checks that any bits in other are set in b
func (b *ComponentBitArray256) IncludesAny(other ComponentBitArray256) bool {
	return ((b[0] & other[0]) > 0) ||
		((b[1] & other[1]) > 0) ||
		((b[2] & other[2]) > 0) ||
		((b[3] & other[3]) > 0)
}

func (b *ComponentBitArray256) AllSet(yield func(ComponentID) bool) {
	var id ComponentID
	var raisedBitsCount int
	for i, v := range b {
		raisedBitsCount = bits.OnesCount(v)
		for range raisedBitsCount {
			index := bits.Len(v) - 1
			v &^= 1 << index
			id = ComponentID(i*bits.UintSize + index)
			if !yield(id) {
				return
			}
		}
	}
}
