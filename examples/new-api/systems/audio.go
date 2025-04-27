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
	"github.com/negrel/assert"
	"gomp/examples/new-api/components"
	"gomp/pkg/ecs"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewAudioSystem() AudioSystem {
	return AudioSystem{}
}

type AudioSystem struct {
	EntityManager *ecs.EntityManager
	SoundEffects  *components.SoundEffectsComponentManager
}

func (s *AudioSystem) Init() {
	rl.InitAudioDevice()
}

func (s *AudioSystem) Run(dt time.Duration) {
	s.SoundEffects.EachEntity()(func(entity ecs.Entity) bool {
		soundEffect := s.SoundEffects.GetUnsafe(entity)
		assert.NotNil(soundEffect)

		clip := soundEffect.Clip

		// check if clip is loaded
		if clip == nil || !rl.IsSoundValid(*clip) {
			return true
		}

		if !soundEffect.IsPlaying {
			if rl.IsSoundPlaying(*clip) {
				rl.StopSound(*clip)
				return true
			} else {
				*clip = rl.LoadSoundAlias(*clip)

				rl.PlaySound(*clip)
				soundEffect.IsPlaying = true
				return true
			}
		}

		// check if sound is over
		if !rl.IsSoundPlaying(*clip) && soundEffect.IsPlaying {
			if soundEffect.IsLooping {
				rl.PlaySound(*clip)
			} else {
				// sound is over, remove entity
				s.EntityManager.Delete(entity)
				// rl.UnloadSoundAlias(*clip) // TODO: this doesn't work https://github.com/gen2brain/raylib-go/issues/494
			}
		}

		return true
	})
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
	s.SoundEffects.EachEntity()(func(entity ecs.Entity) bool {
		soundEffect := s.SoundEffects.GetUnsafe(entity)
		assert.NotNil(soundEffect)

		clip := soundEffect.Clip

		// check if clip is loaded
		if clip == nil {
			return true
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
		return true
	})
}
func (s *AudioSettingsSystem) Destroy() {
}
