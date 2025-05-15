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
	"gomp/internal/wasyan"
	"unsafe"
)

//go:wasmimport env get_game
func get_game(gameRef uint32)

//go:wasmimport env set_game
func set_game(gameRef uint32)

func GetGame(game *wasyan.Game) {
	get_game(uint32(uintptr(unsafe.Pointer(game))))
}

func SetGame(game *wasyan.Game) {
	set_game(uint32(uintptr(unsafe.Pointer(game))))
}
