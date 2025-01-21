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
	spriteRenders := components.SpriteRenderService.GetManager(world)
	sprites := components.SpriteService.GetManager(world)

	sprites.AllDataParallel(func(sprite *components.Sprite) bool {
		if sprite.Texture == nil {
			sprite.Texture = &s.texture
		}

		sprite.TextureRegion = rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(s.texture.Width),
			Height: float32(s.texture.Height),
		}
		return true
	})

	if rl.WindowShouldClose() {
		world.SetShouldDestroy(true)
		return
	}

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	spriteRenders.AllData(func(spriteRender *components.SpriteRender) bool {
		sprite := &spriteRender.Sprite
		dest := spriteRender.Dest
		texture := *sprite.Texture

		rl.DrawTexturePro(texture, sprite.TextureRegion, dest, sprite.Origin, spriteRender.Rotation, sprite.Tint)
		return true
	})

	// rl.DrawRectangle(0, 0, 120, 120, rl.DarkGray)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d", world.Size()), 10, 30, 20, rl.Red)
	rl.DrawText(fmt.Sprintf("%s", world.DtUpdate()), 10, 50, 20, rl.Red)
	rl.DrawText(fmt.Sprintf("%s", world.DtFixedUpdate()), 10, 70, 20, rl.Red)
}

func (s *renderController) FixedUpdate(world *ecs.World) {}

func (s *renderController) Destroy(world *ecs.World) {
	rl.CloseWindow()
}
