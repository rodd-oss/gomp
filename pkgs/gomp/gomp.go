/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"gomp_game/pkgs/gomp/ecs"
	"time"

	"github.com/yohamta/donburi"
)

func NewGame(tickRate time.Duration) *Game {
	game := new(Game)

	game.Init(tickRate)

	return game
}

// func CreateEntity(components ...ecs.IComponent) func() ecs.Entity {
// 	return func() ecs.Entity {
// 		return components
// 	}
// }

type Component struct {
	ecs.IComponent
	set func(*donburi.Entry)
}

func CreateEntity(components ...Component) ecs.Entity {
	cmpnnts := make([]ecs.IComponent, len(components))
	return func(world donburi.World, extra ...ecs.IComponent) {
		for i, c := range components {
			cmpnnts[i] = c.IComponent
		}

		entity := world.Create(append(cmpnnts, extra...)...)
		entry := world.Entry(entity)
		for _, c := range components {
			c.set(entry)
		}
	}
}

type ComponentFactory[T any] struct {
	Query *donburi.ComponentType[T]
}

func CreateComponent[T any]() ComponentFactory[T] {
	return ComponentFactory[T]{Query: donburi.NewComponentType[T]()}
}

func (cf ComponentFactory[T]) New(data T) Component {
	return Component{
		cf.Query,
		func(entity *donburi.Entry) {
			cf.Query.SetValue(entity, data)
		},
	}
}

// func CreateComponent[T any](initData T) *donburi.ComponentType[T] {
// 	return donburi.NewComponentType[T](initData)
// }

var systemId uint16 = 0

func CreateSystem(controller ecs.SystemController) ecs.System {
	sys := ecs.System{
		ID:         systemId,
		Controller: controller,
	}

	systemId++

	return sys
}
