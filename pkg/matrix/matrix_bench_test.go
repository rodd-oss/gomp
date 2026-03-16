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

package matrix

import (
	"testing"
)

// createSampleM4 generates a sample 4x4 identity matrix.
func createSampleM4() M4 {
	return M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// createSampleData generates a slice of sample V4 vectors.
func createSampleData(size int) []V4 {
	data := make([]V4, size)
	for i := range data {
		data[i] = V4{float32(i), float32(i + 1), float32(i + 2), float32(i + 3)}
	}
	return data
}

// BenchmarkMultiply benchmarks the multiply function with a 1000-element slice.
func BenchmarkMultiply(b *testing.B) {
	m := createSampleM4()
	backupData := createSampleData(1000) // Sample data with 1000 vectors
	data := make([]V4, len(backupData))  // Preallocate to avoid allocations in loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(data, backupData) // Reset data to original state
		multiply(data, m)
	}
}

// BenchmarkMultiplyAsm benchmarks the multiply function with a 1000-element slice.
func BenchmarkMultiplyAsm(b *testing.B) {
	m := createSampleM4()
	backupData := createSampleData(1000) // Sample data with 1000 vectors
	data := make([]V4, len(backupData))  // Preallocate to avoid allocations in loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(data, backupData) // Reset data to original state
		multiplyasm(data, m)
	}
}
