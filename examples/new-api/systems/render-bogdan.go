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
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"math"
	"slices"
	"sync"
	"time"
)

func NewRenderBogdanSystem() RenderBogdanSystem {
	return RenderBogdanSystem{}
}

type RenderBogdanSystem struct {
	EntityManager    *ecs.EntityManager
	RlTexturePros    *stdcomponents.RLTextureProComponentManager
	Positions        *stdcomponents.PositionComponentManager
	Rotations        *stdcomponents.RotationComponentManager
	Scales           *stdcomponents.ScaleComponentManager
	AnimationPlayers *stdcomponents.AnimationPlayerComponentManager
	Tints            *stdcomponents.TintComponentManager
	Flips            *stdcomponents.FlipComponentManager
	Renderables      *stdcomponents.RenderableComponentManager
	AnimationStates  *stdcomponents.AnimationStateComponentManager
	SpriteMatrixes   *stdcomponents.SpriteMatrixComponentManager
	RenderOrders     *stdcomponents.RenderOrderComponentManager
	ColliderBoxes    *stdcomponents.BoxColliderComponentManager
	Collisions       *stdcomponents.CollisionComponentManager
	renderList       []renderEntry
	instanceData     []stdcomponents.RLTexturePro
	camera           rl.Camera2D
}

type renderEntry struct {
	Entity    ecs.Entity
	TextureId int
	ZIndex    float32
}

func (s *RenderBogdanSystem) Init() {
	s.camera = rl.Camera2D{
		Target:   rl.NewVector2(0, 0),
		Offset:   rl.NewVector2(0, 0),
		Rotation: 0,
		Zoom:     1,
	}
}
func (s *RenderBogdanSystem) Run(dt time.Duration) bool {
	if rl.WindowShouldClose() {
		return false
	}

	s.prepareRender(dt)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	// draw grid
	const gridSize = 256
	for i := int32(1); i < 1024/gridSize; i++ {
		rl.DrawLine(i*gridSize, 0, i*gridSize, 768, rl.Green)
	}
	for i := int32(1); i < 768/gridSize; i++ {
		rl.DrawLine(0, i*gridSize, 1024, i*gridSize, rl.Green)
	}
	s.render()
	s.ColliderBoxes.EachEntity(func(e ecs.Entity) bool {
		box := s.ColliderBoxes.GetUnsafe(e)
		pos := s.Positions.GetUnsafe(e)

		rl.DrawRectangleLines(int32(pos.XY.X), int32(pos.XY.Y), int32(box.WH.X), int32(box.WH.Y), rl.Red)
		return true
	})
	s.Collisions.EachEntity(func(entity ecs.Entity) bool {
		pos := s.Positions.GetUnsafe(entity)
		rl.DrawRectangle(int32(pos.XY.X), int32(pos.XY.X), 16, 16, rl.Red)
		return true
	})
	rl.DrawRectangle(0, 0, 200, 60, rl.DarkBrown)
	rl.DrawFPS(10, 10)
	rl.DrawText(fmt.Sprintf("%d entities", s.EntityManager.Size()), 10, 30, 20, rl.RayWhite)
	rl.EndDrawing()

	return true
}

func (s *RenderBogdanSystem) Destroy() {}

func (s *RenderBogdanSystem) render() {
	// Extract and sort entities
	if cap(s.renderList) < s.Renderables.Len() {
		s.renderList = append(s.renderList, make([]renderEntry, 0, s.Renderables.Len()-cap(s.renderList))...)
	}
	s.Renderables.EachEntity(func(e ecs.Entity) bool {
		sprite := s.SpriteMatrixes.Get(e)
		renderOrder := s.RenderOrders.GetUnsafe(e)
		s.renderList = append(s.renderList, renderEntry{
			Entity:    e,
			TextureId: int(sprite.Texture.ID),
			ZIndex:    renderOrder.CalculatedZ,
		})
		return true
	})

	slices.SortStableFunc(s.renderList, func(a, b renderEntry) int {
		if a.TextureId == b.TextureId {
			return int(math.Floor(float64(a.ZIndex - b.ZIndex)))
		}
		return int(a.TextureId - b.TextureId)
	})

	// Batch and render
	var currentTex = -1
	var instanceData []stdcomponents.RLTexturePro = make([]stdcomponents.RLTexturePro, 0, 8192)
	for i := range s.renderList {
		entry := &s.renderList[i]
		if entry.TextureId != currentTex || len(instanceData) >= 8192 {
			if len(instanceData) > 0 {
				s.submitBatch(currentTex, instanceData)
				instanceData = instanceData[:0]
			}
			currentTex = entry.TextureId
		}
		instanceData = append(instanceData, s.getInstanceData(entry.Entity))
	}
	s.submitBatch(currentTex, instanceData) // Submit last batch
	s.renderList = s.renderList[:0]
}

