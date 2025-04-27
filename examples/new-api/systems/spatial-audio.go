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
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"time"
)

const (
	MinDistance = 10
)

func NewSpatialAudioSystem() SpatialAudioSystem {
	return SpatialAudioSystem{}
}

type SpatialAudioSystem struct {
	EntityManager *ecs.EntityManager
	SoundEffects  *components.SoundEffectsComponentManager
	Positions     *stdcomponents.PositionComponentManager
	SpatialAudio  *components.SpatialAudioComponentManager
	Cameras       *stdcomponents.CameraComponentManager
}

func (s *SpatialAudioSystem) Init() {
}
func (s *SpatialAudioSystem) Run(dt time.Duration) {
	var mainCamera ecs.Entity

	// TODO: Add listener component? Then we need position component on it...
	s.Cameras.EachEntity()(func(entity ecs.Entity) bool {
		camera := s.Cameras.GetUnsafe(entity)
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
	var mainCameraPosition vectors.Vec2 = vectors.Vec2{
		X: mainCameraComponent.Camera2D.Target.X,
		Y: mainCameraComponent.Camera2D.Target.Y,
	}

	s.SoundEffects.EachEntity()(func(entity ecs.Entity) bool {
		soundEffect := s.SoundEffects.GetUnsafe(entity)

		clip := soundEffect.Clip

		if clip == nil {
			return true
		}

		position := s.Positions.GetUnsafe(entity)

		if position == nil {
			return true
		}

		spatialAudio := s.SpatialAudio.GetUnsafe(entity)

		if spatialAudio == nil {
			spatialAudio = s.SpatialAudio.Create(entity, components.SpatialAudio{
				Volume: 0,
				Pan:    0.5,
			})
		}

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

		return true
	})
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
