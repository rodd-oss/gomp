/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"cmp"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/negrel/assert"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"math"
	"runtime"
	"slices"
	"time"
)

func NewRender2DCamerasSystem() Render2DCamerasSystem {
	return Render2DCamerasSystem{}
}

type Render2DCamerasSystem struct {
	Renderables         *stdcomponents.RenderableComponentManager
	RenderVisibles      *stdcomponents.RenderVisibleComponentManager
	RenderOrders        *stdcomponents.RenderOrderComponentManager
	Textures            *stdcomponents.RLTextureProComponentManager
	Tints               *stdcomponents.TintComponentManager
	AABBs               *stdcomponents.AABBComponentManager
	Cameras             *stdcomponents.CameraComponentManager
	RenderTexture2D     *stdcomponents.FrameBuffer2DComponentManager
	AnimationPlayers    *stdcomponents.AnimationPlayerComponentManager
	Flips               *stdcomponents.FlipComponentManager
	Positions           *stdcomponents.PositionComponentManager
	Scales              *stdcomponents.ScaleComponentManager
	Rotations           *stdcomponents.RotationComponentManager
	renderObjects       []renderObject
	renderObjectsSorted []renderObjectSorted
	numWorkers          int

	Engine *core.Engine
}

type renderObject struct {
	texture stdcomponents.RLTexturePro
	mask    stdcomponents.CameraLayer
	order   float32
}

type renderObjectSorted struct {
	entity ecs.Entity
	order  float32
}

func (s *Render2DCamerasSystem) Init() {
	s.renderObjects = make([]renderObject, 0, s.RenderVisibles.Len())
	s.numWorkers = runtime.NumCPU() - 2
}

func (s *Render2DCamerasSystem) Run(dt time.Duration) {
	s.prepareRender(dt)

	// Collect and sort render objects
	if cap(s.renderObjects) < s.RenderVisibles.Len() {
		s.renderObjects = make([]renderObject, 0, max(s.RenderVisibles.Len(), cap(s.renderObjects)*2))
	}

	if cap(s.renderObjectsSorted) < s.RenderVisibles.Len() {
		s.renderObjectsSorted = make([]renderObjectSorted, 0, max(s.RenderVisibles.Len(), cap(s.renderObjects)*2))
	}

	s.RenderVisibles.EachEntity(func(entity ecs.Entity) bool {
		o := s.RenderOrders.GetUnsafe(entity)
		assert.NotNil(o)

		s.renderObjectsSorted = append(s.renderObjectsSorted, renderObjectSorted{
			entity: entity,
			order:  o.CalculatedZ,
		})

		return true
	})

	slices.SortFunc(s.renderObjectsSorted, func(a, b renderObjectSorted) int {
		return cmp.Compare(a.order, b.order)
	})

	for i := range s.renderObjectsSorted {
		obj := &s.renderObjectsSorted[i]

		t := s.Textures.GetUnsafe(obj.entity)
		assert.NotNil(t)

		//TODO: rework this with future new assets manager
		if t.Texture == nil {
			continue
		}

		r := s.Renderables.GetUnsafe(obj.entity)
		assert.NotNil(r)

		s.renderObjects = append(s.renderObjects, renderObject{
			texture: *t,
			mask:    r.CameraMask,
			order:   obj.order,
		})
	}

	s.Cameras.EachEntity(func(cameraEntity ecs.Entity) bool {
		camera := s.Cameras.GetUnsafe(cameraEntity)
		assert.NotNil(camera)
		renderTexture := s.RenderTexture2D.GetUnsafe(cameraEntity)
		assert.NotNil(renderTexture)

		// Draw render objects
		rl.BeginTextureMode(renderTexture.Texture)
		rl.BeginMode2D(camera.Camera2D)
		rl.ClearBackground(camera.BGColor)

		for i := range s.renderObjects {
			obj := &s.renderObjects[i]
			if camera.Layer&obj.mask == 0 {
				continue
			}
			rl.DrawTexturePro(*obj.texture.Texture, obj.texture.Frame, obj.texture.Dest, obj.texture.Origin, obj.texture.Rotation, obj.texture.Tint)
		}

		rl.EndMode2D()
		rl.EndTextureMode()

		return true
	})

	s.renderObjects = s.renderObjects[:0]
	s.renderObjectsSorted = s.renderObjectsSorted[:0]
}

