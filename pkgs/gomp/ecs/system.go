/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type System interface {
	Init()
	Run(*ECS)
	Destroy()
}

type SystemBuilder struct {
	ecs *ECS
}

func (b *SystemBuilder) Sequential(systems ...System) *SystemBuilder {
	for i := 0; i < len(systems); i++ {
		system := systems[i]
		parallelSystems := make([]System, 0)
		parallelSystems = append(parallelSystems, system)
		b.ecs.Systems = append(b.ecs.Systems, parallelSystems)
	}
	return b
}

func (b *SystemBuilder) Parallel(systems ...System) *SystemBuilder {
	b.ecs.Systems = append(b.ecs.Systems, systems)
	return b
}

func runSystemAsync(system System, e *ECS) {
	defer e.wg.Done()
	system.Run(e)
}
