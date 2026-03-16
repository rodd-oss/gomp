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
	"github.com/negrel/assert"
	"gomp/pkg/worker"
)

// NewAccumulator creates a new accumulator instance
func NewAccumulator[T any](
	initFn func() T,
	resetFn func(T) T,
	mergeFn func([]T) T,
	processFn func(T),
) Accumulator[T] {
	return Accumulator[T]{
		initFn:    initFn,
		resetFn:   resetFn,
		mergeFn:   mergeFn,
		processFn: processFn,
	}
}

// Accumulator handles multi-worker accumulation with flexible merge strategies
type Accumulator[T any] struct {
	workerData    []T         // Per-worker storage
	initFn        func() T    // Initializes worker-specific storage
	resetFn       func(T) T   // Clears worker data while preserving allocations
	mergeFn       func([]T) T // Combines all worker data into final result
	processFn     func(T)     // Processes worker data
	isInitialized bool
}

// Init initializes the accumulator
func (a *Accumulator[T]) Init(pool *worker.Pool) {
	a.workerData = make([]T, pool.NumWorkers())
	for i := range a.workerData {
		a.workerData[i] = a.initFn()
	}
	a.isInitialized = true
}

// Update modifies worker-specific data in a thread-safe manner
func (a *Accumulator[T]) Update(workerID worker.WorkerId, update func(acc T) T) {
	assert.True(a.isInitialized)
	a.workerData[workerID] = update(a.workerData[workerID])
}

// Merge combines all worker data into final result
func (a *Accumulator[T]) Merge() T {
	assert.True(a.isInitialized)
	return a.mergeFn(a.workerData)
}

// Process applies finalization function to all collected data
func (a *Accumulator[T]) Process() {
	assert.True(a.isInitialized)
	for i := range a.workerData {
		a.processFn(a.workerData[i])
	}
}

// Reset prepares the accumulator for new data collection
func (a *Accumulator[T]) Reset() {
	assert.True(a.isInitialized)
	for i := range a.workerData {
		a.workerData[i] = a.resetFn(a.workerData[i])
	}
}

// Destroy releases all resources associated with the accumulator
func (a *Accumulator[T]) Destroy() {
	a.workerData = nil
	a.initFn = nil
	a.resetFn = nil
	a.mergeFn = nil
	a.processFn = nil
	a.isInitialized = false
}
