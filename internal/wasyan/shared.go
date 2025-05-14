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

package wasyan

import (
	"gomp/vectors"
	"sync"
)

func Add(a, b int32) int32 {
	var c float64
	c += float64(a)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 1_000_000 {
			c += float64(b)
		}
	}()
	wg.Wait()

	return int32(c)
}

type Game struct {
	Position vectors.Vec2
	Velocity vectors.Vec2
}
