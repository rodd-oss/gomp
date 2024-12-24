/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "github.com/hajimehoshi/ebiten/v2"

type AnyDrawSystem[W any] interface {
	Init(*W)
	Run(*W, *ebiten.Image)
	Destroy(*W)
}

type DrawSystemBuilder[W any] struct {
	ecs     *W
	systems *[][]AnyDrawSystem[W]
}

func (b *DrawSystemBuilder[W]) Sequential(systems ...AnyDrawSystem[W]) *DrawSystemBuilder[W] {
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.ecs)
		parallelSystems := make([]AnyDrawSystem[W], 0)
		parallelSystems = append(parallelSystems, systems[i])
		*b.systems = append(*b.systems, parallelSystems)
	}
	return b
}

func (b *DrawSystemBuilder[W]) Parallel(systems ...AnyDrawSystem[W]) *DrawSystemBuilder[W] {
	*b.systems = append(*b.systems, systems)
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.ecs)
	}
	return b
}
