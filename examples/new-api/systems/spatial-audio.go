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
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewSpatialAudioSystem() SpatialAudioSystem {
	return SpatialAudioSystem{}
}

type SpatialAudioSystem struct {
	EntityManager *ecs.EntityManager
	SoundEffects  *components.SoundEffectsComponentManager
	Positions     *stdcomponents.PositionComponentManager
	Player        *components.PlayerTagComponentManager
}

func (s *SpatialAudioSystem) Init() {
}
func (s *SpatialAudioSystem) Run(dt time.Duration) {
	var player ecs.Entity = 0

	s.Player.EachEntity(func(entity ecs.Entity) bool {
		player = entity
		return false
	})

	if player == 0 {
		return
	}

	playerPos := s.Positions.Get(player)

	if playerPos == nil {
		return
	}

	s.SoundEffects.EachEntity(func(entity ecs.Entity) bool {
		soundEffect := s.SoundEffects.Get(entity)

		clip := soundEffect.Clip

		if clip.FrameCount == 0 {
			return true
		}

		if !soundEffect.IsPlaying {
			return true
		}

		position := s.Positions.Get(entity)

		if position == nil {
			return true
		}

		pan := s.calculatePan(playerPos.XY, position.XY)
		rl.SetSoundPan(clip, pan)

		return true
	})
}
func (s *SpatialAudioSystem) Destroy() {
}

func (s *SpatialAudioSystem) calculatePan(listener vectors.Vec2, source vectors.Vec2) float32 {
	dx := float64(source.X - listener.X)
	dy := float64(source.Y - listener.Y)
	distanceSq := dx*dx + dy*dy

	// Если источник и слушатель в одной точке
	if distanceSq < 1e-9 { // Используем квадрат расстояния для оптимизации
		return 0.5
	}

	distance := math.Sqrt(distanceSq)
	pan := 0.5 - (dx / (2 * distance))

	return float32(math.Max(0, math.Min(1, pan)))
}
