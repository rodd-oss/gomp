/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/entities"
	"log"

	"github.com/jakecoffman/cp/v2"
	"github.com/yohamta/donburi"
)

var PhysicsSystem = engine.CreateSystem(&physicsSystemController{})

// physicsSystemController is a system that updates the physics of a game
type physicsSystemController struct {
	world donburi.World
	space *cp.Space
}

func (c *physicsSystemController) Init(scene *engine.Scene) {
	c.space = scene.Space
	c.world = scene.World

	entities.PlayerPhysics.Each(c.world, func(e *donburi.Entry) {
		component := entities.PlayerPhysics.Get(e)

		c.space.AddBody(component.Body)

		component.Body.SetVelocity(1, 0)
	})
}

func (c *physicsSystemController) Update(dt float64) {
	entities.PlayerPhysics.Each(c.world, func(e *donburi.Entry) {
		p := entities.PlayerPhysics.Get(e)

		vel := p.Body.Position()

		log.Println(vel)
	})

	c.space.Step(dt)
}
