/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- HromRu Donated 1 500 RUB

Thank you for your support!
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"time"
)

func NewMinimapRenderSystem() MinimapRenderSystem {
	return MinimapRenderSystem{}
}

type MinimapRenderSystem struct {
	EntityManager                      *ecs.EntityManager
	RlTexturePros                      *stdcomponents.RLTextureProComponentManager
	Positions                          *stdcomponents.PositionComponentManager
	Rotations                          *stdcomponents.RotationComponentManager
	Scales                             *stdcomponents.ScaleComponentManager
	AnimationPlayers                   *stdcomponents.AnimationPlayerComponentManager
	Tints                              *stdcomponents.TintComponentManager
	Flips                              *stdcomponents.FlipComponentManager
	Renderables                        *stdcomponents.RenderableComponentManager
	AnimationStates                    *stdcomponents.AnimationStateComponentManager
	Sprites                            *stdcomponents.SpriteComponentManager
	SpriteMatrixes                     *stdcomponents.SpriteMatrixComponentManager
	RenderOrders                       *stdcomponents.RenderOrderComponentManager
	BoxColliders                       *stdcomponents.BoxColliderComponentManager
	CircleColliders                    *stdcomponents.CircleColliderComponentManager
	AABBs                              *stdcomponents.AABBComponentManager
	Collisions                         *stdcomponents.CollisionComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	BvhTrees                           *stdcomponents.BvhTreeComponentManager
	MinimapCameras                     *stdcomponents.MinimapCameraComponentManager
	RenderTexture2D                    *stdcomponents.RenderTexture2DComponentManager
	renderList                         []renderEntry
	instanceData                       []stdcomponents.RLTexturePro

	monitorWidth  int
	monitorHeight int

	debug           bool
	renderTexture2D *stdcomponents.RenderTexture2D
}

func (s *MinimapRenderSystem) Init() {
	//rl.InitWindow(1280, 720, "GOMP")
	s.monitorWidth = rl.GetScreenWidth()
	s.monitorHeight = rl.GetScreenHeight()
	s.monitorWidth = rl.GetScreenWidth()
	s.monitorHeight = rl.GetScreenHeight()

	s.MinimapCameras.EachComponent(func(c *stdcomponents.MinimapCamera) bool {
		//c.Camera2D = stdcomponents.Camera2D{
		c.Target = vectors.Vec2{X: float32(s.monitorWidth / 2), Y: float32(s.monitorHeight / 2)}
		c.Offset = vectors.Vec2{X: float32(s.monitorWidth / 2), Y: float32(s.monitorHeight / 2)}
		c.Rotation = 0
		c.Zoom = 0.5
		//}
		return false
	})

	screenCamera := s.EntityManager.Create()
	s.renderTexture2D = s.RenderTexture2D.Create(screenCamera, stdcomponents.RenderTexture2D{})
	s.renderTexture2D.Texture = rl.LoadRenderTexture(int32(s.monitorWidth), int32(s.monitorHeight))
	s.renderTexture2D.Frame = rl.NewRectangle(0, 0, float32(s.monitorWidth/6), float32(s.monitorHeight/6))
	s.renderTexture2D.Position = rl.Vector2{X: 0, Y: float32(s.monitorHeight - s.monitorHeight/6)}
}

