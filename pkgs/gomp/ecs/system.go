/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type System interface {
	Init(*ECS)
	Run(*ECS)
	Destroy(*ECS)
}

type SystemBuilder struct {
	ecs *ECS
}

func (b *SystemBuilder) Sequential(systems ...System) *SystemBuilder {
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.ecs)
		parallelSystems := make([]System, 0)
		parallelSystems = append(parallelSystems, systems[i])
		b.ecs.systems = append(b.ecs.systems, parallelSystems)
	}
	return b
}

func (b *SystemBuilder) Parallel(systems ...System) *SystemBuilder {
	b.ecs.systems = append(b.ecs.systems, systems)
	for i := 0; i < len(systems); i++ {
		systems[i].Init(b.ecs)
	}
	return b
}

func runSystemAsync(system System, e *ECS) {
	defer e.wg.Done()
	system.Run(e)
}
