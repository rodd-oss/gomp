/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type renderController struct {
	width, height int32
}

func (s *renderController) Init(world *ecs.World) {
	rl.InitWindow(s.width, s.height, "raylib [core] example - basic window")
}
func (s *renderController) Run(world *ecs.World) {
	colors := components.ColorService.GetManager(world)
	transforms := components.TransformService.GetManager(world)

	if rl.WindowShouldClose() {
		world.SetShouldDestroy(true)
		return
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	colors.All(func(entity ecs.EntityID, color *color.RGBA) bool {
		if color == nil {
			return true
		}

		transform := transforms.GetPtr(entity)
		if transform == nil {
			return true
		}

		rl.DrawRectangle(int32(transform.X), int32(transform.Y), 8, 8, rl.NewColor(color.R, color.G, color.B, color.A))
		return true
	})

	rl.DrawFPS(10, 10)
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.EndDrawing()
}

func (s *renderController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
