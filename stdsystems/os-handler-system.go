/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func NewOSHandlerSystem() OSHandlerSystem {
	return OSHandlerSystem{}
}

type OSHandlerSystem struct{}

func (s *OSHandlerSystem) Init() {
	// TODO: pass parameters, resize or reinit.
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "GOMP")
}

func (s *OSHandlerSystem) Run(dt time.Duration) bool {
	if rl.IsKeyPressed(rl.KeyEscape) {
		return true
	}
	if rl.WindowShouldClose() {
		return true
	}
	return false
}

func (s *OSHandlerSystem) Destroy() {
	rl.CloseWindow()
}
