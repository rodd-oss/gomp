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
	"time"
)

type AnyGame interface {
	Init()
	Update(dt time.Duration)
	FixedUpdate(dt time.Duration)
	Render(dt time.Duration)
	Destroy()
	ShouldDestroy() bool
}

func NewEngine(game AnyGame) Engine {
	engine := Engine{
		Game: game,
	}

	return engine
}

type Engine struct {
	Game AnyGame

	lastUpdateAt      time.Time
	lastFixedUpdateAt time.Time
	lastRenderAt      time.Time
}

func (e *Engine) Run(tickrate uint, framerate uint) {
	fixedUpdDuration := time.Second / time.Duration(tickrate)
	framerateDuration := time.Second / time.Duration(framerate)

	fixedUpdTicker := time.NewTicker(fixedUpdDuration)
	defer fixedUpdTicker.Stop()

	renderTicker := time.NewTicker(framerateDuration)
	defer renderTicker.Stop()

	e.Game.Init()
	defer e.Game.Destroy()

	e.lastUpdateAt = time.Now()
	e.lastFixedUpdateAt = time.Now()
	e.lastRenderAt = time.Now()

	for !e.Game.ShouldDestroy() {
		// Update
		dt := time.Since(e.lastUpdateAt)
		e.Game.Update(dt)
		e.lastUpdateAt = time.Now()

		// Fixed Update
		select {
		case <-fixedUpdTicker.C:
			dt := time.Since(e.lastFixedUpdateAt)
			e.Game.FixedUpdate(dt)
			e.lastFixedUpdateAt = time.Now()
		default:
			break
		}

		// Render
		select {
		case <-renderTicker.C:
			dt := time.Since(e.lastRenderAt)
			e.Game.Render(dt)
			e.lastRenderAt = time.Now()
		default:
			break
		}
	}
}
