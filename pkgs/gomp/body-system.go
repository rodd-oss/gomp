/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"log"
	"math/rand/v2"

	"github.com/jakecoffman/cp/v2"
	"github.com/yohamta/donburi"
)

var BodySystem = CreateSystem(new(bodySystemController))

// physicsSystemController is a system that updates the physics of a game
type bodySystemController struct {
	world donburi.World
	space *cp.Space
}

func (c *bodySystemController) Init(world donburi.World) {
	c.space = cp.NewSpace()
	c.world = world

	BodyComponent.Query.Each(c.world, func(e *donburi.Entry) {
		body := BodyComponent.Query.Get(e)

		randX := (rand.Float64()) * 1000
		randY := (rand.Float64()) * 1000

		body.SetPosition(cp.Vector{X: randX, Y: randY})

		// randX = (rand.Float64() - 0.5) * 10
		// randY = (rand.Float64() - 0.5) * 10

		// body.SetVelocity(randX, randY)

		c.space.AddBody(body)
	})
}

func (c *bodySystemController) Update(dt float64) {
	BodyComponent.Query.Each(c.world, func(e *donburi.Entry) {
		body := BodyComponent.Query.Get(e)

		if body.IsSleeping() {
			log.Println("is sleeping")
			return
		}
	})

	c.space.Step(dt)
}
