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
	t             rl.Texture2D
}

func (s *renderController) Init(world *ecs.World) {
	rl.InitWindow(s.width, s.height, "raylib [core] example - basic window")

	// currentMonitorRefreshRate := rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())
	// // rl.SetTargetFPS(int32(currentMonitorRefreshRate))

	s.t = rl.LoadTexture("assets/star.png")
}

func (s *renderController) Update(world *ecs.World) {
	spriteManager := components.SpriteService.GetManager(world)

	if rl.WindowShouldClose() {
		world.SetShouldDestroy(true)
		return
	}

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	spriteManager.AllData(s.drawSprite)

	rl.DrawRectangle(0, 0, 120, 60, rl.DarkGray)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d", world.Size()), 10, 30, 20, rl.Red)
}

func (s *renderController) FixedUpdate(world *ecs.World) {}

func (s *renderController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}

func (s *renderController) drawSprite(sprite *components.Sprite) bool {
	rl.DrawTextureV(s.t, sprite.Pos, sprite.Tint)
	return true
}
