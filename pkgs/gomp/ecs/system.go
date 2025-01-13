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

type AnySystemManagerPtr interface {
	register(*World) AnySystemControllerPtr
}

type SystemManager[T AnySystemPtr[World]] struct {
	instance  map[*World]T
	initValue T
}

func (m *SystemManager[T]) register(word *World) AnySystemControllerPtr {
	m.instance[word] = m.initValue
	return m.instance[word]
}

func CreateSystem[T AnySystemPtr[World]](controller T, dependsOn ...AnySystemManagerPtr) SystemManager[T] {
	return SystemManager[T]{
		instance:  make(map[*World]T),
		initValue: controller,
	}
}
