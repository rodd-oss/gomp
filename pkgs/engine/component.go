/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	capnp "capnproto.org/go/capnp/v3"
	ecs "github.com/yohamta/donburi"
)

type Controller interface {
	Init()
	Update(dt float64)
}

type NetworkController[T any] interface {
	Init()
	Update(dt float64)
	CreateState(s *capnp.Segment) (T, error)
	OnStateRequest(T) T
	OnStateUpdate(T)
}

type WrappedNetworkController[T any] struct {
	Controller NetworkController[T]
	State      T
	Patch      T
}

type Component[T any] struct {
	System  *ecs.ComponentType[WrappedNetworkController[T]]
	Type_ID uint64
}

func CreateNetworkComponent[State any](stateType_ID uint64, controller NetworkController[State]) *Component[State] {
	component := new(Component[State])

	stateArena := capnp.SingleSegment(nil)
	_, stateSeg, err := capnp.NewMessage(stateArena)
	if err != nil {
		panic(err)
	}
	state, err := controller.CreateState(stateSeg)
	if err != nil {
		panic(err)
	}

	patchArena := capnp.SingleSegment(nil)
	_, patchSeg, err := capnp.NewMessage(patchArena)
	if err != nil {
		panic(err)
	}
	patch, err := controller.CreateState(patchSeg)
	if err != nil {
		panic(err)
	}

	wrappedController := WrappedNetworkController[State]{
		Controller: controller,
		State:      state,
		Patch:      patch,
	}

	component.System = ecs.NewComponentType[WrappedNetworkController[State]](wrappedController)
	component.Type_ID = stateType_ID

	return component
}

func CreateComponent(controller Controller) *ecs.ComponentType[Controller] {
	component := ecs.NewComponentType[Controller]()
	return component
}
