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

import "unsafe"

func main() {
	println("main")
}

//go:wasmexport add
func add(a, b int32) uint64 {
	game := GetGame()
	return uint64(uintptr(unsafe.Pointer(&game.Velocity.X)))
}
