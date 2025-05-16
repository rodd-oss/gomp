/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"image/color"
	"math"
	"time"
)

func NewDebugInfoSystem() DebugInfoSystem {
	return DebugInfoSystem{}
}

type DebugInfoSystem struct {
	EntityManager                      *ecs.EntityManager
	Positions                          *stdcomponents.PositionComponentManager
	Rotations                          *stdcomponents.RotationComponentManager
	Scales                             *stdcomponents.ScaleComponentManager
	ColliderSleepStateComponentManager *stdcomponents.ColliderSleepStateComponentManager
	CircleColliders                    *stdcomponents.CircleColliderComponentManager
	BoxColliders                       *stdcomponents.BoxColliderComponentManager
	Cameras                            *stdcomponents.CameraComponentManager
	RenderTexture2D                    *stdcomponents.FrameBuffer2DComponentManager
	TextureRect                        *components.TextureRectComponentManager
	Texture                            *stdcomponents.RLTextureProComponentManager
	Renderable                         *stdcomponents.RenderableComponentManager
	RenderOrders                       *stdcomponents.RenderOrderComponentManager
	AABBs                              *stdcomponents.AABBComponentManager
	Circle                             *components.PrimitiveCircleComponentManager
	CollisionChunks                    *stdcomponents.CollisionChunkComponentManager
	Tints                              *stdcomponents.TintComponentManager
	BvhTrees                           *stdcomponents.BvhTreeComponentManager

	debug       bool
	children    children
	liveParents []ecs.Entity
}

type children map[ecs.Entity]childType
type childType struct {
	id      ecs.Entity
	isAlive bool
}

func (s *DebugInfoSystem) Init() {
	s.children = make(children, s.BoxColliders.Len()+s.CircleColliders.Len())
	s.liveParents = make([]ecs.Entity, 0, s.BoxColliders.Len()+s.CircleColliders.Len())
}

func (s *DebugInfoSystem) Run(dt time.Duration) bool {
	if rl.IsKeyPressed(rl.KeyF6) {
		if !s.debug {
			s.BoxColliders.EachEntity(func(e ecs.Entity) bool {
				col := s.BoxColliders.GetUnsafe(e)
				scale := s.Scales.GetUnsafe(e)
				position := s.Positions.GetUnsafe(e)
				rotation := s.Rotations.GetUnsafe(e)

				x := position.XY.X
				y := position.XY.Y
				width := col.WH.X * scale.XY.X
				height := col.WH.Y * scale.XY.Y
				s.spawnRect(x, y, width, height, rl.Vector2{
					X: col.Offset.X * scale.XY.X,
					Y: col.Offset.Y * scale.XY.Y,
				}, float32(rotation.Degrees()), rl.DarkGreen, e)
				return true
			})
			s.CircleColliders.EachEntity(func(e ecs.Entity) bool {
				col := s.CircleColliders.GetUnsafe(e)
				scale := s.Scales.GetUnsafe(e)
				pos := s.Positions.GetUnsafe(e)

				circleColor := rl.DarkGreen
				isSleeping := s.ColliderSleepStateComponentManager.GetUnsafe(e)
				if isSleeping != nil {
					circleColor = rl.Blue
				}

				posWithOffset := pos.XY.Add(col.Offset.Mul(scale.XY))
				s.spawnCircle(posWithOffset.X, posWithOffset.Y, col.Radius*scale.XY.X, circleColor, e)
				return true
			})
			s.debug = true
		} else {
			s.Destroy()
		}

	}

	if s.debug {
		s.liveParents = s.liveParents[:0]

		// TODO: Parallelize this with future batches feature
		// Follow child to texture of parent box collider
		s.BoxColliders.EachEntity(func(e ecs.Entity) bool {
			parentAABB := s.AABBs.GetUnsafe(e)
			parentPosition := s.Positions.GetUnsafe(e)
			col := s.BoxColliders.GetUnsafe(e)
			scale := s.Scales.GetUnsafe(e)
			rotation := s.Rotations.GetUnsafe(e)

			s.liveParents = append(s.liveParents, e)
			child, ok := s.children[e]

			if ok {
				childAABB := s.AABBs.GetUnsafe(child.id)
				childRect := s.TextureRect.GetUnsafe(child.id)

				if parentAABB != nil && childAABB != nil {
					childAABB.Min = parentAABB.Min
					childAABB.Max = parentAABB.Max
					childRect.Dest.X = parentPosition.XY.X
					childRect.Dest.Y = parentPosition.XY.Y
					childRect.Dest.Width = col.WH.X * scale.XY.X
					childRect.Dest.Height = col.WH.Y * scale.XY.Y
					childRect.Origin.X = col.Offset.X * scale.XY.X
					childRect.Origin.Y = col.Offset.Y * scale.XY.Y
					childRect.Rotation = float32(rotation.Degrees())
				}

			} else {
				//TODO: defer spawning to non parallel function
				x := parentPosition.XY.X
				y := parentPosition.XY.Y
				width := col.WH.X * scale.XY.X
				height := col.WH.Y * scale.XY.Y
				s.spawnRect(x, y, width, height, rl.Vector2{
					X: col.Offset.X * scale.XY.X,
					Y: col.Offset.Y * scale.XY.Y,
				}, float32(rotation.Degrees()), rl.DarkGreen, e)
			}

			return true
		})
		// TODO: Parallelize this with future batches feature
		// Follow child to texture of parent circle collider
		s.CircleColliders.EachEntity(func(e ecs.Entity) bool {
			parentAABB := s.AABBs.GetUnsafe(e)
			pos := s.Positions.GetUnsafe(e)
			col := s.CircleColliders.GetUnsafe(e)
			scale := s.Scales.GetUnsafe(e)

			s.liveParents = append(s.liveParents, e)
			child, ok := s.children[e]
			if ok {
				childAABB := s.AABBs.GetUnsafe(child.id)
				childCircle := s.Circle.GetUnsafe(child.id)
				circleColor := rl.DarkGreen
				isSleeping := s.ColliderSleepStateComponentManager.GetUnsafe(e)
				if isSleeping != nil {
					circleColor = rl.Blue
				}
				if parentAABB != nil && childAABB != nil {
					childAABB.Min = parentAABB.Min
					childAABB.Max = parentAABB.Max
					posWithOffset := pos.XY.Add(col.Offset.Mul(scale.XY))
					childCircle.CenterX = posWithOffset.X
					childCircle.CenterY = posWithOffset.Y
					childCircle.Radius = col.Radius * scale.XY.X
					childCircle.Color = circleColor
				}

			} else {
				//TODO: defer spawning to non parallel function
				s.spawnCircle(pos.XY.X, pos.XY.Y, col.Radius*scale.XY.X, rl.DarkGreen, e)
			}
			return true
		})

		//Remove children that are not alive anymore
		for _, parent := range s.liveParents {
			if child, ok := s.children[parent]; ok {
				child.isAlive = true
				s.children[parent] = child
			}
		}
		for key, entity := range s.children {
			if !entity.isAlive {
				s.EntityManager.Delete(entity.id)
				delete(s.children, key)
			} else {
				entity.isAlive = false
				s.children[key] = entity
			}
		}
	}

	return false
}

