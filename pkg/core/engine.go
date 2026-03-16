/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- тажефигня Donated 500 RUB

Thank you for your support!
*/

package core

import (
	"gomp/pkg/worker"
	"log"
	"runtime"
	"time"
)

const (
	MaxFrameSkips = 5
)

func NewEngine(game AnyGame) Engine {
	numCpu := max(runtime.NumCPU()-1, 1)
	engine := Engine{
		Game: game,
		pool: worker.NewPool(numCpu),
	}
	return engine
}

type Engine struct {
	Game AnyGame
	pool worker.Pool
}

func (e *Engine) Run(tickrate uint, framerate uint) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fixedUpdDuration := time.Second / time.Duration(tickrate)

	var renderTicker *time.Ticker
	if framerate > 0 {
		renderTicker = time.NewTicker(time.Second / time.Duration(framerate))
		defer renderTicker.Stop()
	}

	e.Game.Init(e)
	defer e.Game.Destroy()

	var lastUpdateAt = time.Now() // TODO: REMOVE?
	var nextFixedUpdateAt = time.Now()
	var dt = time.Since(lastUpdateAt)

	e.pool.Start()

	for !e.Game.ShouldDestroy() {
		if renderTicker != nil {
			<-renderTicker.C
		}
		dt = time.Since(lastUpdateAt)
		lastUpdateAt = time.Now()

		// Update
		e.Game.Update(dt)

		// Fixed Update
		loops := 0
		// TODO: Refactor to work without for loop
		for nextFixedUpdateAt.Compare(time.Now()) == -1 && loops < MaxFrameSkips {
			e.Game.FixedUpdate(fixedUpdDuration)
			nextFixedUpdateAt = nextFixedUpdateAt.Add(fixedUpdDuration)
			loops++
		}
		if loops >= MaxFrameSkips {
			nextFixedUpdateAt = time.Now()
			log.Println("Too many updates detected")
		}

		// RenderAssterodd
		e.Game.Render(dt)
	}
}

func (e *Engine) Pool() *worker.Pool {
	return &e.pool
}
