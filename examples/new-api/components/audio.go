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

package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
)

type SoundEffect struct {
	Clip      *rl.Sound
	IsPlaying bool
	IsLooping bool
	// base is 1.0
	Volume float32
	// base is 1.0
	Pitch float32
	// base is 0.5. 1.0 is left, 0.0 is right
	Pan float32
}

type SoundEffectsComponentManager = ecs.ComponentManager[SoundEffect]

func NewSoundEffectsComponentManager() SoundEffectsComponentManager {
	return ecs.NewComponentManager[SoundEffect](SoundEffectManagerComponentId)
}

type SpatialAudio struct {
	// follows raylib rules. Base is 1.0
	Volume float32
	// follows raylib rules. Base is 0.5, 1.0 is left, 0.0 is right
	Pan float32
}

type SpatialAudioComponentManager = ecs.ComponentManager[SpatialAudio]

func NewSpatialAudioComponentManager() SpatialAudioComponentManager {
	return ecs.NewComponentManager[SpatialAudio](SpatialAudioManagerComponentId)
}