func (s *Render2DCamerasSystem) Destroy() {
}

func (s *Render2DCamerasSystem) prepareRender(dt time.Duration) {
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		assert.NotNil(texturePro)

		animation := s.AnimationPlayers.GetUnsafe(entity)
		if animation != nil {
			frame := &texturePro.Frame
			if animation.Vertical {
				frame.Y += frame.Height * float32(animation.Current)
			} else {
				frame.X += frame.Width * float32(animation.Current)
			}
		}

		flipped := s.Flips.GetUnsafe(entity)
		if flipped != nil {
			if flipped.X {
				texturePro.Frame.Width *= -1
			}
			if flipped.Y {
				texturePro.Frame.Height *= -1
			}
		}

		position := s.Positions.GetUnsafe(entity)
		if position != nil {
			//decay := 40.0 // DECAY IS TICKRATE DEPENDENT
			//x := float32(s.expDecay(float64(texturePro.Dest.X), float64(position.XY.X), decay, dts))
			//y := float32(s.expDecay(float64(texturePro.Dest.Y), float64(position.XY.Y), decay, dts))
			texturePro.Dest.X = position.XY.X
			texturePro.Dest.Y = position.XY.Y
		}

		rotation := s.Rotations.GetUnsafe(entity)
		if rotation != nil {
			texturePro.Rotation = float32(rotation.Degrees())
		}

		scale := s.Scales.GetUnsafe(entity)
		if scale != nil {
			texturePro.Dest.Width *= scale.XY.X
			texturePro.Dest.Height *= scale.XY.Y
		}

		tint := s.Tints.GetUnsafe(entity)
		if tint != nil {
			trTint := &texturePro.Tint
			trTint.A = tint.A
			trTint.R = tint.R
			trTint.G = tint.G
			trTint.B = tint.B
		}
	})
}

func (s *Render2DCamerasSystem) prepareAnimations() {
	// TODO revert loop to process animations but not a textures?
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		animation := s.AnimationPlayers.GetUnsafe(entity)
		if animation == nil {
			return
		}
		frame := &texturePro.Frame
		if animation.Vertical {
			frame.Y += frame.Height * float32(animation.Current)
		} else {
			frame.X += frame.Width * float32(animation.Current)
		}
	})
}

func (s *Render2DCamerasSystem) prepareFlips() {
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		flipped := s.Flips.GetUnsafe(entity)
		if flipped == nil {
			return
		}
		if flipped.X {
			texturePro.Frame.Width *= -1
		}
		if flipped.Y {
			texturePro.Frame.Height *= -1
		}
	})
}

func (s *Render2DCamerasSystem) preparePositions(dt time.Duration) {
	//dts := dt.Seconds()
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		position := s.Positions.GetUnsafe(entity)
		if position == nil {
			return
		}
		//decay := 40.0 // DECAY IS TICKRATE DEPENDENT
		//x := float32(s.expDecay(float64(texturePro.Dest.X), float64(position.XY.X), decay, dts))
		//y := float32(s.expDecay(float64(texturePro.Dest.Y), float64(position.XY.Y), decay, dts))
		texturePro.Dest.X = position.XY.X
		texturePro.Dest.Y = position.XY.Y
	})
}

func (s *Render2DCamerasSystem) prepareRotations() {
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		rotation := s.Rotations.GetUnsafe(entity)
		if rotation == nil {
			return
		}
		texturePro.Rotation = float32(rotation.Degrees())
	})
}

func (s *Render2DCamerasSystem) prepareScales() {
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		texturePro := s.Textures.GetUnsafe(entity)
		scale := s.Scales.GetUnsafe(entity)
		if scale == nil {
			return
		}
		texturePro.Dest.Width *= scale.XY.X
		texturePro.Dest.Height *= scale.XY.Y
	})
}

func (s *Render2DCamerasSystem) prepareTints() {
	s.Textures.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		tr := s.Textures.GetUnsafe(entity)
		tint := s.Tints.GetUnsafe(entity)
		if tint == nil {
			return
		}
		trTint := &tr.Tint
		trTint.A = tint.A
		trTint.R = tint.R
		trTint.G = tint.G
		trTint.B = tint.B
	})
}

func (s *Render2DCamerasSystem) expDecay(a, b, decay, dt float64) float64 {
	return b + (a-b)*(math.Exp(-decay*dt))
}