func (s *DebugInfoSystem) spawnRect(x float32, y float32, width float32, height float32, origin rl.Vector2, rotation float32, tint color.RGBA, parent ecs.Entity) {
	if _, ok := s.children[parent]; !ok {
		childEntity := s.EntityManager.Create()
		s.TextureRect.Create(childEntity, components.TextureRect{
			Dest: rl.Rectangle{
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
			},
			Origin:   origin,
			Rotation: rotation,
			Color:    tint,
		})
		s.AABBs.Create(childEntity, stdcomponents.AABB{})
		s.Texture.Create(childEntity, stdcomponents.RLTexturePro{})
		s.Renderable.Create(childEntity, stdcomponents.Renderable{
			CameraMask: config.MinimapCameraLayer | config.MainCameraLayer,
			Type:       stdcomponents.SpriteRenderableType,
		})
		s.RenderOrders.Create(childEntity, stdcomponents.RenderOrder{
			CalculatedZ: math.MaxInt,
		})
		s.children[parent] = childType{
			id:      childEntity,
			isAlive: false,
		}
	}
}

func (s *DebugInfoSystem) spawnCircle(x float32, y float32, radius float32, circleColor color.RGBA, e ecs.Entity) {
	if _, ok := s.children[e]; !ok {
		childEntity := s.EntityManager.Create()
		s.Circle.Create(childEntity, components.TextureCircle{
			CenterX:  x,
			CenterY:  y,
			Radius:   radius,
			Rotation: 0,
			Origin:   rl.Vector2{},
			Color:    circleColor,
		})
		s.AABBs.Create(childEntity, stdcomponents.AABB{})
		s.Texture.Create(childEntity, stdcomponents.RLTexturePro{})
		s.Renderable.Create(childEntity, stdcomponents.Renderable{
			CameraMask: config.MinimapCameraLayer | config.MainCameraLayer,
			Type:       stdcomponents.SpriteRenderableType,
		})
		s.RenderOrders.Create(childEntity, stdcomponents.RenderOrder{
			CalculatedZ: math.MaxInt,
		})
		s.children[e] = childType{
			id:      childEntity,
			isAlive: false,
		}
	}
}

func (s *DebugInfoSystem) Destroy() {
	for _, child := range s.children {
		s.EntityManager.Delete(child.id)
	}
	clear(s.children)
	s.debug = false
}
