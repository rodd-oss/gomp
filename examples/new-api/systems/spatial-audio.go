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
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"time"

	"github.com/negrel/assert"
)

const (
	MinDistance = 10
)

func NewSpatialAudioSystem() SpatialAudioSystem {
	return SpatialAudioSystem{}
}

type SpatialAudioSystem struct {
	EntityManager         *ecs.EntityManager
	SoundEffects          *components.SoundEffectComponentManager
	Positions             *stdcomponents.PositionComponentManager
	SpatialAudio          *components.SpatialAudioComponentManager
	Cameras               *stdcomponents.CameraComponentManager
	numWorkers            int
	accSpatialAudioCreate [][]ecs.Entity
	accSpatialAudioDelete [][]ecs.Entity
	Engine                *core.Engine
}

func (s *SpatialAudioSystem) Init() {
	s.numWorkers = s.Engine.Pool().NumWorkers()
	s.accSpatialAudioCreate = make([][]ecs.Entity, s.numWorkers)
	s.accSpatialAudioDelete = make([][]ecs.Entity, s.numWorkers)
}

func (s *SpatialAudioSystem) Run(dt time.Duration) {
	var mainCamera ecs.Entity

	// TODO: Add listener component? Then we need position component on it...
	s.Cameras.EachEntity(func(entity ecs.Entity) bool {
		camera := s.Cameras.GetUnsafe(entity)
		assert.NotNil(camera)
		if camera.Layer == config.MainCameraLayer {
			mainCamera = entity
			return false
		}

		return true
	})

	if mainCamera == 0 {
		return
	}

	mainCameraComponent := s.Cameras.GetUnsafe(mainCamera)
	assert.NotNil(mainCameraComponent)

	var mainCameraPosition vectors.Vec2 = vectors.Vec2{
		X: mainCameraComponent.Camera2D.Target.X,
		Y: mainCameraComponent.Camera2D.Target.Y,
	}

	s.SoundEffects.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		position := s.Positions.Has(entity)

		if s.SpatialAudio.Has(entity) {
			if !position {
				s.accSpatialAudioDelete[workerId] = append(s.accSpatialAudioDelete[workerId], entity)
			}
		} else {
			if position {
				s.accSpatialAudioCreate[workerId] = append(s.accSpatialAudioCreate[workerId], entity)
			}
		}
	})

	for a := range s.accSpatialAudioCreate {
		for _, entity := range s.accSpatialAudioCreate[a] {
			s.SpatialAudio.Create(entity, components.SpatialAudio{
				Volume: 0,
				Pan:    0.5,
			})
		}
	}

	for a := range s.accSpatialAudioDelete {
		for _, entity := range s.accSpatialAudioDelete[a] {
			s.SpatialAudio.Delete(entity)
		}
	}

	s.SpatialAudio.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		spatialAudio := s.SpatialAudio.GetUnsafe(entity)
		assert.NotNil(spatialAudio)

		position := s.Positions.GetUnsafe(entity)
		assert.NotNil(position)

		spatialAudio.Volume = s.calculateVolume(
			mainCameraPosition,
			position.XY,
			mainCameraComponent.Camera2D.Offset.X*2,
		)
		spatialAudio.Pan = s.calculatePan(
			mainCameraPosition,
			position.XY,
			mainCameraComponent.Camera2D.Offset.X*2,
		)
	})

	for i := range s.accSpatialAudioCreate {
		s.accSpatialAudioCreate[i] = s.accSpatialAudioCreate[i][:0]
	}

	for i := range s.accSpatialAudioDelete {
		s.accSpatialAudioDelete[i] = s.accSpatialAudioDelete[i][:0]
	}

}
func (s *SpatialAudioSystem) Destroy() {
}

func (s *SpatialAudioSystem) calculatePan(listener vectors.Vec2, source vectors.Vec2, maxDistance float32) float32 {
	distance := listener.Distance(source)

	if distance < MinDistance {
		return 0.5
	}

	distanceX := float64(listener.X - source.X)

	pan := (1 + (distanceX / float64(maxDistance))) / 2

	return float32(math.Max(0, math.Min(1, pan)))
}

func (s *SpatialAudioSystem) calculateVolume(listener vectors.Vec2, source vectors.Vec2, maxDistance float32) float32 {
	distance := float64(listener.Distance(source))

	if distance < MinDistance {
		return 1
	}

	spatialVolume := 1 - (distance / float64(maxDistance)) // TODO: add ability to configure volume hearing distance
	volume := math.Max(0, math.Min(1, spatialVolume))

	return float32(math.Sin((volume * math.Pi) / 2)) // TODO: add ability to configure volume falloff. Current is easeOutSine
}
