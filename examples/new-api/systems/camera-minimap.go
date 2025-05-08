package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"image/color"
	"time"
)

func NewMinimapSystem() MinimapSystem {
	return MinimapSystem{}
}

type MinimapSystem struct {
	EntityManager *ecs.EntityManager
	Cameras       *stdcomponents.CameraComponentManager
	Position      *stdcomponents.PositionComponentManager
	Rotation      *stdcomponents.RotationComponentManager
	RlTexturePros *stdcomponents.RLTextureProComponentManager
	Renderables   *stdcomponents.RenderableComponentManager
	FrameBuffer2D *stdcomponents.FrameBuffer2DComponentManager
	Player        *components.PlayerTagComponentManager

	minimapCamera        ecs.Entity
	frameBufferComponent stdcomponents.FrameBuffer2D
	cameraComponent      stdcomponents.Camera
	disabled             bool
}

func (s *MinimapSystem) Init() {
	width, height := rl.GetScreenWidth(), rl.GetScreenHeight()

	s.minimapCamera = s.EntityManager.Create()
	s.Cameras.Create(s.minimapCamera, stdcomponents.Camera{
		Camera2D: rl.Camera2D{
			Target:   rl.Vector2{},
			Offset:   rl.Vector2(vectors.Vec2{X: float32(width), Y: float32(height)}.Scale(0.5)),
			Rotation: 0,
			Zoom:     .5,
		},
		Dst:       vectors.Rectangle{X: 0, Y: float32(width) - float32(height)*0.1666666666666667, Width: float32(width) * 0.1666666666666667, Height: float32(height) * 0.1666666666666667},
		Layer:     config.MinimapCameraLayer,
		Order:     1,
		Culling:   stdcomponents.Culling2DFullscreenBB,
		BlendMode: rl.BlendAlpha,
		BGColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		Tint:      color.RGBA{R: 255, G: 255, B: 255, A: 255},
	})
	s.FrameBuffer2D.Create(s.minimapCamera, stdcomponents.FrameBuffer2D{
		Position:  rl.Vector2{},
		Frame:     rl.NewRectangle(0, 0, float32(width), float32(height)),
		Texture:   rl.LoadRenderTexture(int32(width), int32(height)),
		Layer:     config.MinimapCameraLayer,
		BlendMode: rl.BlendAlpha,
		Rotation:  0,
		Tint:      rl.White,
		Dst:       rl.Rectangle{Y: float32(height) - float32(height)*0.1666666666666667, Width: float32(width) * 0.1666666666666667, Height: float32(height) * 0.1666666666666667},
	})
}

func (s *MinimapSystem) Run(dt time.Duration) bool {
	c := s.Cameras.GetUnsafe(s.minimapCamera)

	if rl.IsKeyPressed(rl.KeyM) {
		if s.disabled {
			s.disabled = false
			c = s.Cameras.Create(s.minimapCamera, s.cameraComponent)
			s.FrameBuffer2D.Create(s.minimapCamera, s.frameBufferComponent)
		} else {
			s.disabled = true
			s.cameraComponent = *c
			s.frameBufferComponent = *s.FrameBuffer2D.GetUnsafe(s.minimapCamera)
			s.Cameras.Delete(s.minimapCamera)
			s.FrameBuffer2D.Delete(s.minimapCamera)
		}
	}
	if s.disabled {
		return false
	}
	s.Player.EachEntity()(func(entity ecs.Entity) bool {
		playerPosition := s.Position.GetUnsafe(entity)
		rotation := s.Rotation.GetUnsafe(entity)

		c.Camera2D.Target.X = playerPosition.XY.X
		c.Camera2D.Target.Y = playerPosition.XY.Y
		c.Camera2D.Rotation = -float32(rotation.Degrees())
		return false
	})
	if rl.IsWindowResized() {
		width, height := rl.GetScreenWidth(), rl.GetScreenHeight()
		c := s.Cameras.GetUnsafe(s.minimapCamera)
		c.Dst = vectors.Rectangle{X: 0, Y: float32(width) - float32(height)*0.1666666666666667, Width: float32(width) * 0.1666666666666667, Height: float32(height) * 0.1666666666666667}
		c.Offset = rl.Vector2(vectors.Vec2{X: float32(width), Y: float32(height)}.Scale(0.5))
		fb := s.FrameBuffer2D.GetUnsafe(s.minimapCamera)
		rl.UnloadRenderTexture(fb.Texture)
		fb.Texture = rl.LoadRenderTexture(int32(width), int32(height))
		fb.Frame = rl.NewRectangle(0, 0, float32(width), float32(height))
		fb.Dst = rl.Rectangle{Y: float32(height) - float32(height)*0.1666666666666667, Width: float32(width) * 0.1666666666666667, Height: float32(height) * 0.1666666666666667}
	}
	return false
}

func (s *MinimapSystem) Destroy() {
}
