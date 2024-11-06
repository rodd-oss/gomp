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

var PhysicsSystem = CreateSystem(new(physicsSystemController))

// physicsSystemController is a system that updates the physics of a game
type physicsSystemController struct {
	world donburi.World
	space *cp.Space
}

func (c *physicsSystemController) Init(world donburi.World) {
	c.space = cp.NewSpace()
	c.world = world

	PhysicsComponent.Each(c.world, func(e *donburi.Entry) {
		component := PhysicsComponent.Get(e)

		// body := cp.NewKinematicBody()

		// randX := 100 + (rand.Float64()+0.5)*100
		// randY := 100 + (rand.Float64()-0.5)*100

		// body.SetPosition(cp.Vector{X: randX, Y: randY})

		c.space.AddBody(component.Body)
		// component.Body = body
	})
}

func (c *physicsSystemController) Update(dt float64) {
	PhysicsComponent.Each(c.world, func(e *donburi.Entry) {
		log.Println(e)
		p := PhysicsComponent.Get(e)

		if p.Body.IsSleeping() {
			log.Println("is sleeping")
			return
		}

		randX := (rand.Float64()) * 100
		randY := (rand.Float64()) * 100

		p.Body.SetPosition(cp.Vector{X: randX, Y: randY})

		pos := p.Body.Position()

		log.Println(pos)
	})

	c.space.EachBody(func(body *cp.Body) {
		log.Println(body)
	})

	c.space.Step(dt)
	log.Println(c.world.Len())
}
