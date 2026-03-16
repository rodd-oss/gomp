/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"math"
	"math/bits"
	"testing"
)

func TestCalcIndex(t *testing.T) {
	for i := uint64(0); i <= 100_000_000; i++ {
		var value uint64 = i>>10 + 1
		want := int(math.Log2(float64(value)))
		want2 := int(bits.LeadingZeros64(value) ^ 63)
		have := bits.Len64(value) - 1
		if want != have {
			t.Fatalf("i: %v, want: %v, got: %v", i, want, have)
		}
		if want2 != have {
			t.Fatalf("i: %v, want: %v, got: %v", i, want2, have)
		}
	}
}

func BenchmarkFastestLog2(b *testing.B) {
	var i uint64 = 1
	for b.Loop() {
		_ = FastestIntLog2(i/10 + 1)
		i++
	}
}

func BenchmarkFastLog2(b *testing.B) {
	var i uint64 = 1
	for b.Loop() {
		_ = FastIntLog2(i/10 + 1)
		i++
	}
}

func BenchmarkStdMathLog2(b *testing.B) {
	var i uint64 = 1
	for b.Loop() {
		_ = math.Log2(float64(i/10 + 1))
		i++
	}
}
