/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"gomp_game/pkgs/gomp/ecs"
	"gomp_game/pkgs/gomp/engine"
	"time"

	"github.com/yohamta/donburi"
)

func NewGame(tickRate time.Duration) *engine.Game {
	game := new(engine.Game)

	game.Init(tickRate)

	return game
}

func CreateEntity(components ...ecs.IComponent) ecs.Entity {
	return components
}

func CreateComponent[T any](initData T) *donburi.ComponentType[T] {
	return donburi.NewComponentType[T](initData)
}

func CreateSystem(controller ecs.SystemController) ecs.System {
	return ecs.System{
		Controller: controller,
	}
}
