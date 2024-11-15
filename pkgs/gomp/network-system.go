/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"github.com/yohamta/donburi"
)

var NetworkSystem = CreateSystem(new(networkController))

type networkController struct {
	world         donburi.World
	addPlayers    chan int
	removePlayers chan int
	newEvent      chan int
	sendEvent     chan int
}

func (c *networkController) Init(game *Game) {
	c.world = game.World
}

func (c *networkController) Update(dt float64) {
	for i := 0; i < len(c.addPlayers); i++ {
		playerId := <-c.addPlayers

		network := NetworkComponent.New(NetworkData{
			PlayerID: playerId,
		})

		player := CreateEntity(network)(1)
		playerLen := len(player)
		for i := 0; i < playerLen; i++ {
			player[i](c.world)
		}

		if i >= 9 {
			break
		}
	}

	for i := 0; i < len(c.removePlayers); i++ {
		playerId := <-c.removePlayers

		NetworkComponent.Query.Each(c.world, func(e *donburi.Entry) {
			network := NetworkComponent.Query.Get(e)

			if network.PlayerID != playerId {
				return
			}

			ent := e.Entity()
			c.world.Remove(ent)
		})

		if i >= 9 {
			break
		}
	}

	for i := 0; i < len(c.newEvent); i++ {
		// event := <-c.newEvent

	}

	for i := 0; i < len(c.sendEvent); i++ {
		// event := <-c.sendEvent

	}
}