func (s *RenderBogdanSystem) getInstanceData(e ecs.Entity) stdcomponents.RLTexturePro {
	return *s.RlTexturePros.GetUnsafe(e)
}

func (s *RenderBogdanSystem) prepareRender(dt time.Duration) {
	wg := new(sync.WaitGroup)
	wg.Add(6)
	s.prepareAnimations(wg)
	go s.prepareFlips(wg)
	go s.preparePositions(wg, dt)
	go s.prepareRotations(wg)
	go s.prepareScales(wg)
	go s.prepareTints(wg)
	wg.Wait()
}

func (s *RenderBogdanSystem) prepareAnimations(wg *sync.WaitGroup) {
	defer wg.Done()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		texturePro := s.RlTexturePros.GetUnsafe(entity)
		animation := s.AnimationPlayers.GetUnsafe(entity)
		if animation == nil {
			return true
		}
		frame := &texturePro.Frame
		if animation.Vertical {
			frame.Y += frame.Height * float32(animation.Current)
		} else {
			frame.X += frame.Width * float32(animation.Current)
		}
		return true
	})
}

func (s *RenderBogdanSystem) prepareFlips(wg *sync.WaitGroup) {
	defer wg.Done()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		texturePro := s.RlTexturePros.GetUnsafe(entity)
		mirrored := s.Flips.GetUnsafe(entity)
		if mirrored == nil {
			return true
		}
		if mirrored.X {
			texturePro.Frame.Width *= -1
		}
		if mirrored.Y {
			texturePro.Frame.Height *= -1
		}
		return true
	})
}

func (s *RenderBogdanSystem) preparePositions(wg *sync.WaitGroup, dt time.Duration) {
	defer wg.Done()
	//dts := dt.Seconds()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		texturePro := s.RlTexturePros.GetUnsafe(entity)
		position := s.Positions.GetUnsafe(entity)
		if position == nil {
			return true
		}
		//decay := 16.0 // DECAY IS TICKRATE DEPENDENT
		//texturePro.Dest.X = float32(s.expDecay(float64(texturePro.Dest.X), float64(position.X), decay, dts))
		//texturePro.Dest.Y = float32(s.expDecay(float64(texturePro.Dest.Y), float64(position.Y), decay, dts))
		texturePro.Dest.X = position.XY.X
		texturePro.Dest.Y = position.XY.Y

		return true
	})
}

func (s *RenderBogdanSystem) prepareRotations(wg *sync.WaitGroup) {
	defer wg.Done()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		texturePro := s.RlTexturePros.GetUnsafe(entity)
		rotation := s.Rotations.GetUnsafe(entity)
		if rotation == nil {
			return true
		}
		texturePro.Rotation = float32(rotation.Angle)
		return true
	})
}

func (s *RenderBogdanSystem) prepareScales(wg *sync.WaitGroup) {
	defer wg.Done()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		texturePro := s.RlTexturePros.GetUnsafe(entity)
		scale := s.Scales.GetUnsafe(entity)
		if scale == nil {
			return true
		}
		texturePro.Dest.Width *= scale.XY.X
		texturePro.Dest.Height *= scale.XY.Y
		return true
	})
}

func (s *RenderBogdanSystem) prepareTints(wg *sync.WaitGroup) {
	defer wg.Done()
	s.RlTexturePros.EachEntity(func(entity ecs.Entity) bool {
		tr := s.RlTexturePros.GetUnsafe(entity)
		tint := s.Tints.GetUnsafe(entity)
		if tint == nil {
			return true
		}
		trTint := &tr.Tint
		trTint.A = tint.A
		trTint.R = tint.R
		trTint.G = tint.G
		trTint.B = tint.B
		return true
	})
}

func (s *RenderBogdanSystem) submitBatch(texID int, data []stdcomponents.RLTexturePro) {
	rl.BeginMode2D(s.camera)
	for i := range data {
		rl.DrawTexturePro(*data[i].Texture, data[i].Frame, data[i].Dest, data[i].Origin, data[i].Rotation, data[i].Tint)
	}
	rl.EndMode2D()
}
