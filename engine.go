/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package gomp

import (
	"gomp/pkg/ecs"
	"time"
)

type EngineSystems interface {
	Init()
	Update(dt time.Duration)
	FixedUpdate(dt time.Duration)
	Destroy()
}

type Engine[C any, S EngineSystems] struct {
	World      *ecs.World
	Components C
	Systems    S
}

func NewEngine[C any, S EngineSystems](world *ecs.World, components C, systems S) *Engine[C, S] {

	newGame := Engine[C, S]{
		World:      world,
		Components: components,
		Systems:    systems,
	}

	return &newGame
}

func (g *Engine[C, S]) Run(tickrate uint) {
	duration := time.Second / time.Duration(tickrate)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	var (
		t       time.Time
		dt      time.Duration
		fixedDt time.Duration
	)

	g.Systems.Init()
	defer g.Systems.Destroy()

	for !g.World.ShouldDestroy() {
		needFixedUpdate := true
		for needFixedUpdate {
			select {
			default:
				needFixedUpdate = false
			case <-ticker.C:
				t = time.Now()
				g.Systems.FixedUpdate(fixedDt)
				fixedDt = time.Since(t)
			}
		}
		t = time.Now()
		g.Systems.Update(dt)
		dt = time.Since(t)
	}
}
