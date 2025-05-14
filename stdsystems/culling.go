/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"github.com/negrel/assert"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"gomp/vectors"
	"time"
)

func NewCullingSystem() CullingSystem {
	return CullingSystem{}
}

type CullingSystem struct {
	Renderables            *stdcomponents.RenderableComponentManager
	RenderVisible          *stdcomponents.RenderVisibleComponentManager
	RenderOrders           *stdcomponents.RenderOrderComponentManager
	Textures               *stdcomponents.RLTextureProComponentManager
	AABBs                  *stdcomponents.AABBComponentManager
	Cameras                *stdcomponents.CameraComponentManager
	RenderTexture2D        *stdcomponents.FrameBuffer2DComponentManager
	numWorkers             int
	accRenderVisibleCreate [][]ecs.Entity
	accRenderVisibleDelete [][]ecs.Entity
	Engine                 *core.Engine
}

func (s *CullingSystem) Init() {
	s.numWorkers = s.Engine.Pool().NumWorkers()
	s.accRenderVisibleCreate = make([][]ecs.Entity, s.numWorkers)
	s.accRenderVisibleDelete = make([][]ecs.Entity, s.numWorkers)
}

func (s *CullingSystem) Run(dt time.Duration) {
	for i := range s.accRenderVisibleCreate {
		s.accRenderVisibleCreate[i] = s.accRenderVisibleCreate[i][:0]
	}
	for i := range s.accRenderVisibleDelete {
		s.accRenderVisibleDelete[i] = s.accRenderVisibleDelete[i][:0]
	}

	s.Renderables.ProcessComponents(func(r *stdcomponents.Renderable, workerId worker.WorkerId) {
		r.Observed = false
	})

	s.Cameras.EachEntity(func(entity ecs.Entity) bool {
		camera := s.Cameras.GetUnsafe(entity)
		cameraRect := camera.Rect()

		s.Renderables.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
			renderable := s.Renderables.GetUnsafe(entity)
			assert.NotNil(renderable)

			texture := s.Textures.GetUnsafe(entity)
			assert.NotNil(texture)

			textureRect := texture.Rect()

			if s.intersects(cameraRect, textureRect) {
				renderable.Observed = true
			}
		})
		return true
	})

	s.Renderables.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		renderable := s.Renderables.GetUnsafe(entity)
		assert.NotNil(renderable)
		if !s.RenderVisible.Has(entity) {
			if renderable.Observed {
				s.accRenderVisibleCreate[workerId] = append(s.accRenderVisibleCreate[workerId], entity)
			}
		} else {
			if !renderable.Observed {
				s.accRenderVisibleDelete[workerId] = append(s.accRenderVisibleDelete[workerId], entity)
			}
		}
	})
	for a := range s.accRenderVisibleCreate {
		for _, entity := range s.accRenderVisibleCreate[a] {
			s.RenderVisible.Create(entity, stdcomponents.RenderVisible{})
		}
	}
	for a := range s.accRenderVisibleDelete {
		for _, entity := range s.accRenderVisibleDelete[a] {
			s.RenderVisible.Delete(entity)
		}
	}
}

func (_ *CullingSystem) intersects(rect1, rect2 vectors.Rectangle) bool {
	return rect1.X < rect2.X+rect2.Width &&
		rect1.X+rect1.Width > rect2.X &&
		rect1.Y < rect2.Y+rect2.Height &&
		rect1.Y+rect1.Height > rect2.Y
}

func (s *CullingSystem) Destroy() {
}
