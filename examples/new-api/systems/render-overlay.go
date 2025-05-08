/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package systems

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/negrel/assert"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"image/color"
	"slices"
	"time"
)

const (
	fpsAvgSamples   = 100
	fontSize        = 20
	fpsGraphWidth   = 160
	fpsGraphHeight  = 60
	msGraphMaxValue = 33.33 // Max ms to show on graph (30 FPS)
)

func NewRenderOverlaySystem() RenderOverlaySystem {
	return RenderOverlaySystem{}
}

type RenderOverlaySystem struct {
	EntityManager                      *ecs.EntityManager
	SceneManager                       *components.AsteroidSceneManagerComponentManager
	Cameras                            *stdcomponents.CameraComponentManager
	FrameBuffer2D                      *stdcomponents.FrameBuffer2DComponentManager
	CollisionChunks                    *stdcomponents.CollisionChunkComponentManager
	CollisionCells                     *stdcomponents.CollisionCellComponentManager
	Tints                              *stdcomponents.TintComponentManager
	BvhTrees                           *stdcomponents.BvhTreeComponentManager
	Positions                          *stdcomponents.PositionComponentManager
	AABBs                              *stdcomponents.AABBComponentManager
	Collisions                         *stdcomponents.CollisionComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	Textures                           *stdcomponents.RLTextureProComponentManager
	frameBuffer                        ecs.Entity
	monitorWidth                       int
	monitorHeight                      int
	debugLvl                           int
	debug                              bool
	lastFPSTime                        time.Time
	frameCount                         int
	currentFPS                         int
	fpsSamples                         []int
	fpsSampleSum                       int
	fpsSampleIdx                       int
	avgFPS                             float64
	lowestFps                          int
	lastFrameDuration                  time.Duration
	msHistory                          []float64 // ms per frame, ring buffer
	msHistoryIdx                       int
}

func (s *RenderOverlaySystem) Init() {
	s.monitorWidth = rl.GetScreenWidth()
	s.monitorHeight = rl.GetScreenHeight()

	s.frameBuffer = s.EntityManager.Create()
	s.FrameBuffer2D.Create(s.frameBuffer, stdcomponents.FrameBuffer2D{
		Frame:     rl.Rectangle{X: 0, Y: 0, Width: float32(s.monitorWidth), Height: float32(s.monitorHeight)},
		Texture:   rl.LoadRenderTexture(int32(s.monitorWidth), int32(s.monitorHeight)),
		Layer:     config.DebugLayer,
		BlendMode: rl.BlendAlphaPremultiply,
		Tint:      rl.White,
		Dst:       rl.Rectangle{Width: float32(s.monitorWidth), Height: float32(s.monitorHeight)},
	})

	// Initialize ms history buffer for graph
	s.msHistory = make([]float64, fpsGraphWidth)
}

