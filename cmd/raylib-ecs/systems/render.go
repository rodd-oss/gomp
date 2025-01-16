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

func (s *renderController) Update(world *ecs.World) {
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

		rec := rl.NewRectangle(float32(transform.X), float32(transform.Y), 8, 8)
		origin := rl.NewVector2(0, 0)

		rl.DrawRectanglePro(rec, origin, 0, *color)
		return true
	})

	rl.DrawFPS(10, 10)
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.EndDrawing()
}
func (s *renderController) FixedUpdate(world *ecs.World) {}

func (s *renderController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
