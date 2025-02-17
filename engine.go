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
	Destroy()
	ShouldDestroy() bool
}

func NewEngine(game AnyGame) *Engine {
	newGame := Engine{
		Game: game,
	}

	return &newGame
}

type Engine struct {
	Game AnyGame
}

func (e *Engine) Run(tickrate uint) {
	duration := time.Second / time.Duration(tickrate)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	var (
		t       time.Time
		dt      time.Duration
		fixedDt time.Duration
	)

	e.Game.Init()
	defer e.Game.Destroy()

	for !e.Game.ShouldDestroy() {
		needFixedUpdate := true
		for needFixedUpdate {
			select {
			default:
				needFixedUpdate = false
			case <-ticker.C:
				t = time.Now()
				e.Game.FixedUpdate(fixedDt)
				fixedDt = time.Since(t)
			}
		}
		t = time.Now()
		e.Game.Update(dt)
		dt = time.Since(t)
	}
}
