/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"math"
	"testing"
)

func TesCalcIndex(t *testing.T) {
	for i := 0; i <= 100_000_000; i++ {
		value := i>>10 + 1
		want := int(math.Log2(float64(value)))
		have := FastIntLog2(value)
		if want != have {
			t.Fatalf("i: %v, want: %v, got: %v", i, want, have)
		}
	}
}

func BenchmarkFastLog2(t *testing.B) {
	for range t.N {
		_ = FastIntLog2(t.N/10 + 1)
	}
}

func BenchmarkStdMathLog2(t *testing.B) {
	for range t.N {
		_ = int(math.Log2(float64(t.N/10 + 1)))
	}
}
