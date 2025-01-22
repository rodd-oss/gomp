/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type sineWaveController struct {
	sineWave ecs.EntityID
}

func (s *sineWaveController) Init(world *ecs.World) {
	s.sineWave = world.CreateEntity("sine wave")
	audibles := components.AudibleService.GetManager(world)
	audibles.Create(s.sineWave, components.Audible{
		State: components.Stoped,
		AudioStream: &components.SineWaveGenerator{
			SampleRate:   48000,
			SampleFactor: 220.0 / float64(48000),
			Phase:        0.0,
		},
		ChannelId: SFXChannelId,
		Volume:    1.0,
		Pan:       0.0,
	})
}
func (s *sineWaveController) Update(world *ecs.World) {
	audibles := components.AudibleService.GetManager(world)
	audible := audibles.Get(s.sineWave)

	if rl.IsKeyDown(rl.KeySpace) {
		audible.State = components.Playing
	} else {
		audible.State = components.Stoped
	}
}
func (s *sineWaveController) FixedUpdate(world *ecs.World) {}
func (s *sineWaveController) Destroy(world *ecs.World)     {}
