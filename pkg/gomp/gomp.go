/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yohamta/donburi"
)

func NewGame(tickRate time.Duration) *Game {
	game := new(Game)

	game.Init(tickRate)

	return game
}

func CreateEntity(components ...Component) func(amount int) []Entity {
	return func(amount int) []Entity {
		if amount <= 0 {
			panic(fmt.Sprint("Adding Entity to scene with (", amount, ") amount failed. Amount must be greater than 0."))
		}

		entArr := make([]Entity, amount)
		ent := func(world donburi.World, extra ...Component) {
			components := append(components, extra...)
			cmpnnts := make([]IComponent, len(components))
			for i, c := range components {
				cmpnnts[i] = c.ComponentType
			}

			entity := world.Create(cmpnnts...)
			entry := world.Entry(entity)
			for _, c := range components {
				c.Set(entry)
			}
		}

		for i := 0; i < amount; i++ {
			entArr[i] = ent
		}

		return entArr
	}
}

type ComponentFactory[T any] struct {
	Query *donburi.ComponentType[T]
}

func CreateComponent[T any]() ComponentFactory[T] {
	typeFor := reflect.TypeFor[T]()

	if typeFor.Kind() == reflect.Interface {
		panic(fmt.Sprint("CreateComponent[", typeFor.String(), "] failed. Type must not be an interface."))
	}

	return ComponentFactory[T]{Query: donburi.NewComponentType[T]()}
}

func (cf ComponentFactory[T]) New(data T) Component {
	return Component{
		ComponentType: cf.Query,
		Set: func(entity *donburi.Entry) {
			cf.Query.SetValue(entity, data)
		},
	}
}

// func CreateComponent[T any](initData T) *donburi.ComponentType[T] {
// 	return donburi.NewComponentType[T](initData)
// }

var systemId uint16 = 0

func CreateSystem(controller SystemController) System {
	sys := System{
		ID:         systemId,
		Controller: controller,
	}

	systemId++

	return sys
}
