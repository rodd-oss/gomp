/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewRenderSystem() RenderSystem {
	return RenderSystem{}
}

type RenderSystem struct {
	EntityManager  *ecs.EntityManager
	TextureRenders *stdcomponents.TextureRenderComponentManager
	Positions      *stdcomponents.PositionComponentManager
}

func (s *RenderSystem) Init() {
	rl.InitWindow(1024, 768, "raylib [core] ebiten-ecs - basic window")
	//currentMonitorRefreshRate := int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor()))
	//rl.SetTargetFPS(12)
}
func (s *RenderSystem) Run() bool {
	if rl.WindowShouldClose() {
		return false
	}
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	s.TextureRenders.AllData(func(tr *stdcomponents.TextureRender) bool {
		texture := *tr.Texture
		rl.DrawTexturePro(texture, tr.Frame, tr.Dest, tr.Origin, tr.Rotation, tr.Tint)
		return true
	})

	// rl.DrawRectangle(0, 0, 120, 120, rl.DarkGray)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d", s.EntityManager.Size()), 10, 30, 20, rl.Red)

	rl.EndDrawing()
	return true
}

func (s *RenderSystem) Destroy() {
	rl.CloseWindow()
}
