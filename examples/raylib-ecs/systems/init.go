/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkgs/ecs"
)

type initController struct {
	windowWidth, windowHeight int32
}

func (s *initController) Init(world *ecs.World) {
	rl.InitWindow(s.windowWidth, s.windowHeight, "raylib [core] ebiten-ecs - basic window")
}

func (s *initController) Update(world *ecs.World)      {}
func (s *initController) FixedUpdate(world *ecs.World) {}
func (s *initController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
