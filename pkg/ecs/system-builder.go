/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type SystemBuilder struct {
	world   *World
	systems *[][]*SystemServiceInstance
}

func (b *SystemBuilder) Sequential(systems ...AnySystemServicePtr) *SystemBuilder {
	for i := range systems {
		instance := systems[i].register(b.world)
		*b.systems = append(*b.systems, []*SystemServiceInstance{instance})
		instance.controller.Init(b.world)
	}

	return b
}

func (b *SystemBuilder) Parallel(systems ...AnySystemServicePtr) *SystemBuilder {
	systemInstances := []*SystemServiceInstance{}

	for i := range systems {
		instance := systems[i].register(b.world)
		systemInstances = append(systemInstances, instance)
		go instance.Run()
	}

	*b.systems = append(*b.systems, systemInstances)
	return b
}
