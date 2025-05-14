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
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"image/color"
	"math"
	"time"
)

// MainCameraSystem is simple system responsible for managing the camera entities in the game.
type MainCameraSystem struct {
	EntityManager *ecs.EntityManager
	Cameras       *stdcomponents.CameraComponentManager
	Position      *stdcomponents.PositionComponentManager
	Rotation      *stdcomponents.RotationComponentManager
	RlTexturePros *stdcomponents.RLTextureProComponentManager
	Renderables   *stdcomponents.RenderableComponentManager
	FrameBuffer2D *stdcomponents.FrameBuffer2DComponentManager
	Player        *components.PlayerTagComponentManager

	mainCamera   ecs.Entity
	shouldRotate bool
}

func NewMainCameraSystem() MainCameraSystem {
	return MainCameraSystem{}
}

func (s *MainCameraSystem) Init() {
	width, height := rl.GetScreenWidth(), rl.GetScreenHeight()

	s.mainCamera = s.EntityManager.Create()
	s.Cameras.Create(s.mainCamera, stdcomponents.Camera{
		Camera2D: rl.Camera2D{
			Target:   rl.Vector2{},
			Offset:   rl.Vector2(vectors.Vec2{X: float32(width), Y: float32(height)}.Scale(0.5)),
			Rotation: 0,
			Zoom:     1.0,
		},
		Dst:       vectors.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)},
		Layer:     config.MainCameraLayer,
		Order:     0,
		Culling:   stdcomponents.Culling2DFullscreenBB,
		BlendMode: rl.BlendAlphaPremultiply,
		BGColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		Tint:      color.RGBA{R: 255, G: 255, B: 255, A: 255},
	})
	s.FrameBuffer2D.Create(s.mainCamera, stdcomponents.FrameBuffer2D{
		Position:  rl.Vector2{},
		Frame:     rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)},
		Texture:   rl.LoadRenderTexture(int32(width), int32(height)),
		Layer:     config.MainCameraLayer,
		BlendMode: rl.BlendColor,
		Rotation:  0,
		Tint:      rl.White,
		Dst:       rl.Rectangle{Width: float32(width), Height: float32(height)},
	})
}

func (s *MainCameraSystem) Run(dt time.Duration) {
	// Update camera positions or other logic here
	// For example, you might want to move the main camera based on player input
	// or other game events.
	// This is just a placeholder for the actual camera logic.

	scroll := rl.GetMouseWheelMove()
	if scroll != 0.0 {
		c := s.Cameras.GetUnsafe(s.mainCamera)
		c.Zoom += scroll * 0.1
		if c.Zoom < 0.1 {
			c.Zoom = 0.1
		}
	}

	// Follow player for main camera and minimap camera
	if rl.IsKeyPressed(rl.KeyR) {
		s.shouldRotate = !s.shouldRotate
	}
	s.Player.EachEntity(func(entity ecs.Entity) bool {
		playerPosition := s.Position.GetUnsafe(entity)
		c := s.Cameras.GetUnsafe(s.mainCamera)
		//decay := 40.0 // DECAY IS TICKRATE DEPENDENT
		//c.Camera2D.Target.X = float32(s.expDecay(float64(c.Camera2D.Target.X), float64(playerPosition.XY.X), decay, float64(dt)))
		//c.Camera2D.Target.Y = float32(s.expDecay(float64(c.Camera2D.Target.Y), float64(playerPosition.XY.Y), decay, float64(dt)))
		c.Camera2D.Target = rl.Vector2(playerPosition.XY)
		if s.shouldRotate {
			rotation := s.Rotation.GetUnsafe(entity)
			c.Camera2D.Rotation = -float32(rotation.Degrees())
		} else {
			c.Camera2D.Rotation = 0
		}
		return false
	})

	if rl.IsWindowResized() {
		width, height := rl.GetScreenWidth(), rl.GetScreenHeight()
		c := s.Cameras.GetUnsafe(s.mainCamera)
		c.Dst = vectors.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}
		c.Camera2D.Offset = rl.Vector2(vectors.Vec2{X: float32(width), Y: float32(height)}.Scale(0.5))
		fb := s.FrameBuffer2D.GetUnsafe(s.mainCamera)
		rl.UnloadRenderTexture(fb.Texture)
		fb.Texture = rl.LoadRenderTexture(int32(width), int32(height))
		fb.Frame = rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}
		fb.Dst = rl.Rectangle{Width: float32(width), Height: float32(height)}
	}
}

// TODO: check and do better
func (s *MainCameraSystem) expDecay(a, b, decay, dt float64) float64 {
	return b + (a-b)*(math.Exp(-decay*dt))
}

func (s *MainCameraSystem) Destroy() {
	s.EntityManager.Delete(s.mainCamera)
}
