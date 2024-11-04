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
	"math/rand/v2"

	"github.com/jakecoffman/cp/v2"
	"github.com/yohamta/donburi"
)

func PhysicsSystem() engine.System {
	return engine.CreateSystem(new(physicsSystemController))
}

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

		body := cp.NewKinematicBody()

		randX := 100 + (rand.Float64()-0.5)*10
		randY := 100 + (rand.Float64()-0.5)*10

		body.SetPosition(cp.Vector{X: randX, Y: randY})

		c.space.AddBody(body)
		component.Body = body
	})
}

func (c *physicsSystemController) Update(dt float64) {
	entities.PlayerPhysics.Each(c.world, func(e *donburi.Entry) {
		p := entities.PlayerPhysics.Get(e)

		if p.Body.IsSleeping() {
			log.Println("is sleeping")
			return
		}

		randX := 100 + (rand.Float64()-0.5)*10
		randY := 100 + (rand.Float64()-0.5)*10

		p.Body.SetVelocity(randX, randY)

		pos := p.Body.Position()

		log.Println(pos)
	})

	c.space.Step(dt)
	log.Println(c.world.Len())
}
