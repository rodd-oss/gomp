/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"flag"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 1000
	height = 1000
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	newGame := newGameClient()

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("1 mil pixel ecs")

	if err := ebiten.RunGame(&newGame); err != nil {
		panic(err)
	}
}
