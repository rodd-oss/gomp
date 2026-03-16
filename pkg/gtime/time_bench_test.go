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

package gtime

import (
	"gomp/pkg/timeasm"
	"testing"
	"time"
)

const loadsize = 0

func workLoad(a int64) int64 {
	var result int64 = 0
	for range a {
		result += a * a
	}
	return result
}

func BenchmarkAsm2(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		n := timeasm.Now()
		_ = workLoad(loadsize)
		_ = timeasm.Now() - n
	}
}

func BenchmarkAsm(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		n := timeasm.Cputicks()
		_ = workLoad(loadsize)
		_ = timeasm.Cputicks() - n
	}
}

func BenchmarkC(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		n := GetTimestampC()
		_ = workLoad(loadsize)
		_ = GetTimestampC() - n
	}
}

func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		n := time.Now()
		_ = workLoad(loadsize)
		_ = time.Since(n)
	}
}

func BenchmarkTimeNowUnix(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		n := time.Now().Unix()
		_ = workLoad(loadsize)
		_ = time.Now().Unix() - n
	}
}

// Опционально: сравнение с наносекундной точностью
func BenchmarkTimeNowUnixNano(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		n := time.Now().UnixNano()
		_ = workLoad(loadsize)
		_ = time.Now().UnixNano() - n
	}
}
