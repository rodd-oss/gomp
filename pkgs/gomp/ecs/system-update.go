/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type AnyUpdateSystem[W any] interface {
	Init(*W)
	Run(*W)
	Destroy(*W)
}

type UpdateSystemBuilder[W any] struct {
	world   *W
	systems *[][]AnyUpdateSystem[W]
}

func (b *UpdateSystemBuilder[W]) Sequential(systems ...AnyUpdateSystem[W]) *UpdateSystemBuilder[W] {
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.world)
		parallelSystems := make([]AnyUpdateSystem[W], 0)
		parallelSystems = append(parallelSystems, systems[i])
		*b.systems = append(*b.systems, parallelSystems)
	}
	return b
}

func (b *UpdateSystemBuilder[W]) Parallel(systems ...AnyUpdateSystem[W]) *UpdateSystemBuilder[W] {
	*b.systems = append(*b.systems, systems)
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.world)
	}
	return b
}
