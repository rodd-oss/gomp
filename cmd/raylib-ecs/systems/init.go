/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type initController struct {
	windowWidth, windowHeight int32
}

func (s *initController) Init(world *ecs.World) {
	rl.InitWindow(s.windowWidth, s.windowHeight, "raylib [core] example - basic window")

	currentMonitorRefreshRate := rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())
	rl.SetTargetFPS(int32(currentMonitorRefreshRate))
}

func (s *initController) Update(world *ecs.World)      {}
func (s *initController) FixedUpdate(world *ecs.World) {}
func (s *initController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
