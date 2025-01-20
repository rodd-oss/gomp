/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"fmt"
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type renderController struct {
	width, height int32
	texture       rl.Texture2D
}

func (s *renderController) Init(world *ecs.World) {
	rl.InitWindow(s.width, s.height, "raylib [core] example - basic window")

	// currentMonitorRefreshRate := rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())
	// // rl.SetTargetFPS(int32(currentMonitorRefreshRate))

	s.texture = rl.LoadTexture("assets/star.png")
}

func (s *renderController) Update(world *ecs.World) {
	spriteManager := components.SpriteService.GetManager(world)

	spriteManager.AllDataParallel(func(sprite *components.Sprite) bool {
		sprite.Texture = s.texture
		return true
	})

	if rl.WindowShouldClose() {
		world.SetShouldDestroy(true)
		return
	}

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	spriteManager.AllData(func(sprite *components.Sprite) bool {
		sprite.Draw()
		return true
	})

	rl.DrawRectangle(0, 0, 120, 60, rl.DarkGray)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d", world.Size()), 10, 30, 20, rl.Red)
}

func (s *renderController) FixedUpdate(world *ecs.World) {}

func (s *renderController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
