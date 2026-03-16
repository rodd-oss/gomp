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
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"time"

	"github.com/negrel/assert"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewAudioSystem() AudioSystem {
	return AudioSystem{}
}

type AudioSystem struct {
	EntityManager         *ecs.EntityManager
	SoundEffects          *components.SoundEffectsComponentManager
	accSoundEffectsDelete [][]ecs.Entity
	numWorkers            int
	Engine                *core.Engine
}

func (s *AudioSystem) Init() {
	s.numWorkers = s.Engine.Pool().NumWorkers()
	s.accSoundEffectsDelete = make([][]ecs.Entity, s.numWorkers)
	rl.InitAudioDevice()
}

func (s *AudioSystem) Run(dt time.Duration) {
	for a := range s.accSoundEffectsDelete {
		s.accSoundEffectsDelete[a] = s.accSoundEffectsDelete[a][:0]
	}

	s.SoundEffects.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		soundEffect := s.SoundEffects.GetUnsafe(entity)
		assert.NotNil(soundEffect)

		clip := soundEffect.Clip

		// check if clip is loaded
		if clip == nil || !rl.IsSoundValid(*clip) {
			return
		}

		if !soundEffect.IsPlaying {
			if rl.IsSoundPlaying(*clip) {
				rl.StopSound(*clip)
				return
			} else {
				*clip = rl.LoadSoundAlias(*clip)

				rl.PlaySound(*clip)
				soundEffect.IsPlaying = true
				return
			}
		}

		// check if sound is over
		if !rl.IsSoundPlaying(*clip) && soundEffect.IsPlaying {
			if soundEffect.IsLooping {
				rl.PlaySound(*clip)
			} else {
				// sound is over, remove entity
				s.accSoundEffectsDelete[workerId] = append(s.accSoundEffectsDelete[workerId], entity)

				// rl.UnloadSoundAlias(*clip) // TODO: this doesn't work https://github.com/gen2brain/raylib-go/issues/494
			}
		}
	})

	for a := range s.accSoundEffectsDelete {
		for _, entity := range s.accSoundEffectsDelete[a] {
			s.EntityManager.Delete(entity)
		}
	}
}
func (s *AudioSystem) Destroy() {
	rl.CloseAudioDevice()
}

func NewAudioSettingsSystem() AudioSettingsSystem {
	return AudioSettingsSystem{}
}

type AudioSettingsSystem struct {
	EntityManager *ecs.EntityManager
	SoundEffects  *components.SoundEffectsComponentManager
	SpatialAudio  *components.SpatialAudioComponentManager
}

func (s *AudioSettingsSystem) Init() {}

func (s *AudioSettingsSystem) Run(dt time.Duration) {
	s.SoundEffects.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		soundEffect := s.SoundEffects.GetUnsafe(entity)
		assert.NotNil(soundEffect)

		clip := soundEffect.Clip

		// check if clip is loaded
		if clip == nil {
			return
		}

		spatialSettings := s.SpatialAudio.GetUnsafe(entity)

		if spatialSettings != nil {
			rl.SetSoundVolume(*clip, spatialSettings.Volume*soundEffect.Volume)
			rl.SetSoundPan(*clip, spatialSettings.Pan)
		} else {
			rl.SetSoundVolume(*clip, soundEffect.Volume)
			rl.SetSoundPan(*clip, soundEffect.Pan)
		}

		rl.SetSoundPitch(*clip, soundEffect.Pitch)
	})
}
func (s *AudioSettingsSystem) Destroy() {
}
