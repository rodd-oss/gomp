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

package stdcomponents

import (
	"gomp/pkg/ecs"
	"time"
)

type AnimationPlayer struct {
	First         uint8
	Last          uint8
	Current       uint8
	Speed         float32
	Loop          bool
	Vertical      bool
	ElapsedTime   time.Duration
	FrameDuration time.Duration
	State         AnimationState
	IsInitialized bool
}

type AnimationPlayerComponentManager = ecs.ComponentManager[AnimationPlayer]

func NewAnimationPlayerComponentManager(world *ecs.World) *ecs.ComponentManager[AnimationPlayer] {
	return ecs.NewComponentManager[AnimationPlayer](world, ANIMATION_PLAYER_ID)
}
