/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type AnySystemPtr[W any] interface {
	Init(*W)
	Run(*W)
	Destroy(*W)
}

type AnySystemControllerPtr interface {
	Init(*World)
	Run(*World)
	Destroy(*World)
}

type AnySystemServicePtr interface {
	register(*World) AnySystemControllerPtr
}

type SystemService[T AnySystemPtr[World]] struct {
	instance  map[*World]T
	initValue T
}

func (m *SystemService[T]) register(word *World) AnySystemControllerPtr {
	m.instance[word] = m.initValue
	return m.instance[word]
}

// TODO: dependsOn
func CreateSystem[T AnySystemPtr[World]](controller T, dependsOn ...AnySystemServicePtr) SystemService[T] {
	return SystemService[T]{
		instance:  make(map[*World]T),
		initValue: controller,
	}
}