func (s *RenderOverlaySystem) Run(dt time.Duration) bool {
	if rl.IsKeyPressed(rl.KeyF6) {
		if !s.debug {
			s.debug = true
		} else {
			s.debug = false
		}
	}
	if rl.IsKeyPressed(rl.KeyF7) {
		s.debugLvl--
		if s.debugLvl < 0 {
			s.debugLvl = 63
		}
	}
	if rl.IsKeyPressed(rl.KeyF8) {
		s.debugLvl++
		if s.debugLvl > 63 {
			s.debugLvl = 0
		}
	}

	// FPS calculation (custom)
	now := time.Now()
	if s.lastFPSTime.IsZero() {
		s.lastFPSTime = now
		s.fpsSamples = make([]int, fpsAvgSamples)
		s.msHistory = make([]float64, fpsGraphWidth)
	}
	s.frameCount++
	s.lastFrameDuration = dt

	{ // Store current frame FPS in samples
		frameFPS := 0
		if dt > 0 {
			// Correct calculation: convert duration to frames per second
			frameFPS = int(time.Second / dt)
		}
		s.fpsSampleSum -= s.fpsSamples[s.fpsSampleIdx]
		s.fpsSamples[s.fpsSampleIdx] = frameFPS
		s.fpsSampleSum += frameFPS
		s.fpsSampleIdx = (s.fpsSampleIdx + 1) % len(s.fpsSamples)

		// Calculate average FPS over samples
		s.avgFPS = float64(s.fpsSampleSum) / float64(len(s.fpsSamples))

		// Calculate 1% FPS (lowest 1% frame in the sample window)
		s.lowestFps = slices.Min(s.fpsSamples)

		// Update frame time history (ms) on every frame
		// Use average of last two frames for smoother graph
		var ms float64
		if s.lastFrameDuration > 0 {
			ms = float64(s.lastFrameDuration.Microseconds()) / 1000.0
			s.msHistory[s.msHistoryIdx] = ms
			s.msHistoryIdx = (s.msHistoryIdx + 1) % len(s.msHistory)
		}

		if now.Sub(s.lastFPSTime) >= time.Second {
			s.currentFPS = s.frameCount
			s.frameCount = 0
			s.lastFPSTime = now
		}
	}

	s.Cameras.EachEntity()(func(entity ecs.Entity) bool {
		camera := s.Cameras.GetUnsafe(entity)
		fb := s.FrameBuffer2D.GetUnsafe(entity)
		switch fb.Layer {
		case config.MainCameraLayer:
			overlayFrame := s.FrameBuffer2D.GetUnsafe(s.frameBuffer)
			rl.BeginTextureMode(overlayFrame.Texture)
			rl.ClearBackground(rl.Blank)

			// Debug mode: BVH tree and dots
			if s.debug {
				rl.BeginMode2D(camera.Camera2D)
				cameraRect := camera.Rect()

				s.CollisionCells.EachEntity()(func(e ecs.Entity) bool {
					cell := s.CollisionCells.GetUnsafe(e)
					assert.NotNil(cell)

					if cell.Layer != stdcomponents.CollisionLayer(s.debugLvl) {
						return true
					}

					tint := s.Tints.GetUnsafe(e)
					assert.NotNil(tint)

					position := s.Positions.GetUnsafe(e)
					assert.NotNil(position)

					clr := color.RGBA{
						R: tint.R,
						G: tint.G,
						B: tint.B,
						A: 255,
					}

					// Simple AABB culling
					if s.intersects(cameraRect, vectors.Rectangle{
						X:      position.XY.X,
						Y:      position.XY.Y,
						Width:  cell.Size,
						Height: cell.Size,
					}) {
						rl.DrawRectangleLines(int32(position.XY.X), int32(position.XY.Y), int32(cell.Size), int32(cell.Size), clr)
					}
					return true
				})
				s.CollisionChunks.EachEntity()(func(e ecs.Entity) bool {
					chunk := s.CollisionChunks.GetUnsafe(e)
					assert.NotNil(chunk)

					if chunk.Layer != stdcomponents.CollisionLayer(s.debugLvl) {
						return true
					}

					tint := s.Tints.GetUnsafe(e)
					assert.NotNil(tint)

					position := s.Positions.GetUnsafe(e)
					assert.NotNil(position)

					tree := s.BvhTrees.GetUnsafe(e)
					assert.NotNil(tree)

					tree.AabbNodes.EachData()(func(a *stdcomponents.AABB) bool {
						// Simple AABB culling
						if s.intersects(cameraRect, a.Rect()) {
							rl.DrawRectangleRec(rl.Rectangle{
								X:      a.Min.X,
								Y:      a.Min.Y,
								Width:  a.Max.X - a.Min.X,
								Height: a.Max.Y - a.Min.Y,
							}, *tint)
						}
						return true
					})

					clr := color.RGBA{
						R: tint.R,
						G: tint.G,
						B: tint.B,
						A: 255,
					}

					// Simple AABB culling
					if s.intersects(cameraRect, vectors.Rectangle{
						X:      position.XY.X,
						Y:      position.XY.Y,
						Width:  chunk.Size,
						Height: chunk.Size,
					}) {
						rl.DrawRectangleLines(int32(position.XY.X), int32(position.XY.Y), int32(chunk.Size), int32(chunk.Size), clr)
					}
					return true
				})
				s.AABBs.EachEntity()(func(e ecs.Entity) bool {
					aabb := s.AABBs.GetUnsafe(e)
					clr := rl.Green
					isSleeping := s.ColliderSleepStateComponentManager.GetUnsafe(e)
					if isSleeping != nil {
						clr = rl.Blue
					}
					if s.intersects(cameraRect, aabb.Rect()) {
						rl.DrawRectangleLinesEx(rl.Rectangle{
							X:      aabb.Min.X,
							Y:      aabb.Min.Y,
							Width:  aabb.Max.X - aabb.Min.X,
							Height: aabb.Max.Y - aabb.Min.Y,
						}, 1, clr)
					}
					return true
				})
				s.Collisions.EachEntity()(func(entity ecs.Entity) bool {
					pos := s.Positions.GetUnsafe(entity)
					rec := vectors.Rectangle{
						X:      pos.XY.X - 8,
						Y:      pos.XY.Y - 8,
						Width:  16,
						Height: 16,
					}
					if s.intersects(cameraRect, rec) {
						rl.DrawRectangleRec(rl.Rectangle(rec), rl.Red)
					}
					return true
				})
				s.Textures.EachComponent()(func(r *stdcomponents.RLTexturePro) bool {
					rec := vectors.Rectangle{
						X:      r.Dest.X - 2,
						Y:      r.Dest.Y - 2,
						Width:  4,
						Height: 4,
					}
					if s.intersects(cameraRect, rec) {
						rl.DrawRectanglePro(rl.Rectangle(rec), rl.Vector2{}, r.Rotation, rl.Red)
					}
					return true
				})

				rl.EndMode2D()
			}

			// Print stats
			const x = 10
			const y = 10
			statsPanelWidth := float32(200 + fpsGraphWidth)
			statsPanelHeight := float32(y + fontSize*8)
			rl.DrawRectangleRec(rl.Rectangle{Height: statsPanelHeight, Width: statsPanelWidth}, rl.Black)
			s.drawCustomFPS(x, y)
			rl.DrawText(fmt.Sprintf("%d entities", s.EntityManager.Size()), x, y+fontSize*6, fontSize, rl.RayWhite)
			rl.DrawText(fmt.Sprintf("%d debugLvl", s.debugLvl), x, y+fontSize*7, 20, rl.RayWhite)
			// Game over
			s.SceneManager.EachComponent()(func(a *components.AsteroidSceneManager) bool {
				rl.DrawText(fmt.Sprintf("Player HP: %d", a.PlayerHp), x, y+fontSize*4, 20, rl.RayWhite)
				rl.DrawText(fmt.Sprintf("Score: %d", a.PlayerScore), x, y+fontSize*5, 20, rl.RayWhite)
				if a.PlayerHp <= 0 {
					text := "Game Over"
					textSize := rl.MeasureTextEx(rl.GetFontDefault(), text, 96, 0)
					x := (s.monitorWidth - int(textSize.X)) / 2
					y := (s.monitorHeight - int(textSize.Y)) / 2
					rl.DrawText(text, int32(x), int32(y), 96, rl.Red)

				}
				return false
			})
			rl.EndTextureMode()

		case config.MinimapCameraLayer:
			rl.BeginTextureMode(fb.Texture)
			rl.DrawRectangleLines(2, 2, fb.Texture.Texture.Width-2, fb.Texture.Texture.Height-2, rl.Green)
			rl.EndTextureMode()
		}

		return true
	})

	if rl.IsWindowResized() {
		s.monitorWidth = rl.GetScreenWidth()
		s.monitorHeight = rl.GetScreenHeight()
		fb := s.FrameBuffer2D.GetUnsafe(s.frameBuffer)
		rl.UnloadRenderTexture(fb.Texture)
		fb.Texture = rl.LoadRenderTexture(int32(s.monitorWidth), int32(s.monitorHeight))
		fb.Frame = rl.Rectangle{X: 0, Y: 0, Width: float32(s.monitorWidth), Height: float32(s.monitorHeight)}
		fb.Dst = rl.Rectangle{Width: float32(s.monitorWidth), Height: float32(s.monitorHeight)}
	}
	return true
}

