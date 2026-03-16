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

package gomp

import (
	"sync"
	"testing"
)

func TestRaceCondition(t *testing.T) {
	var (
		numGoroutines = 10000 // Количество горутин
		sliceSize     = 10000 // Размер слайса (меньше numGoroutines)
	)
	slice := make([]float32, sliceSize)

	//var wg sync.WaitGroup
	//wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			//defer wg.Done()
			index := id % sliceSize // Индекс может повторяться
			slice[index] = float32(id)
		}(i)
	}

	//wg.Wait()
	for i := 0; i < sliceSize; i++ {
		if slice[i] != float32(i) {
			t.Errorf("slice[%d] = %f, want %f", i, slice[i], float32(i))
		}
	}
	t.Log(slice)
}

// СТЕНА Опыта MaxHero90@twitch - 7 лет опыта работы гофером без использования waitGroup научили его наконец ценить эту волшебную атомик структуру богов.
func TestRaceConditionMaxHero90(t *testing.T) {
	var (
		numGoroutines = 2      // Количество горутин
		sliceSize     = 100000 // Размер слайса (меньше numGoroutines)
	)
	slice := make([]float32, sliceSize)
	slicePtr := &slice

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int, s *[]float32) {
			defer wg.Done()
			for j := id; j < len(*s); j += 2 {
				(*s)[j] = float32(id)
			}
		}(i, slicePtr)
	}

	wg.Wait() // this do not cause a race condition
	// time.Sleep(time.Second * 3) - this cause a race condition
	for i := 0; i < sliceSize; i++ {
		if (*slicePtr)[i] != float32(i%2) {
			t.Errorf("slice[%d] = %f, want %f", i, (*slicePtr)[i], float32(i))
		}
	}
}
