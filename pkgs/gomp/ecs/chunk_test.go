/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"testing"
)

// const (
// 	BUFFER_SIZE    = 1
// 	CHUNK_CAPACITY = 2
// )

const (
	BUFFER_SIZE    = 8
	BUFFER_SIZE1   = (221>>5 | 1)
	CHUNK_CAPACITY = 1 << 14
)

func BenchmarkChunkDelete(b *testing.B) {
	b.ReportAllocs()

	w := testStruct2{
		ID: 0,
	}

	test := testStruct{
		W: &w,
	}

	chunk := NewChunkArray[testStruct](BUFFER_SIZE, CHUNK_CAPACITY)
	for range b.N {
		test.X = float32(b.N)
		test.Y = float32(b.N * 2)
		test.Z = float32(b.N * 3)

		chunk.Append(test)
	}

	b.ResetTimer()
	for range b.N {
		chunk.SoftReduce()
	}
	chunk.Clean()
}

func BenchmarkChunkUpdate(b *testing.B) {
	b.ReportAllocs()

	w := testStruct2{
		ID: 0,
	}

	test := testStruct{
		W: &w,
	}

	chunk := NewChunkArray[testStruct](BUFFER_SIZE, CHUNK_CAPACITY)
	for range b.N {
		test.X = float32(b.N)
		test.Y = float32(b.N * 2)
		test.Z = float32(b.N * 3)

		chunk.Append(test)
	}

	b.ResetTimer()
	for _, v := range chunk.Iter() {
		v.X = 0
		v.Y = 0
		v.Z = 0
	}
}

func BenchmarkChunk10mil(b *testing.B) {
	w := testStruct2{
		ID: 0,
	}

	test := testStruct{
		W: &w,
	}

	b.ResetTimer()
	for range b.N {
		chunk := NewChunkArray[testStruct](BUFFER_SIZE, CHUNK_CAPACITY)
		for range 10_000_000 {
			test.X = float32(1)
			test.Y = float32(2 * 2)
			test.Z = float32(3 * 3)

			chunk.Append(test)
		}
	}
}

func BenchmarkChunkAppend(b *testing.B) {
	b.ReportAllocs()

	w := testStruct2{
		ID: 0,
	}

	test := testStruct{
		W: &w,
	}

	chunk := NewChunkArray[testStruct](BUFFER_SIZE, CHUNK_CAPACITY)

	b.ResetTimer()
	for range b.N {
		test.X = float32(b.N)
		test.Y = float32(b.N * 2)
		test.Z = float32(b.N * 3)

		chunk.Append(test)
	}
}