// Draws FPS stats: low, current frame, average, and current FPS
func (s *RenderOverlaySystem) drawCustomFPS(x, y int32) {
	fps := int32(s.currentFPS)

	// Frame time in milliseconds
	frameTimeMs := 0.0
	if s.lastFrameDuration > 0 {
		frameTimeMs = float64(s.lastFrameDuration.Microseconds()) / 1000.0
	}

	avgFPS := int32(s.avgFPS)

	// Colors
	fontColor := rl.Lime
	if fps < 30 {
		fontColor = rl.Red
	} else if fps < 60 {
		fontColor = rl.Yellow
	}

	// Frame time color (lower is better)
	frameTimeColor := rl.Lime
	if frameTimeMs > 33.33 { // 30 FPS threshold (33.33ms)
		frameTimeColor = rl.Red
	} else if frameTimeMs > 16.67 { // 60 FPS threshold (16.67ms)
		frameTimeColor = rl.Yellow
	}

	// Draw all stats
	rl.DrawText(fmt.Sprintf("FPS: %d", fps), x, y, fontSize, fontColor)
	rl.DrawText(fmt.Sprintf("Frame: %.2f ms", frameTimeMs), x, y+fontSize, fontSize, frameTimeColor)
	rl.DrawText(fmt.Sprintf("Avg %d: %d", fpsAvgSamples, avgFPS), x, y+fontSize*2, fontSize, fontColor)
	rl.DrawText(fmt.Sprintf("Low: %d", s.lowestFps), x, y+fontSize*3, fontSize, fontColor)

	// Draw ms graph
	s.drawMsGraph(x+180, y)
}

