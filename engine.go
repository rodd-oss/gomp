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
	"log"
	"time"
)

const (
	MaxFrameSkips = 10
)

func NewEngine(game AnyGame) Engine {
	engine := Engine{
		Game: game,
	}

	return engine
}

type Engine struct {
	Game AnyGame

	lastFixedUpdateAt time.Time
}

func (e *Engine) Run(tickrate uint, framerate uint) {

	fixedUpdDuration := time.Second / time.Duration(tickrate)

	var renderTicker *time.Ticker
	if framerate > 0 {
		renderTicker = time.NewTicker(time.Second / time.Duration(framerate))
		defer renderTicker.Stop()
	}

	e.Game.Init()
	defer e.Game.Destroy()

	lastUpdateAt := time.Now() // TODO: REMOVE?
	nextFixedUpdateAt := time.Now()

	for !e.Game.ShouldDestroy() {
		if renderTicker != nil {
			<-renderTicker.C
		}

		// Update
		e.Game.Update(time.Since(lastUpdateAt))

		// Fixed Update
		loops := 0
		for nextFixedUpdateAt.Compare(time.Now()) == -1 && loops < MaxFrameSkips {
			e.Game.FixedUpdate(fixedUpdDuration)
			e.lastFixedUpdateAt = nextFixedUpdateAt
			nextFixedUpdateAt = nextFixedUpdateAt.Add(fixedUpdDuration)
			loops++
		}
		if loops >= MaxFrameSkips {
			log.Println("Too many updates detected")
		}

		// Render

		sinceLastFixedUpdateAt := time.Since(e.lastFixedUpdateAt)
		interpolation := float32(sinceLastFixedUpdateAt.Microseconds()) / float32(fixedUpdDuration.Microseconds())
		e.Game.Render(interpolation)
		lastUpdateAt = time.Now()
	}
}
