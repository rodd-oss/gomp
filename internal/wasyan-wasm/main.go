//go:build wasm

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

package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOARCH, runtime.GOOS)
}

//go:wasmexport add
func add(a, b uint32) {
	//panic("lol")
	//var game wasyan.Game
	//GetGame(&game)
	//wasyan.UpdateGame(&game, vectors.Vec2{X: float32(a), Y: float32(b)})

	//SetGame(&game)
	var _ = a + b
}
