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

	BodyComponent.Each(c.world, func(e *donburi.Entry) {
		body := cp.NewKinematicBody()
		BodyComponent.Set(e, body)

		randX := 100 + (rand.Float64()+0.5)*100
		randY := 100 + (rand.Float64()-0.5)*100

		body.SetPosition(cp.Vector{X: randX, Y: randY})

		c.space.AddBody(body)
	})
}

func (c *bodySystemController) Update(dt float64) {
	BodyComponent.Each(c.world, func(e *donburi.Entry) {
		body := BodyComponent.Get(e)

		if body.IsSleeping() {
			log.Println("is sleeping")
			return
		}

		randX := (rand.Float64()) * 100
		randY := (rand.Float64()) * 100

		body.SetPosition(cp.Vector{X: randX, Y: randY})
	})

	c.space.Step(dt)
}
