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
	"math/rand"
	"sync"
	"testing"
)

const (
	keysNumberOfBases = 1 << 16
	keysOffsetRange   = 2 << 16
)

// Generate a key that is close to others but spans a wide range
func generateKeys(numBases int) func() uint32 {
	// Create an array of base keys
	bases := make([]uint32, numBases)
	for i := 0; i < numBases; i++ {
		bases[i] = rand.Uint32()
	}

	// Offset range (e.g., keys will be within ±1000 of a base)
	offsetRange := uint32(keysOffsetRange)

	return func() uint32 {
		// Randomly select a base key
		base := bases[rand.Intn(numBases)]
		// Generate a random offset within the range
		offset := rand.Uint32() % offsetRange
		// Randomly add or subtract the offset
		if rand.Intn(2) == 0 {
			return base + offset
		}
		return base - offset
	}
}

// Benchmark for default Go map Set with close keys (multiple bases)
func Benchmark_Map_Set(b *testing.B) {
	b.ReportAllocs()

	m := make(map[uint32]int)
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[keys[i]] = i
	}
}

// Benchmark for LookupMap Set with close keys (multiple bases)
func Benchmark_LookupMap_Set(b *testing.B) {
	b.ReportAllocs()

	lm := &LookupMap[int]{}
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lm.Set(keys[i], i)
	}
}

// Benchmark for default Go map Get with close keys (multiple bases)
func Benchmark_Map_Get(b *testing.B) {
	b.ReportAllocs()

	m := make(map[uint32]int)
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
		m[keys[i]] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[keys[i]]
	}
}

// Benchmark for LookupMap Get with close keys (multiple bases)
func Benchmark_LookupMap_Get(b *testing.B) {
	b.ReportAllocs()

	lm := &LookupMap[int]{}
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
		lm.Set(keys[i], i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lm.Get(keys[i])
	}
}

// Benchmark for default Go map concurrent Set with close keys (multiple bases, with sync.Mutex)
func Benchmark_Map_Concurrent_Set(b *testing.B) {
	b.ReportAllocs()

	m := make(map[uint32]int)
	var mx sync.Mutex
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mx.Lock()
			m[keys[i]] = i
			mx.Unlock()
			i++
		}
	})
}

// Benchmark for LookupMap concurrent Set with close keys (multiple bases)
func Benchmark_LookupMap_Concurrent_Set(b *testing.B) {
	b.ReportAllocs()

	lm := &LookupMap[int]{}
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lm.Set(keys[i], i)
			i++
		}
	})
}

// Benchmark for default Go map concurrent Get with close keys (multiple bases, with sync.RWMutex)
func Benchmark_Map_Concurrent_Get(b *testing.B) {
	b.ReportAllocs()

	m := make(map[uint32]int)
	var mx sync.RWMutex
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
		m[keys[i]] = i
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mx.RLock()
			_ = m[keys[i]]
			mx.RUnlock()
			i++
		}
	})
}

// Benchmark for LookupMap concurrent Get with close keys (multiple bases)
func Benchmark_LookupMap_Concurrent_Get(b *testing.B) {
	b.ReportAllocs()

	lm := &LookupMap[int]{}
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	keys := make([]uint32, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = gen()
		lm.Set(keys[i], i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lm.Get(keys[i])
			i++
		}
	})
}

func BenchmarkAll(b *testing.B) {
	b.Run("Set", BenchmarkSet)
	b.Run("Get", BenchmarkGet)
	b.Run("Set_Concurrent", BenchmarkSetConcurrent)
	b.Run("Get_Concurrent", BenchmarkGetConcurrent)
}

func BenchmarkSet(b *testing.B) {
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	b.Run("Map_Set", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		m := make(map[uint32]int)
		b.ResetTimer()
		for i := 0; i < len(keys); i++ {
			m[keys[i]] = i
		}
	})
	b.Run("LookupMap_Set", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		lm := NewLookupMap[int]()
		b.ResetTimer()
		for i := 0; i < len(keys); i++ {
			lm.Set(keys[i], i)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	b.Run("Map_Get", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		m := make(map[uint32]int)
		for i := 0; i < len(keys); i++ {
			m[keys[i]] = i
		}
		b.ResetTimer()
		for i := 0; i < len(keys); i++ {
			_ = m[keys[i]]
		}
	})
	b.Run("LookupMap_Get", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		lm := NewLookupMap[int]()
		for i := 0; i < len(keys); i++ {
			lm.Set(keys[i], i)
		}
		b.ResetTimer()
		for i := 0; i < len(keys); i++ {
			lm.Get(keys[i])
		}
	})
}

func BenchmarkSetConcurrent(b *testing.B) {
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	b.Run("Map_Concurrent_Set", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		m := make(map[uint32]int)
		var mx sync.Mutex
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				mx.Lock()
				m[keys[i]] = i
				mx.Unlock()
				i++
			}
		})
	})
	b.Run("LookupMap_Concurrent_Set", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		lm := NewLookupMap[int]()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				lm.Set(keys[i], i)
				i++
			}
		})
	})
}

func BenchmarkGetConcurrent(b *testing.B) {
	gen := generateKeys(keysNumberOfBases) // Use keysNumberOfBases base keys
	b.Run("Map_Concurrent_Get", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		m := make(map[uint32]int)
		var mx sync.RWMutex
		for i := 0; i < len(keys); i++ {
			m[keys[i]] = i
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				mx.RLock()
				_ = m[keys[i]]
				mx.RUnlock()
				i++
			}
		})
	})
	b.Run("LookupMap_Concurrent_Get", func(b *testing.B) {
		keys := make([]uint32, b.N)
		for i := 0; i < b.N; i++ {
			keys[i] = gen()
		}
		b.ReportAllocs()
		lm := NewLookupMap[int]()
		for i := 0; i < len(keys); i++ {
			lm.Set(keys[i], i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				lm.Get(keys[i])
				i++
			}
		})
	})
}
