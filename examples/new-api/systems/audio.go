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
	s.SoundEffects.EachEntity(func(entity ecs.Entity) bool {
		soundEffect := s.SoundEffects.Get(entity)
		clip := soundEffect.Clip

		// check if clip is valid
		if clip.FrameCount == 0 {
			return true
		}

		if !soundEffect.IsPlaying {
			if rl.IsSoundPlaying(clip) {
				rl.StopSound(clip)
				return true
			} else {
				clip = rl.LoadSoundAlias(clip)

				rl.SetSoundVolume(clip, soundEffect.Volume)
				rl.SetSoundPitch(clip, soundEffect.Pitch)
				rl.SetSoundPan(clip, soundEffect.Pan)

				rl.PlaySound(clip)
				soundEffect.IsPlaying = true
				return true
			}
		}

		rl.SetSoundVolume(clip, soundEffect.Volume)
		rl.SetSoundPitch(clip, soundEffect.Pitch)
		rl.SetSoundPan(clip, soundEffect.Pan)

		// check if sound is over
		if !rl.IsSoundPlaying(clip) && soundEffect.IsPlaying {
			if soundEffect.IsLooping {
				rl.PlaySound(clip)
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
