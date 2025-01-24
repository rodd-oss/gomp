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

type renderController struct{}

func (s *renderController) Init(world *ecs.World) {
	currentMonitorRefreshRate := int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor()))
	rl.SetTargetFPS(currentMonitorRefreshRate)
}
func (s *renderController) Update(world *ecs.World) {
	spriteRenders := components.TextureRenderService.GetManager(world)

	if rl.WindowShouldClose() {
		world.SetShouldDestroy(true)
		return
	}

	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	spriteRenders.AllData(func(tr *components.TextureRender) bool {
		texture := *tr.Texture
		rl.DrawTexturePro(texture, tr.Frame, tr.Dest, tr.Origin, tr.Rotation, tr.Tint)
		return true
	})

	// rl.DrawRectangle(0, 0, 120, 120, rl.DarkGray)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d", world.Size()), 10, 30, 20, rl.Red)
	rl.DrawText(fmt.Sprintf("%s", world.DtUpdate()), 10, 50, 20, rl.Red)
	rl.DrawText(fmt.Sprintf("%s", world.DtFixedUpdate()), 10, 70, 20, rl.Red)

	rl.EndDrawing()
}

func (s *renderController) FixedUpdate(world *ecs.World) {}
func (s *renderController) Destroy(world *ecs.World)     {}