func (s *MinimapRenderSystem) Run(dt time.Duration) bool {
	if rl.IsKeyPressed(rl.KeyF12) {
		s.debug = !s.debug
	}

	var minimapCamera *stdcomponents.MinimapCamera
	s.MinimapCameras.EachComponent(func(c *stdcomponents.MinimapCamera) bool {
		minimapCamera = c
		return false
	})

	rl.BeginTextureMode(s.renderTexture2D.Texture)
	rl.BeginMode2D(minimapCamera.ToRaylibCamera())
	rl.ClearBackground(rl.White)
	s.BoxColliders.EachEntity(func(e ecs.Entity) bool {
		col := s.BoxColliders.Get(e)
		scale := s.Scales.Get(e)
		pos := s.Positions.Get(e)
		rot := s.Rotations.Get(e)

		rl.DrawRectanglePro(rl.Rectangle{
			X:      pos.XY.X,
			Y:      pos.XY.Y,
			Width:  col.WH.X * scale.XY.X,
			Height: col.WH.Y * scale.XY.Y,
		}, rl.Vector2{
			X: col.Offset.X * scale.XY.X,
			Y: col.Offset.Y * scale.XY.Y,
		}, float32(rot.Degrees()), rl.DarkGreen)
		return true
	})
	s.CircleColliders.EachEntity(func(e ecs.Entity) bool {
		col := s.CircleColliders.Get(e)
		scale := s.Scales.Get(e)
		pos := s.Positions.Get(e)

		color := rl.DarkGreen
		isSleeping := s.ColliderSleepStateComponentManager.Get(e)
		if isSleeping != nil {
			color = rl.Blue
		}

		posWithOffset := pos.XY.Add(col.Offset.Mul(scale.XY))
		rl.DrawCircle(int32(posWithOffset.X), int32(posWithOffset.Y), col.Radius*scale.XY.X, color)
		return true
	})

	rl.EndMode2D()
	rl.EndTextureMode()

	return false
}

func (s *MinimapRenderSystem) Destroy() {
}

type MinimapRenderInjector struct {
	EntityManager                      *ecs.EntityManager
	RlTexturePros                      *stdcomponents.RLTextureProComponentManager
	Positions                          *stdcomponents.PositionComponentManager
	Rotations                          *stdcomponents.RotationComponentManager
	Scales                             *stdcomponents.ScaleComponentManager
	AnimationPlayers                   *stdcomponents.AnimationPlayerComponentManager
	Tints                              *stdcomponents.TintComponentManager
	Flips                              *stdcomponents.FlipComponentManager
	Renderables                        *stdcomponents.RenderableComponentManager
	AnimationStates                    *stdcomponents.AnimationStateComponentManager
	Sprites                            *stdcomponents.SpriteComponentManager
	SpriteMatrixes                     *stdcomponents.SpriteMatrixComponentManager
	RenderOrders                       *stdcomponents.RenderOrderComponentManager
	BoxColliders                       *stdcomponents.BoxColliderComponentManager
	CircleColliders                    *stdcomponents.CircleColliderComponentManager
	AABBs                              *stdcomponents.AABBComponentManager
	Collisions                         *stdcomponents.CollisionComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	BvhTrees                           *stdcomponents.BvhTreeComponentManager
	MinimapCameras                     *stdcomponents.MinimapCameraComponentManager
	RenderTexture2D                    *stdcomponents.RenderTexture2DComponentManager
}

func (s *MinimapRenderSystem) InjectWorld(injector *MinimapRenderInjector) {
	s.EntityManager = injector.EntityManager
	s.RlTexturePros = injector.RlTexturePros
	s.Positions = injector.Positions
	s.Rotations = injector.Rotations
	s.Scales = injector.Scales
	s.AnimationPlayers = injector.AnimationPlayers
	s.Tints = injector.Tints
	s.Flips = injector.Flips
	s.Renderables = injector.Renderables
	s.AnimationStates = injector.AnimationStates
	s.Sprites = injector.Sprites
	s.SpriteMatrixes = injector.SpriteMatrixes
	s.RenderOrders = injector.RenderOrders
	s.BoxColliders = injector.BoxColliders
	s.CircleColliders = injector.CircleColliders
	s.AABBs = injector.AABBs
	s.Collisions = injector.Collisions
	s.ColliderSleepStateComponentManager = injector.ColliderSleepStateComponentManager
	s.BvhTrees = injector.BvhTrees
	s.MinimapCameras = injector.MinimapCameras
	s.RenderTexture2D = injector.RenderTexture2D
}
