/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "testing"

type testType uint32

type testStruct struct {
	X, Y, Z float32
	W       *testStruct2
}

type testStruct2 struct {
	ID testType
}

func BenchmarkBucketDelete(b *testing.B) {
	b.ReportAllocs()

	w := testStruct2{
		ID: 0,
	}

	test := testStruct{
		W: &w,
	}

	bucket := NewBucket[testStruct](1000)
	for range b.N {
		test.X = float32(b.N)
		test.Y = float32(b.N * 2)
		test.Y = float32(b.N * 3)

		bucket.Append(test)
	}

	b.ResetTimer()
	for range b.N {
		bucket.SoftReduce()
	}
	bucket.Clean()
	b.StopTimer()

	if len(bucket.data) != 0 {
		b.Fatal("bucket should be empty")
	}
}