// Draw a graph of historical frame times in milliseconds
func (s *RenderOverlaySystem) drawMsGraph(x, y int32) {
	// Draw graph border
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(fpsGraphWidth),
		Height: float32(fpsGraphHeight),
	}, 1, rl.Gray)

	// Draw graph background
	rl.DrawRectangle(x+1, y+1, fpsGraphWidth-2, fpsGraphHeight-2, rl.Black)

	// Draw horizontal reference lines (33.33ms, 16.67ms, 8.33ms)
	// These correspond to 30 FPS, 60 FPS, and 120 FPS
	refLines := []struct {
		ms    float32
		color rl.Color
		label string
	}{
		{33.33, rl.Red, "33.33 (30 FPS)"},
		{16.67, rl.Yellow, "16.67 (60 FPS)"},
		{8.33, rl.Green, "8.33 (120 FPS)"},
	}
	for _, ref := range refLines {
		refY := y + int32(float32(fpsGraphHeight)*(ref.ms/msGraphMaxValue))
		rl.DrawLineEx(
			rl.NewVector2(float32(x), float32(refY)),
			rl.NewVector2(float32(x+int32(fpsGraphWidth)), float32(refY)),
			1.0,
			ref.color,
		)
		rl.DrawText(ref.label, x+int32(fpsGraphWidth)+2, refY-8, 10, rl.Fade(ref.color, 0.8))
	}

	// Start from the oldest sample and move forward
	startIdx := s.msHistoryIdx % len(s.msHistory)

	// Draw frame time data points and connect with lines
	for i := 0; i < len(s.msHistory)-1; i++ {
		// Calculate indices in a way that we're drawing from left to right,
		// with the newest data on the right
		idx := (startIdx + i) % len(s.msHistory)
		nextIdx := (startIdx + i + 1) % len(s.msHistory)

		ms1 := float32(s.msHistory[idx])
		ms2 := float32(s.msHistory[nextIdx])

		// Clamp values to max
		if ms1 > msGraphMaxValue {
			ms1 = msGraphMaxValue
		}
		if ms2 > msGraphMaxValue {
			ms2 = msGraphMaxValue
		}

		// Calculate positions (note: for ms, higher value = worse performance, so we scale directly)
		x1 := x + int32(i)
		y1 := y + int32(float32(fpsGraphHeight)*(ms1/msGraphMaxValue))
		x2 := x + int32(i+1)
		y2 := y + int32(float32(fpsGraphHeight)*(ms2/msGraphMaxValue))

		// Choose color based on frame time
		lineColor := rl.Green
		if ms2 > 16.67 { // 60 FPS threshold
			lineColor = rl.Yellow
		}
		if ms2 > 33.33 { // 30 FPS threshold
			lineColor = rl.Red
		}

		// Skip drawing if either value is zero (not yet initialized)
		if ms1 > 0 && ms2 > 0 {
			rl.DrawLineEx(
				rl.NewVector2(float32(x1), float32(y1)),
				rl.NewVector2(float32(x2), float32(y2)),
				2.0,
				lineColor,
			)
		}
	}

	// Draw a vertical line indicating the current position in the buffer
	currentX := x + int32(len(s.msHistory)-1)
	rl.DrawLineEx(
		rl.NewVector2(float32(currentX), float32(y)),
		rl.NewVector2(float32(currentX), float32(y+int32(fpsGraphHeight))),
		1.0,
		rl.White,
	)
}

func (s *RenderOverlaySystem) intersects(rect1, rect2 vectors.Rectangle) bool {
	return rect1.X < rect2.X+rect2.Width &&
		rect1.X+rect1.Width > rect2.X &&
		rect1.Y < rect2.Y+rect2.Height &&
		rect1.Y+rect1.Height > rect2.Y
}

func (s *RenderOverlaySystem) Destroy() {}
